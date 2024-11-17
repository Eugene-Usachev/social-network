package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Eugune-Usachev/social-network/src/internal/filestorage"
	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/gin-gonic/gin"
)

func (handler *Handler) UploadAvatar(ctx *gin.Context) {
	clientID := GetClientID(ctx)

	file, err := ctx.FormFile("avatar")
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)
	}

	fileContent, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})

		return
	}
	defer fileContent.Close()

	uploadedFile := filestorage.UploadedFile{
		File:       fileContent,
		FileHeader: file,
	}

	err = handler.service.UploadAvatar(ctx, clientID, uploadedFile)
	if err != nil {
		var statusCode int

		if errors.Is(err, filestorage.ErrFileNameEmpty) || errors.Is(err, filestorage.ErrFileNameTooLong) {
			statusCode = http.StatusBadRequest
		} else {
			statusCode = http.StatusInternalServerError
		}

		handler.AbortWithError(ctx, statusCode, err)

		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) GetSmallProfile(ctx *gin.Context) {
	userID64, err := strconv.ParseInt(ctx.Param("userID"), 10, 64)
	if err != nil || userID64 < 1 {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	userID := int(userID64)

	profile, err := handler.service.GetSmallProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			handler.AbortWithError(ctx, http.StatusNotFound, err)

			return
		}

		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":        profile.GetName(),
		"second_name": profile.GetSecondName(),
		"avatar":      profile.GetAvatar(),
		"description": profile.GetDescription(),
		"birthday":    profile.GetBirthday(),
		"gender":      profile.GetGender(),
		"email":       profile.GetEmail(),
	})
}

func (handler *Handler) UpdateSmallProfile(ctx *gin.Context) {
	clientID := GetClientID(ctx)

	var newSmallProfile model.UpdateSmallProfile

	err := ctx.BindJSON(&newSmallProfile)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	err = handler.service.UpdateSmallProfile(ctx, clientID, &newSmallProfile)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) GetInfo(ctx *gin.Context) {
	userID64, err := strconv.ParseInt(ctx.Param("userID"), 10, 64)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	if userID64 < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID := int(userID64)

	info, err := handler.service.GetInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			handler.AbortWithError(ctx, http.StatusNotFound, err)

			return
		}

		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.String(http.StatusOK, "%s", info)
}

func (handler *Handler) UpdateInfo(ctx *gin.Context) {
	clientID := GetClientID(ctx)

	var newInfo string

	err := ctx.Bind(&newInfo)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusBadRequest, err)

		return
	}

	err = handler.service.UpdateInfo(ctx, clientID, newInfo)
	if err != nil {
		handler.AbortWithError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.Status(http.StatusOK)
}
