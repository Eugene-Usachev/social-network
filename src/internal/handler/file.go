package handler

import (
	"errors"
	"net/http"

	"github.com/Eugune-Usachev/social-network/src/internal/filestorage"
	"github.com/gin-gonic/gin"
)

var ErrPathEmpty = errors.New("path is empty")

func (handler *Handler) UploadFile(ctx *gin.Context) {
	path := ctx.Query("path")
	if len(path) == 0 {
		handler.AbortWithError(ctx, http.StatusBadRequest, ErrPathEmpty)

		return
	}

	url, err := handler.service.GetPresignedURL(ctx, 0, path)
	if err != nil {
		var status int
		if errors.Is(err, filestorage.ErrNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		handler.AbortWithError(ctx, status, err)

		return
	}

	ctx.String(http.StatusOK, url)
}

func (handler *Handler) UploadFileWithAuth(ctx *gin.Context) {
	clientID := GetClientID(ctx)
	path := ctx.Query("path")

	if len(path) == 0 {
		handler.AbortWithError(ctx, http.StatusBadRequest, ErrPathEmpty)

		return
	}

	url, err := handler.service.GetPresignedURL(ctx, clientID, path)
	if err != nil {
		var status int
		if errors.Is(err, filestorage.ErrNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		handler.AbortWithError(ctx, status, err)

		return
	}

	ctx.String(http.StatusOK, url)
}
