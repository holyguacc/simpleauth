package api

import (
	"errors"
	"fmt"
	"net/http"
	"simpleauth/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizatationHeaderKey = "authorization"
	authoriaztionTypeBearer  = "bearer"
	authoriaztionPayloadKey  = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizatationHeader := ctx.GetHeader(authorizatationHeaderKey)
		if len(authorizatationHeader) == 0 {
			err := errors.New("authoriaztion header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizatationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authoriaztion header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizatationType := strings.ToLower(fields[0])
		if authorizatationType != authoriaztionTypeBearer {
			err := fmt.Errorf("wrong authorization type %s", authorizatationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accsessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accsessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authoriaztionPayloadKey, payload)
		ctx.Next()
	}
}
