package upload

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImage(ctx context.Context, file *multipart.FileHeader, cloudName, apiKey, apiSecret string) (string, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload file to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
