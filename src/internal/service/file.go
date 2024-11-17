package service

import (
	"context"

	"github.com/Eugune-Usachev/social-network/src/internal/filestorage"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
)

type FileService struct {
	repository *repository.Repository
	fs         filestorage.FileStorage
}

var _ File = (*FileService)(nil)

func NewFileService(
	repository *repository.Repository,
	fs filestorage.FileStorage,
) *FileService {
	return &FileService{
		repository,
		fs,
	}
}

func (fileService *FileService) GetPresignedURL(ctx context.Context, userID int, filePath string) (string, error) {
	return fileService.fs.LoadFile(ctx, userID, filePath)
}
