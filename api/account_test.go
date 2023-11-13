package api

import (
	"bytes"
	"database/sql"
	"db"
	"dbmock"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"token"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	userTest, _ := createFakeUser()
	coin := db.Coin{
		Name: "BTC",
	}
	accountTest := db.Account{
		ID:       int64(utils.CreateRandomInt(100)),
		Username: userTest.Username,
		Coin:     coin.Name,
	}
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *dbmock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": accountTest.Username,
				"coin":     accountTest.Coin,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizatoin(t, request, tokenMaker, authorizationTypeBearer, userTest.Username, time.Minute)
			},
			buildStubs: func(store *dbmock.MockStore) {
				arg := db.CreateAccountParams{
					Username: accountTest.Username,
					Coin:     accountTest.Coin,
				}
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(accountTest.Coin)).
					Times(1).
					Return(coin, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accountTest, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, accountTest)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"username": accountTest.Username,
				"coin":     accountTest.Coin,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizatoin(t, request, tokenMaker, authorizationTypeBearer, userTest.Username, time.Minute)
			},
			buildStubs: func(store *dbmock.MockStore) {
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(accountTest.Coin)).
					Times(1).
					Return(coin, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Not Valid Coin",
			body: gin.H{
				"username": accountTest.Username,
				"coin":     "NotValidCoin",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizatoin(t, request, tokenMaker, authorizationTypeBearer, userTest.Username, time.Minute)
			},
			buildStubs: func(store *dbmock.MockStore) {
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(accountTest.Coin)).
					Times(0).
					Return(coin, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			crtl := gomock.NewController(t)
			defer crtl.Finish()

			store := dbmock.NewMockStore(crtl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/account/create", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account.Username, gotAccount.Username)
	require.Equal(t, account.Coin, gotAccount.Coin)
}
