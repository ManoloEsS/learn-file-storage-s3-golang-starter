package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	//get videoID from url path
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	//validate bearer token from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//start upload
	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	//parse multipart request into int64
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse request body", err)
		return
	}

	//get file data and header from request
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get file", err)
		return
	}
	defer file.Close()

	//get media type from thumbnail request header
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
	}
	//only accept jpeg or png files
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Not a valid file", err)
		return
	}

	//get video data from database using video id
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video in database", err)
		return
	}

	//validate user is author of video
	if userID != video.UserID {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to update this video", err)
		return
	}

	//get file path
	videoFilePath := getAssetPath(mediaType)
	fullPath := cfg.getAssetDiskPath(videoFilePath)

	//create file
	videoFile, err := os.Create(fullPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file in disk", err)
		return
	}

	//copy image file contents into videoFile
	defer videoFile.Close()
	if _, err = io.Copy(videoFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file", err)
		return
	}

	//assign path to file to ThumbnailURL in the retrieved video
	thumbnailURL := cfg.getAssetURL(videoFilePath)
	video.ThumbnailURL = &thumbnailURL

	//update video data with added thumbnail data url
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
