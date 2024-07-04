package service

import (
	"context"
	"github.com/Eugene-Usachev/fst"
	"social-network/src/internal/model"
	"social-network/src/internal/repository"
)

type Auth interface {
	SignUp(ctx context.Context, model model.SignUp) (id int, accessToken, refreshToken string, err error)
	SignIn(ctx context.Context, email, password string) (id int, accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, id int, password string) (accessToken, refreshToken string, err error)
}

type Profile interface{}

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
		Auth: NewAuthService(repository, accessConverter, refreshConverter),
	}
}
