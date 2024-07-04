package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-network/src/internal/model"
	"social-network/src/internal/repository"
	"social-network/src/internal/service"
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

	id, accessToken, refreshToken, err = handler.service.Auth.SignUp(ctx, signUpModel)
	if err != nil {
		if errors.Is(err, service.EmailIsBusy) {
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

	id, accessToken, refreshToken, err = handler.service.Auth.SignIn(ctx, signInModel.Email, signInModel.Password)
	if err != nil {
		if errors.Is(err, repository.NotFound) || errors.Is(err, repository.InvalidPassword) {
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

func (handler *Handler) Refresh(ctx *gin.Context) {
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

	accessToken, refreshToken, err = handler.service.Auth.RefreshTokens(ctx, refreshModel.ID, refreshModel.Password)
	if err != nil {
		if errors.Is(err, service.Unauthorized) {
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
