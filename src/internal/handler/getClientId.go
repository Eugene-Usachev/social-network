package handler

import (
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	keyId      = "id"
	AuthHeader = "Authorization"
)

func (handler *Handler) getClientIdFromHeaders(ctx *gin.Context) int {
	accessToken := ctx.GetHeader(AuthHeader)
	if accessToken == "" {
		return 0
	}

	idBytes, err := handler.accessTokenConverter.ParseToken(accessToken)
	if err != nil {
		return 0
	}

	return fb.B2I(idBytes)
}

func (handler *Handler) CheckAuth(ctx *gin.Context) {
	id := handler.getClientIdFromHeaders(ctx)
	if id == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set(keyId, id)
	ctx.Next()
}

func GetClientId(ctx *gin.Context) int {
	return ctx.MustGet(keyId).(int)
}
