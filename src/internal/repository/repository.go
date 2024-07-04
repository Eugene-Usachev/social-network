package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"social-network/src/internal/model"
)

type Auth interface {
	IsEmailBusy(ctx context.Context, email string) (isExists bool, err error)
	SignUp(ctx context.Context, model model.SignUp) (id int, err error)
	SignIn(ctx context.Context, email, password string) (id int, err error)
	AuthenticateUser(ctx context.Context, id int, password string) (wasAuthenticated bool, err error)
}

type Profile interface{}

type Playlist interface{}

type Chat interface{}

type Song interface{}

type Message interface{}

type Post interface{}

type Repository struct {
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
}

func NewRepository(postgres *pgxpool.Pool) *Repository {
	return &Repository{
		Auth: NewAuthRepository(postgres),
	}
}
