package repository

import (
	"context"
	"errors"

	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	postgres *pgxpool.Pool
}

var _ Auth = (*AuthRepository)(nil)

func NewAuthRepository(postgres *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		postgres: postgres,
	}
}

func (authRepository AuthRepository) IsEmailBusy(ctx context.Context, email string) (isExists bool, err error) {
	const query = "SELECT TRUE FROM users WHERE email = $1"

	if err = authRepository.postgres.QueryRow(ctx, query, email).Scan(&isExists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return isExists, nil
}

func (authRepository AuthRepository) SignUp(ctx context.Context, model *model.SignUp) (id int, err error) {
	const query = "INSERT INTO users (name, second_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id"

	row := authRepository.postgres.QueryRow(
		ctx,
		query,
		model.GetName(),
		model.GetSecondName(),
		model.GetEmail(),
		model.GetPassword(),
	)

	if err = row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

var ErrInvalidPassword = errors.New("invalid password")

func (authRepository AuthRepository) SignIn(ctx context.Context, email, password string) (id int, err error) {
	const query = "SELECT id, password FROM users WHERE email = $1"

	var selectedPassword string

	row := authRepository.postgres.QueryRow(ctx, query, email)

	if err = row.Scan(&id, &selectedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return id, ErrNotFound
		}

		return id, err
	}

	if password != selectedPassword {
		return id, ErrInvalidPassword
	}

	return id, err
}

func (authRepository AuthRepository) AuthenticateUser(
	ctx context.Context,
	id int,
	password string,
) (wasAuthenticated bool, err error) {
	const query = "SELECT TRUE FROM users WHERE id = $1 AND password = $2"

	if err = authRepository.postgres.QueryRow(ctx, query, id, password).Scan(&wasAuthenticated); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return wasAuthenticated, nil
}
