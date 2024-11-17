package handler

import (
	"log"
	"net/http"

	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/gin-gonic/gin"
)

const (
	keyID      = "id"
	AuthHeader = "Authorization"
)

func (handler *Handler) getClientIDFromHeaders(ctx *gin.Context) int {
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
	id := handler.getClientIDFromHeaders(ctx)
	if id == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	ctx.Set(keyID, id)
	ctx.Next()
}

func GetClientID(ctx *gin.Context) int {
	id, isExist := ctx.Get(keyID)

	if !isExist {
		log.Println("[BUG] Client id is not exist!")

		return 0
	}

	idInt, isValid := id.(int)

	if !isValid {
		log.Println("[BUG] Client id is not int!")

		return 0
	}

	return idInt
}
