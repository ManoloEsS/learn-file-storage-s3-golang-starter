package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	//create a random 32-byte slice
	videoBytesSrc := make([]byte, 32)
	rand.Read(videoBytesSrc)

	//encode the 32-byte slice into base 64 URL encoding
	videoBytesDst := base64.RawURLEncoding.EncodeToString(videoBytesSrc)

	//get file extension
	ext := mediaTypeToExt(mediaType)

	return fmt.Sprintf("%s%s", videoBytesDst, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)

	objectInput := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignedRequest, err := presignClient.PresignGetObject(context.Background(), &objectInput, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return presignedRequest.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	paramsSplit := strings.Split(*video.VideoURL, ",")
	if len(paramsSplit) < 2 {
		return video, nil
	}

	presignedURL, err := generatePresignedURL(cfg.s3Client, paramsSplit[0], paramsSplit[1], time.Minute*2)
	if err != nil {
		return video, err
	}

	video.VideoURL = &presignedURL

	return video, nil
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
