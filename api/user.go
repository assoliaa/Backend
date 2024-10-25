package api

import (
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type createUserRequest struct {
	Username		string `json:"username" binding:"required,alphanum"`
	Password		string `json:"password" binding:"required,min=6"`
	FullName		string `json:"full_name" binding:"required,alphanum"`
    Email			string `json:"email" binding:"required,email"`
}
type userResponse struct {
	Username			string `json:"username" binding:"required,alphanum"`
	FullName    		string `json:"full_name" binding:"required,alphanum"`
    Email				string `json:"email" binding:"required,email"`
    PasswordChangedAt 	time.Time `json:"password_changed_at"`
	CreatedAt 			time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err :=utils.HashPassword(req.Password)
    if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := db.InsertUserParams{
		Username:      req.Username,
		HashPassword:  hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.InsertUser(ctx, params)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp:= newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	UserId		uuid.UUID `json:"username" binding:"required,alphanum"`
	Password	string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`
}

func (server *Server)loginUser(ctx *gin.Context){
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.UserId)
	if err!=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
	}
	err =utils.CheckPassword(req.Password, user.HashPassword)
	if err !=nil{
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
	}
	accessToken, err := server. tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp:=loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)

}