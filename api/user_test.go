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
	"reflect"
	"testing"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createFakeUser() (db.User, string) {
	pass := utils.CreateRandomName()
	hash, _ := utils.CreateHashPassword(pass)
	userTest := db.User{
		Username:       utils.CreateRandomName(),
		HashedPassword: hash,
		FullName:       utils.CreateRandomName(),
		Email:          utils.CreateRandomEmail(),
	}
	return userTest, pass
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := utils.VerifyHashPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(arg, e.arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return "Matches with password"
}

func EQCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}

func TestUserCreate(t *testing.T) {
	userTest, pass := createFakeUser()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *dbmock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  userTest.Username,
				"full_name": userTest.FullName,
				"password":  pass,
				"email":     userTest.Email,
			},
			buildStubs: func(store *dbmock.MockStore) {
				arg := db.CreateUserParams{
					Username: userTest.Username,
					FullName: userTest.FullName,
					Email:    userTest.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EQCreateUserParams(arg, pass)).
					Times(1).
					Return(userTest, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, userTest)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"username":  userTest.Username,
				"full_name": userTest.FullName,
				"password":  pass,
				"email":     userTest.Email,
			},
			buildStubs: func(store *dbmock.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Duplicate Username",
			body: gin.H{
				"username":  userTest.Username,
				"full_name": userTest.FullName,
				"password":  pass,
				"email":     userTest.Email,
			},
			buildStubs: func(store *dbmock.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalidUser#123",
				"full_name": userTest.FullName,
				"password":  pass,
				"email":     userTest.Email,
			},
			buildStubs: func(store *dbmock.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
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
			request, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
}
