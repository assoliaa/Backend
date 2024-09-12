package api

import (
	"Backend/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authBearer = "bearer"
	authPayloadKey ="authorization_key"
)

func authMidddleware(tokenMaker token.Maker)gin.HandlerFunc{
	return func(ctx *gin.Context){
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) ==0{
			err := errors.New("authorization header not provided")
		    ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields)<2{
			err := errors.New("Invalid Fornmat")
		    ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorzType := strings.ToLower(fields[0])

		if authorzType != authBearer{
			err := errors.New("Invalid auth type")
		    ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err!=nil{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}


