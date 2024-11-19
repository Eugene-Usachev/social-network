package repository

import (
	"context"
	"errors"

	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrivateFileMetadataRepository struct {
	postgres *pgxpool.Pool
	logger   logger.Logger
}

var _ PrivateFileMetadata = (*PrivateFileMetadataRepository)(nil)

func NewPrivateFileMetadataRepository(postgres *pgxpool.Pool, logger logger.Logger) *PrivateFileMetadataRepository {
	return &PrivateFileMetadataRepository{
		postgres: postgres,
		logger:   logger,
	}
}

func (privateFileMetadataRepository PrivateFileMetadataRepository) CheckAccess(
	ctx context.Context,
	filePath string,
	userID int,
) (bool, error) {
	const query = `
		SELECT TRUE 
		FROM file_metadata 
		WHERE file_path = $1 AND $2 = ANY(authorized_users)
	`

	var hasAccess bool

	row := privateFileMetadataRepository.postgres.QueryRow(ctx, query, filePath, userID)
	if err := row.Scan(&hasAccess); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrNotFound
		}

		return false, err
	}

	return hasAccess, nil
}

func (privateFileMetadataRepository PrivateFileMetadataRepository) SaveFileMetadata(
	ctx context.Context,
	filePath string,
	authorizedUsers []int,
) error {
	const query = `
		INSERT INTO file_metadata 
		(file_path, authorized_users)
		VALUES ($1, $2)
	`

	_, err := privateFileMetadataRepository.postgres.Exec(ctx, query, filePath, authorizedUsers)

	return err
}

func (privateFileMetadataRepository PrivateFileMetadataRepository) CheckFileExists(
	ctx context.Context,
	filePath string,
) (bool, error) {
	const query = `
		SELECT TRUE 
		FROM file_metadata 
		WHERE file_path = $1
	`

	var exists bool

	row := privateFileMetadataRepository.postgres.QueryRow(ctx, query, filePath)
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}
