package handler

import (
	"errors"
	"github.com/Eugune-Usachev/social-network/src/customErrors"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (handler *Handler) GetSmallProfile(ctx *gin.Context) {
	userId64, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if userId64 < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId := int(userId64)

	profile, err := handler.service.GetSmallProfile(ctx, userId)
	if err != nil {
		if errors.Is(err, customErrors.NotFound) {
			handler.AbortWithError(ctx, http.StatusNotFound, err)
			return
		}
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":        profile.Name,
		"second_name": profile.SecondName,
		"avatar":      profile.Avatar,
		"description": profile.Description,
		"Birthday":    profile.Birthday,
		"gender":      profile.Gender,
		"email":       profile.Email,
	})
}

func (handler *Handler) UpdateSmallProfile(ctx *gin.Context) {
	clientId := GetClientId(ctx)

	var newSmallProfile model.UpdateSmallProfile

	err := ctx.BindJSON(&newSmallProfile)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)
		return
	}

	err = handler.service.UpdateSmallProfile(ctx, clientId, &newSmallProfile)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) GetInfo(ctx *gin.Context) {
	userId64, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)
		return
	}

	if userId64 < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId := int(userId64)

	info, err := handler.service.GetInfo(ctx, userId)
	if err != nil {
		if errors.Is(err, customErrors.NotFound) {
			handler.AbortWithError(ctx, http.StatusNotFound, err)
			return
		}
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.String(http.StatusOK, "%s", info)
}

func (handler *Handler) UpdateInfo(ctx *gin.Context) {
	clientId := GetClientId(ctx)

	var newInfo string

	err := ctx.Bind(&newInfo)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)
		return
	}

	err = handler.service.UpdateInfo(ctx, clientId, newInfo)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
