package service

import (
	"context"
	"github.com/Eugene-Usachev/fst"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
)

type Auth interface {
	SignUp(ctx context.Context, model *model.SignUp) (id int, accessToken, refreshToken string, err error)
	SignIn(ctx context.Context, email, password string) (id int, accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, id int, password string) (accessToken, refreshToken string, err error)
}

type Profile interface {
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
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
}

func NewService(repository *repository.Repository, accessConverter *fst.EncodedConverter, refreshConverter *fst.EncodedConverter) *Service {
	return &Service{
		Auth:    NewAuthService(repository, accessConverter, refreshConverter),
		Profile: NewProfileService(repository),
	}
}
