package repository

import (
	"context"
	"github.com/Eugune-Usachev/social-network/src/internal/repository/cache"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth interface {
	IsEmailBusy(ctx context.Context, email string) (isExists bool, err error)
	SignUp(ctx context.Context, model *model.SignUp) (id int, err error)
	SignIn(ctx context.Context, email, password string) (id int, err error)
	AuthenticateUser(ctx context.Context, id int, password string) (wasAuthenticated bool, err error)
}

type Profile interface {
	// TODO Update Avatar
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

type Repository struct {
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
}

func NewRepository(postgres *pgxpool.Pool, cache cache.Cache, logger logger.Logger) *Repository {
	return &Repository{
		Auth:    NewAuthRepository(postgres),
		Profile: NewProfileRepository(postgres, cache, logger),
	}
}
