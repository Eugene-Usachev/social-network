package handler

import (
	"errors"
	"github.com/Eugune-Usachev/social-network/src/customErrors"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/internal/service"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
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

	if signUpModel.Email == "" || signUpModel.Password == "" || signUpModel.Name == "" || signUpModel.SecondName == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, accessToken, refreshToken, err = handler.service.Auth.SignUp(ctx, &signUpModel)
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

	if signInModel.Email == "" || signInModel.Password == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, accessToken, refreshToken, err = handler.service.Auth.SignIn(ctx, signInModel.Email, signInModel.Password)
	if err != nil {
		if errors.Is(err, customErrors.NotFound) || errors.Is(err, repository.InvalidPassword) {
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

	if refreshModel.Id == -1 || refreshModel.Password == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err = handler.service.Auth.RefreshTokens(ctx, int(refreshModel.Id), refreshModel.Password)
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
