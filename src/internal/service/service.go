package service

import (
	"context"

	"github.com/Eugene-Usachev/fst"
	"github.com/Eugune-Usachev/social-network/src/internal/filestorage"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
)

type File interface {
	// GetPresignedURL returns a presigned URL subject to access. If the file `public`, id can be any integer.
	GetPresignedURL(ctx context.Context, userID int, filePath string) (url string, err error)
}

type Auth interface {
	SignUp(ctx context.Context, model *model.SignUp) (id int, accessToken, refreshToken string, err error)
	SignIn(ctx context.Context, email, password string) (id int, accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, id int, password string) (accessToken, refreshToken string, err error)
}

type Profile interface {
	UploadAvatar(ctx context.Context, id int, file filestorage.UploadedFile) (err error)
	GetSmallProfile(ctx context.Context, id int) (profile model.SmallProfile, err error)
	UpdateSmallProfile(ctx context.Context, id int, profile *model.UpdateSmallProfile) (err error)
	GetInfo(ctx context.Context, id int) (info string, err error)
	UpdateInfo(ctx context.Context, id int, info string) (err error)
}

type Playlist interface{}

type Chat interface{}

type Song interface{}

type Message interface{}

type Post interface{}

type Service struct {
	File
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
}

func NewService(
	repository *repository.Repository,
	fs filestorage.FileStorage,
	accessConverter *fst.EncodedConverter,
	refreshConverter *fst.EncodedConverter,
) *Service {
	return &Service{
		File:    NewFileService(repository, fs),
		Auth:    NewAuthService(repository, accessConverter, refreshConverter),
		Profile: NewProfileService(repository, fs),
	}
}
