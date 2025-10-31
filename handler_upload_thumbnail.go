package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

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

	//get file and header from request
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get file", err)
		return
	}
	defer file.Close()

	//get media type from thumbnail request header
	mediaType := header.Header.Get("Content-Type")
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", err)
	}

	//read image data from file into bytes
	imageData, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading file", err)
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

	//encode imageData into base64 string
	encodedImageData := base64.StdEncoding.EncodeToString(imageData)

	//create dataURL
	dataURL := fmt.Sprintf("data:%s;base64,%s", mediaType, encodedImageData)

	//assign dataURL to ThumbnailURL in the retrieved video
	video.ThumbnailURL = &dataURL

	//update video data with added thumbnail data url
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
