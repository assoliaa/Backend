package api

import (
	"errors"
	"fmt"
	"net/http"

	db "Backend/db/sqlc"
	"Backend/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transferRequest struct {
	FromAccountID uuid.UUID  `json:"from_account_id" binding:"required"`
	ToAccountID   uuid.UUID  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}


func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
    fromAccount, valid := server.goodAccountCurrency(ctx, req.FromAccountID, req.Currency)
    if !valid{
		return
	}
	authPayload :=ctx.MustGet(authPayloadKey).(*token.Payload)
	if fromAccount.OwnerID != authPayload.UserId {
		err := errors.New("its a wrong user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	
	_, valid = server.goodAccountCurrency(ctx, req.ToAccountID, req.Currency)
    if !valid{
		return
	}
    
	if fromAccount.Balance <= 0 || fromAccount.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "not enough money"})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
    
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) goodAccountCurrency(ctx *gin.Context, accountID uuid.UUID, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	
	return account, true
}
