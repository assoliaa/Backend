package api

import (
	db "Backend/db/sqlc"
	"Backend/token"

	//"database/sql"
	"errors"

	//"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createAccoutRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
    authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	params := db.InsertAccountParams{
		OwnerID:   authPayload.UserId,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.InsertAccount(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountReq struct{
	ID string `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accountID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
    return
    }
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if account.OwnerID != authPayload.UserId{
		err = errors.New("wrong user to get")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountReq struct {
	PageId int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func(server *Server)listAccounts(ctx *gin.Context){
	var req listAccountReq
	if err:=ctx.ShouldBindQuery(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	params:= db.GetAccountsParams{
		OwnerID: authPayload.UserId,
		Limit:req.PageSize,
		Offset:(req.PageId-1)*req.PageSize,
	}
	accounts, err := server.store.GetAccounts(ctx, params)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
    ctx.JSON(http.StatusOK, accounts)
}