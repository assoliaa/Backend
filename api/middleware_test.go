package api

import (
	"Backend/token"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func addAuthorization(
	t *testing.T, 
	request *http.Request, 
	tokenMaker token.Maker,
	autorizationType string,
	id uuid.UUID,
	duration time.Duration,
){
	token, err := tokenMaker.CreateToken(id, duration)
	assert.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", autorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)

}


func TestAuthMiddleware(t *testing.T){
	userId := uuid.New()
	testCases :=[]struct{
		name string
		setAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
		name: "OK",
		setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			addAuthorization(t, request, tokenMaker, authBearer, userId,  time.Minute)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
			assert.Equal(t, http.StatusOK, recorder.Code)
		}, 
		},
		{
			name: "NoAuthorization",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedType",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, "unsupported", uuid.Nil,  time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthFormat",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, "", uuid.New(),  time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, authBearer, userId,  -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	server, router:= createServer(t)
	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T){
			recorder :=httptest.NewRecorder()
			
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			assert.NoError(t, err)
			
			tc.setAuth(t, req, server.tokenMaker)
			router.ServeHTTP(recorder, req)
			
			tc.checkResponse(t, recorder)
		})

	}
}

func createServer(t *testing.T)(*Server, *gin.Engine){
	server:= newTestServer(t, nil)

			authPath:="/auth"

			server.router.GET(
				authPath,
				authMidddleware(server.tokenMaker),
				func(ctx *gin.Context){
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)
	return server, server.router
}

// почему authMiddleware иногда тутпит и роутер