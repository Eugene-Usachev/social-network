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
	UpdateAvatar(ctx context.Context, id int, avatar string) (err error)
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

type PrivateFileMetadata interface {
	// CheckAccess checks access to the `private` file
	CheckAccess(ctx context.Context, filePath string, userID int) (bool, error)
	// SaveFileMetadata saves metadata for the `private` file (including access control)
	SaveFileMetadata(ctx context.Context, filePath string, authorizedUsers []int) error
	// CheckFileExists checks if the `private` file path already exists in the database
	CheckFileExists(ctx context.Context, filePath string) (bool, error)
}

type Repository struct {
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
	PrivateFileMetadata
}

func NewRepository(postgres *pgxpool.Pool, cache cache.Cache, logger logger.Logger) *Repository {
	return &Repository{
		Auth:                NewAuthRepository(postgres),
		Profile:             NewProfileRepository(postgres, cache, logger),
		PrivateFileMetadata: NewPrivateFileMetadataRepository(postgres, logger),
	}
}
