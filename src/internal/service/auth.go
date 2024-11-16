package service

import (
	"context"
	"errors"

	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
)

type AuthService struct {
	repository       *repository.Repository
	accessConverter  *fst.EncodedConverter
	refreshConverter *fst.EncodedConverter
}

var _ Auth = (*AuthService)(nil)

func NewAuthService(
	repository *repository.Repository,
	accessConverter *fst.EncodedConverter,
	refreshConverter *fst.EncodedConverter,
) *AuthService {
	return &AuthService{
		repository,
		accessConverter,
		refreshConverter,
	}
}

var ErrEmailIsBusy = errors.New("email is busy")

func (authService AuthService) SignUp(ctx context.Context, model *model.SignUp) (int, string, string, error) {
	var (
		accessToken, refreshToken string
		id                        int
	)

	isEmailBusy, err := authService.repository.Auth.IsEmailBusy(ctx, model.GetEmail())
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	if isEmailBusy {
		return id, accessToken, refreshToken, ErrEmailIsBusy
	}

	id, err = authService.repository.Auth.SignUp(ctx, model)
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(model.GetPassword()))

	return id, accessToken, refreshToken, nil
}

func (authService AuthService) SignIn(ctx context.Context, email, password string) (int, string, string, error) {
	var accessToken, refreshToken string

	id, err := authService.repository.Auth.SignIn(ctx, email, password)
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(password))

	return id, accessToken, refreshToken, nil
}

var ErrUnauthorized = errors.New("unauthorized")

func (authService AuthService) RefreshTokens(ctx context.Context, id int, password string) (string, string, error) {
	var (
		wasAuthenticated          bool
		accessToken, refreshToken string
		err                       error
	)

	wasAuthenticated, err = authService.repository.Auth.AuthenticateUser(ctx, id, password)
	if err != nil {
		return accessToken, refreshToken, err
	}

	if !wasAuthenticated {
		return accessToken, refreshToken, ErrUnauthorized
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(password))

	return accessToken, refreshToken, nil
}
