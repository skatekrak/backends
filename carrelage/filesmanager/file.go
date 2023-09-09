package filesmanager

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type FilesManager struct {
	cld *cloudinary.Cloudinary
}

func New(url string) (*FilesManager, error) {
	c, err := cloudinary.NewFromURL(url)

	return &FilesManager{cld: c}, err
}

func (f *FilesManager) Upload(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := f.cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "tmp/profiles"})

	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", errors.New("cloudinary: No response")
	}

	return resp.SecureURL, nil
}
