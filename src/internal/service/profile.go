package service

import (
	"context"

	"github.com/Eugune-Usachev/social-network/src/internal/filestorage"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
)

type ProfileService struct {
	repository *repository.Repository
	fs         filestorage.FileStorage
}

var _ Profile = (*ProfileService)(nil)

func NewProfileService(repository *repository.Repository, fs filestorage.FileStorage) *ProfileService {
	return &ProfileService{
		repository: repository,
		fs:         fs,
	}
}

func (profileService ProfileService) UploadAvatar(ctx context.Context, id int, file filestorage.UploadedFile) error {
	name, err := profileService.fs.UploadFile(ctx, file, id, nil)
	if err != nil {
		return err
	}

	return profileService.repository.UpdateAvatar(ctx, id, name)
}

func (profileService ProfileService) GetSmallProfile(
	ctx context.Context,
	id int,
) (model.SmallProfile, error) {
	return profileService.repository.GetSmallProfile(ctx, id)
}

func (
	profileService ProfileService) UpdateSmallProfile(
	ctx context.Context,
	id int,
	profile *model.UpdateSmallProfile,
) error {
	return profileService.repository.UpdateSmallProfile(ctx, id, profile)
}

func (profileService ProfileService) GetInfo(ctx context.Context, id int) (string, error) {
	return profileService.repository.GetInfo(ctx, id)
}

func (profileService ProfileService) UpdateInfo(ctx context.Context, id int, info string) error {
	return profileService.repository.UpdateInfo(ctx, id, info)
}
