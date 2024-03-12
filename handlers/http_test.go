package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v4"

	"eth-fetcher/database/models"
	"eth-fetcher/handlers"
	"eth-fetcher/handlers/mocks"
)

var tx1 = &models.Transaction{
	TxHash:      "hash1",
	TxStatus:    1,
	BlockHash:   "blockHash1",
	BlockNumber: 7976373,
	From:        "from1",
	To:          null.NewString("to1", true),
	LogsCount:   1,
	Input:       "0x1",
	Value:       50000000000000000,
}

var tx2 = &models.Transaction{
	TxHash:      "hash2",
	TxStatus:    1,
	BlockHash:   "blockHash2",
	BlockNumber: 7976373,
	From:        "from2",
	To:          null.NewString("to2", true),
	LogsCount:   1,
	Input:       "0x2",
	Value:       50000000000000000,
}

func Setup(t *testing.T) (*mocks.APP, *mocks.Auth, *handlers.HTTP) {
	app := mocks.NewAPP(t)
	auth := mocks.NewAuth(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sLog := logger.Sugar()
	httpHandler := handlers.NewHTTP(app, "8080", auth, sLog)
	httpHandler.InitRoutes()
	return app, auth, httpHandler
}

func TestHTTP_GetAllTransactions(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/all", nil)
	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("GetAllTransactions").Return([]*models.Transaction{
		tx1,
	}, nil)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 1)
	assert.Equal(t, tx1, response.Transactions[0])

}

func TestHTTP_GetUserTransactions(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/my", nil)
	r.Header.Set("AUTH_TOKEN", "123")
	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("AUTH_TOKEN") != ""
	})).Return("user1", nil)
	app.On("GetUserTransactions", "user1").Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 2)
	assert.Equal(t, tx1, response.Transactions[0])
	assert.Equal(t, tx2, response.Transactions[1])
}

func TestHTTP_GetTransactionsByRLPHandler(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/eth/f90110b842307839623266366133633265316165643263636366393262613636366332326430353361643064386135646137616131666435343737646364363537376234353234b842307835613537653330353163623932653264343832353135623037653762336431383531373232613734363534363537626436346131346333396361336639636632b842307837316239653262343464343034393863303861363239383866616337373664306561633062356239363133633337663966366639613462383838613862303537b842307863356639366266316235346433333134343235643233373962643737643765643465363434663763366538343961373438333230323862333238643464373938", nil)
	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("GetTransactionsByHashes", []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", "0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057", "0xc5f96bf1b54d3314425d2379bd77d7ed4e644f7c6e849a74832028b328d4d798"}).Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 2)
	assert.Equal(t, tx1, response.Transactions[0])
	assert.Equal(t, tx2, response.Transactions[1])
}

func TestHTTP_GetTransactionsByHashesHandler(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/eth?transactionHashes=0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", nil)
	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("GetTransactionsByHashes", []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524"}).Return([]*models.Transaction{
		tx1,
	}, nil)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 1)
	assert.Equal(t, tx1, response.Transactions[0])
}

func TestHTTP_GetTransactionsByHashesHandler_ForUser(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/eth?transactionHashes=0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", nil)
	r.Header.Set("AUTH_TOKEN", "value")

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("AUTH_TOKEN") != ""
	})).Return("user1", nil)
	app.On("GetTransactionsByHashes", []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524"}).Return([]*models.Transaction{
		tx1,
	}, nil)
	app.On("AddUserTransactions", "user1", []*models.Transaction{
		tx1,
	}).Return(nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 1)
	assert.Equal(t, tx1, response.Transactions[0])
}

func TestHTTP_GetTransactionsByHashesHandler_Error(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/eth?transactionHashes=0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", nil)
	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("GetTransactionsByHashes", []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524"}).Return(nil, assert.AnError)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, assert.AnError.Error(), response.Msg)
}

func TestHTTP_GetTransactionsByRLPHandler_ForUser(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	r, _ := http.NewRequest("GET", "/api/eth/f90110b842307839623266366133633265316165643263636366393262613636366332326430353361643064386135646137616131666435343737646364363537376234353234b842307835613537653330353163623932653264343832353135623037653762336431383531373232613734363534363537626436346131346333396361336639636632b842307837316239653262343464343034393863303861363239383866616337373664306561633062356239363133633337663966366639613462383838613862303537b842307863356639366266316235346433333134343235643233373962643737643765643465363434663763366538343961373438333230323862333238643464373938", nil)
	r.Header.Set("AUTH_TOKEN", "value")

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return r.Header.Get("AUTH_TOKEN") != ""
	})).Return("user1", nil)
	app.On("GetTransactionsByHashes", []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", "0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057", "0xc5f96bf1b54d3314425d2379bd77d7ed4e644f7c6e849a74832028b328d4d798"}).Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)
	app.On("AddUserTransactions", "user1", []*models.Transaction{
		tx1,
		tx2,
	}).Return(nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.GetTransactionsResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Len(t, response.Transactions, 2)
	assert.Equal(t, tx1, response.Transactions[0])
	assert.Equal(t, tx2, response.Transactions[1])
}

func TestHTTP_AuthenticateHandler(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	body := []byte(`{"username": "user1", "password": "password1"}`)
	r, _ := http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("CheckUserCredentials", "user1", "password1").Return(&models.User{
		ID:       "user1",
		Username: "user1",
		Password: "password1",
	}, nil)
	auth.On("GenerateToken", "user1").Return("user1.token", nil)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.AuthenticateResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "user1.token", response.Token)
}

func TestHTTP_AuthenticateHandler_Error(t *testing.T) {
	app, auth, httpHandler := Setup(t)
	httpHandler.InitRoutes()
	body := []byte(`{"username": "user1", "password": "password1"}`)
	r, _ := http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)
	app.On("CheckUserCredentials", "user1", "password1").Return(nil, assert.AnError)
	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, assert.AnError.Error(), response.Msg)
}
