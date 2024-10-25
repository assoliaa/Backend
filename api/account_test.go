package api

import (
	mockdb "Backend/db/mock"
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"Backend/token"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := createUser(t)
	account := randomAccount(user.ID)

	testCases := []struct {
		name          string
		accountID     string
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID.String(),
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID.String(),
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID.String(),
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Формируем URL с accountID
			url := fmt.Sprintf("/accounts/%s", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			tc.setAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}


func TestCreateAccountAPI(t *testing.T) {
	user, _ := createUser(t)
	account := randomAccount(user.ID)
	testCases :=[]struct{
        name string  
		body gin.H
		setAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        buildStubs func(store *mockdb.MockStore)
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
        
    }{
		{
			name: "OK",
			body: gin.H{
			   "currency": account.Currency,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				params := db.InsertAccountParams{
					OwnerID:    account.OwnerID,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
				InsertAccount(gomock.Any(), gomock.Eq(params)).
				Times(1).
				Return(account, nil)
			},
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusOK, recorder.Code)
			    requireBodyMatch(t, recorder.Body, account)
			},
		},
		{
			name:"InternalServerError",
			body: gin.H{
				"currency": account.Currency,
			 },
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().InsertAccount(gomock.Any(), gomock.Any()).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:"NoAuthorization",
			body: gin.H{
				"currency": account.Currency,
			 },
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().InsertAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:"InvalidCurrency",
			body: gin.H{
				"owner_id": account.OwnerID,
				"currency": 8,
			 },
			setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
				addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().InsertAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder){
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i:= range testCases{
		tc:= testCases[i]
	
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
	
			store:= mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
	
			server:=newTestServer(t, store)
	
			recorder :=httptest.NewRecorder()
	
			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			assert.NoError(t, err)

			tc.setAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccounts(t *testing.T){
    accounts :=[]db.Account{}
    n:=5
	user, _ := createUser(t)
    for i:=0; i<n; i++{
        accounts = append(accounts, randomAccount(user.ID))
    }
    type Query struct {
        pageID int
        pageSize int
    }
    testCases :=[]struct{
        name string
        query Query
		setAuth func(t *testing.T, request *http.Request, tokenMaker  token.Maker)
        buildStubs func(store *mockdb.MockStore)
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
        
    }{
    {   name: "OK",
        query: Query{
            pageID:1,
            pageSize:n,
        },
		setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
		},
        buildStubs: func(store *mockdb.MockStore){
			params := db.GetAccountsParams{
				OwnerID:  user.ID,
			    Limit:  int32(n),
				Offset: 0,
			}
            store.EXPECT().
            GetAccounts(gomock.Any(), gomock.Eq(params)).
            Times(1).
            Return(accounts, nil)
        },
        checkResponse : func(t *testing.T, recorder *httptest.ResponseRecorder){
            assert.Equal(t, http.StatusOK, recorder.Code)
            requireBodyMatch(t, recorder.Body, accounts)
        },
    },
    {  name : "InternalServerError",
       query: Query{
        pageID:1,
        pageSize:n,
       }, 
	   setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
		addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
	},
       buildStubs: func(store *mockdb.MockStore){
         store.EXPECT().
         GetAccounts(gomock.Any(), gomock.Any()).
         Times(1). // here
         Return([]db.Account{}, sql.ErrConnDone)
       },
       checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
          assert.Equal(t, http.StatusInternalServerError, recorder.Code)
       },
    },
	
    {
        name:"InvalidPageID",
        query: Query{
            pageID: -1,
            pageSize:n,
        },
		setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
		},
        buildStubs: func(store *mockdb.MockStore){
            store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).
            Times(0)
        },
        checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
            assert.Equal(t, http.StatusBadRequest, recorder.Code)
        },
    },
    {
        name:"InvalidPageSize",
        query: Query{
            pageID:1,
            pageSize:1000,
        },
		setAuth: func(t *testing.T, request *http.Request, tokenMaker  token.Maker){
			addAuthorization(t, request, tokenMaker, authBearer, user.ID,  time.Minute)
		},
        buildStubs: func(store *mockdb.MockStore){
            store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).
            Times(0)
        },
        checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
            assert.Equal(t, http.StatusBadRequest, recorder.Code)
        },
    },
}
   for i:= range testCases{
    tc:= testCases[i]

    t.Run(tc.name, func(t *testing.T){
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()

        store:= mockdb.NewMockStore(ctrl)
        tc.buildStubs(store)

        server:=newTestServer(t, store)

        recorder :=httptest.NewRecorder()

        url:="/accounts"
        req, err :=http.NewRequest(http.MethodGet, url, nil)

        assert.NoError(t, err)

        q := req.URL.Query()
        q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
        q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
        req.URL.RawQuery = q.Encode()
        
		tc.setAuth(t, req, server.tokenMaker)
        server.router.ServeHTTP(recorder, req)
        tc.checkResponse(t, recorder)
    })
}
}



func randomAccount(userId uuid.UUID) db.Account {
	return db.Account{
		ID:       uuid.New(),
		OwnerID:  userId,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}


func requireBodyMatch(t *testing.T, body *bytes.Buffer, expected interface{}) {
    var err error

    switch exp := expected.(type) {
    case []db.Account:
        var actualAccounts []db.Account
        err = json.NewDecoder(body).Decode(&actualAccounts)
        assert.NoError(t, err)
        assert.ElementsMatch(t, exp, actualAccounts)

    case db.Account:
        var actualAccount db.Account
        err = json.NewDecoder(body).Decode(&actualAccount)
        assert.NoError(t, err)
        assert.Equal(t, exp, actualAccount)

    default:
        t.Fatalf("Unsupported type for comparison: %T", expected)
    }
}
