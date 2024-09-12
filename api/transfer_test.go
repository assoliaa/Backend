package api

import (
	mockdb "Backend/db/mock"
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"Backend/token"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)


func TestCreateTransfer(t *testing.T){
	amount := int64(10)

	user1, _ := createUser(t)
	user2, _ := createUser(t)
	user3, _ := createUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user3.Username)

	account1.Currency = utils.USD
	account2.Currency = utils.USD
	account3.Currency = utils.EUR


	testCases := []struct{
			name string
			body gin.H
			setAuth func(t *testing.T, request *http.Request, tokenMaker  token.Maker)
			buildStubs func(store *mockdb.MockStore)		
			checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		
		}{
			{
				name: "OK",
				body: gin.H{
					"from_account_id" : account1.ID,
					"to_account_id" : account2.ID,
					"amount" : amount,
					"currency" : utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore){
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				    params :=db.TransferTxParams{
						FromAccountID: account1.ID,
						ToAccountID: account2.ID,
						Amount: amount,
					}
					store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(params)).Times(1)
				},
				checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder){
					assert.Equal(t, http.StatusOK, recorder.Code)			
				},
			},
			{
				name: "FromAccountNotFound",
				body: gin.H{
					"from_account_id": account1.ID,
					"to_account_id":   account2.ID,
					"amount":          amount,
					"currency":        utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, db.ErrRecordNotFound)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(t *testing. T , recorder *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, recorder.Code)
				},

			},
			{
				name: "ToAccountNotFound",
				body: gin.H{
					"from_account_id": account1.ID,
					"to_account_id":   account2.ID,
					"amount":          amount,
					"currency":        utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, db.ErrRecordNotFound)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(t *testing. T , recorder *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			
			{
				name : "ToAcountCurrencyMisMatch",
				body: gin.H{
					"from_account_id" : account1.ID,
					"to_account_id" : account3.ID,
					"amount" : amount,
					"currency":utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store * mockdb.MockStore){
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
					assert.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
			{
				name: "InvalidCurrency",
				body : gin.H{
					"from_account_id": account1.ID,
					"to_account_id": account2.ID,
					 "amount": amount, 
					 "currency": "ANY",

				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore){
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(t *testing.T, recorder  *httptest.ResponseRecorder){
					assert.Equal(t, http.StatusBadRequest, recorder.Code)
				},

			},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        utils.USD,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
			{
				name: "GetAccountError",
				body: gin.H{
					"from_account_id": account1.ID,
					"to_account_id":   account2.ID,
					"amount":          amount,
					"currency":        utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
			{
				name: "TransferTxError",
				body: gin.H{
					"from_account_id": account1.ID,
					"to_account_id":   account2.ID,
					"amount":          amount,
					"currency":        utils.USD,
				},
				setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
					addAuthorization(t, request, tokenMaker, authBearer, user1.Username,  time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
					store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
		}

		for i := range testCases {
			tc := testCases[i]

			t.Run(tc.name, func(t *testing.T){
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store:=mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := newTestServer(t, store)
				recorder :=httptest.NewRecorder()

				data, err :=json.Marshal(tc.body)
				assert.NoError(t,err)

				url := "/transfers"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				assert.NoError(t, err)
                
				tc.setAuth(t, request, server.tokenMaker)
				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)

			})
		}
	}

