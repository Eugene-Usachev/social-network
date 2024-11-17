package handler

import (
	"errors"
	"net/http"

	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/internal/service"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/gin-gonic/gin"
)

func (handler *Handler) SingUp(ctx *gin.Context) {
	var (
		signUpModel  model.SignUp
		id           int
		accessToken  string
		refreshToken string
		err          error
	)

	err = ctx.BindJSON(&signUpModel)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	isValidRequest := signUpModel.GetEmail() != "" ||
		signUpModel.GetPassword() != "" ||
		signUpModel.GetName() != "" ||
		signUpModel.GetSecondName() != ""
	if !isValidRequest {
		ctx.AbortWithStatus(http.StatusBadRequest)

		return
	}

	id, accessToken, refreshToken, err = handler.service.Auth.SignUp(ctx, &signUpModel)
	if err != nil {
		if errors.Is(err, service.ErrEmailIsBusy) {
			ctx.AbortWithStatus(http.StatusBadRequest)

			return
		}

		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":            id,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (handler *Handler) SignIn(ctx *gin.Context) {
	var (
		signInModel  model.SignIn
		id           int
		accessToken  string
		refreshToken string
		err          error
	)

	err = ctx.BindJSON(&signInModel)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	if signInModel.GetEmail() == "" || signInModel.GetPassword() == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)

		return
	}

	id, accessToken, refreshToken, err = handler.service.Auth.SignIn(
		ctx,
		signInModel.GetEmail(),
		signInModel.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, repository.ErrInvalidPassword) {
			ctx.AbortWithStatus(http.StatusBadRequest)

			return
		}

		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            id,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (handler *Handler) RefreshTokens(ctx *gin.Context) {
	var (
		refreshModel model.RefreshTokens
		accessToken  string
		refreshToken string
		err          error
	)

	err = ctx.BindJSON(&refreshModel)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	if refreshModel.GetId() == -1 || refreshModel.GetPassword() == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)

		return
	}

	accessToken, refreshToken, err = handler.service.Auth.RefreshTokens(
		ctx,
		int(refreshModel.GetId()),
		refreshModel.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
