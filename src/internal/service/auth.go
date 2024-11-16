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

func NewAuthService(repository *repository.Repository, accessConverter *fst.EncodedConverter, refreshConverter *fst.EncodedConverter) *AuthService {
	return &AuthService{
		repository,
		accessConverter,
		refreshConverter,
	}
}

var (
	EmailIsBusy = errors.New("email is busy")
)

func (authService AuthService) SignUp(ctx context.Context, model *model.SignUp) (id int, accessToken, refreshToken string, err error) {
	isEmailBusy, err := authService.repository.Auth.IsEmailBusy(ctx, model.Email)
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	if isEmailBusy {
		return id, accessToken, refreshToken, EmailIsBusy
	}

	id, err = authService.repository.Auth.SignUp(ctx, model)
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(model.Password))

	return id, accessToken, refreshToken, nil
}

func (authService AuthService) SignIn(ctx context.Context, email, password string) (id int, accessToken, refreshToken string, err error) {
	id, err = authService.repository.Auth.SignIn(ctx, email, password)
	if err != nil {
		return id, accessToken, refreshToken, err
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(password))

	return id, accessToken, refreshToken, nil
}

var Unauthorized = errors.New("unauthorized")

func (authService AuthService) RefreshTokens(ctx context.Context, id int, password string) (accessToken, refreshToken string, err error) {
	var wasAuthenticated bool

	wasAuthenticated, err = authService.repository.Auth.AuthenticateUser(ctx, id, password)
	if err != nil {
		return accessToken, refreshToken, err
	}

	if !wasAuthenticated {
		return accessToken, refreshToken, Unauthorized
	}

	accessToken = authService.accessConverter.NewToken(fb.I2B(id))
	refreshToken = authService.refreshConverter.NewToken(fb.S2B(password))

	return accessToken, refreshToken, nil
}
