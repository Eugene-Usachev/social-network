package filestorage

import (
	"context"
	"errors"
	"mime/multipart"
)

type UploadedFile struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type FileStorage interface {
	// UploadFile uploads a file. If len(authorizedUserIDs) > 0, the file will be private
	// and accessible only by the authorized users.
	UploadFile(ctx context.Context, file UploadedFile, ownerID int, authorizedUserIDs []int) (string, error)
	// LoadFile returns a presigned URL for the file, with checking if the user has access to it
	LoadFile(ctx context.Context, userID int, filePath string) (string, error)
}

var ErrForbidden = errors.New("you are not allowed to access this file")
