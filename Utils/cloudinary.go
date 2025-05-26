package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

var cloudinaryURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file for Cloudinary:", err)
	}

	cloudinaryURL = os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatal("CLOUDINARY_URL is not set! Check your .env file.")
	}
}

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", fmt.Errorf("cloudinary config error: %v", err)
	}

	uploadParams := uploader.UploadParams{PublicID: filename}
	resp, err := cld.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload error: %v", err)
	}

	return resp.SecureURL, nil
}
