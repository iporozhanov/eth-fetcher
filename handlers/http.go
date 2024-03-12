package handlers

import (
	"encoding/json"
	"eth-fetcher/app"
	"eth-fetcher/database/models"

	"eth-fetcher/helpers/rlp"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTP struct {
	apiPort string
	app     APP
	auth    Auth
	log     *zap.SugaredLogger
	Router  *mux.Router
}

type APP interface {
	GetTransactionsByHashes(transactionHashes []string) ([]*models.Transaction, error)
	GetAllTransactions() ([]*models.Transaction, error)
	AddUserTransactions(userID string, transactions []*models.Transaction) error
	GetUserTransactions(userID string) ([]*models.Transaction, error)
	CheckUserCredentials(username, password string) (*models.User, error)
}

type Auth interface {
	AuthenticateRequest(r *http.Request) (string, error)
	GenerateToken(subject string) (string, error)
}

type Session struct {
	UserID string
}

type handleFunc func(Session, *http.Request) (any, error)

func (h *HTTP) InitRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/eth", h.HandleHTTPRequest(h.GetTransactionsHandler)).Methods("GET")
	router.HandleFunc("/api/all", h.HandleHTTPRequest(h.GetTransactionsHandler)).Methods("GET")
	router.HandleFunc("/api/eth/{rlphex}", h.HandleHTTPRequest(h.GetTransactionsByRLPHandler)).Methods("GET")
	router.HandleFunc("/api/authenticate", h.HandleHTTPRequest(h.AuthenticateHandler)).Methods("POST")
	router.HandleFunc("/api/my", h.HandleHTTPRequest(h.GetUserTransactions)).Methods("GET")

	h.Router = router
}

func (h *HTTP) Run() {
	h.log.Infof("API server starting on port %s", h.apiPort)
	http.ListenAndServe(fmt.Sprintf(":%s", h.apiPort), h.Router)
}

func NewHTTP(app APP, apiPort string, auth Auth, log *zap.SugaredLogger) *HTTP {
	return &HTTP{app: app, apiPort: apiPort, auth: auth, log: log}
}

func (a *HTTP) GetTransactionsByRLPHandler(s Session, r *http.Request) (any, error) {
	rlphex := mux.Vars(r)["rlphex"]

	transactionHashes, err := rlp.RLPtoStrings(rlphex)
	if err != nil {
		return nil, err
	}

	transactions, err := a.app.GetTransactionsByHashes(transactionHashes)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	if s.UserID != "" {
		err := a.app.AddUserTransactions(s.UserID, transactions)
		if err != nil {
			a.log.Errorf("error adding user transactions: %v", err)
		}
	}

	response := GetTransactionsResponse{Transactions: transactions}

	return response, nil
}

func (a *HTTP) GetTransactionsHandler(s Session, r *http.Request) (any, error) {
	transactionHashes := r.URL.Query()["transactionHashes"]
	var (
		transactions []*models.Transaction
		err          error
	)

	if len(transactionHashes) == 0 {
		transactions, err = a.app.GetAllTransactions()
	} else {
		transactions, err = a.app.GetTransactionsByHashes(transactionHashes)
	}

	if err != nil {
		if err == app.ErrBadRequest {
			return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusBadRequest}
		}

		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if s.UserID != "" {
		err := a.app.AddUserTransactions(s.UserID, transactions)
		if err != nil {
			a.log.Errorf("error adding user transactions: %v", err)
		}
	}

	response := GetTransactionsResponse{Transactions: transactions}

	return response, nil
}

func (a *HTTP) GetUserTransactions(s Session, r *http.Request) (any, error) {
	if s.UserID == "" {
		return nil, &ErrorResponse{Msg: "invalid credentials", Code: http.StatusUnauthorized}
	}

	transactions, err := a.app.GetUserTransactions(s.UserID)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	response := GetTransactionsResponse{Transactions: transactions}

	return response, nil
}

func (a *HTTP) HandleHTTPRequest(fn handleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := Session{}

		userID, err := a.auth.AuthenticateRequest(r)
		if err == nil {
			session.UserID = userID
		}

		response, err := fn(session, r)
		if err != nil {
			if e, ok := err.(*ErrorResponse); ok {
				w.WriteHeader(e.Code)
				json.NewEncoder(w).Encode(e)
				return
			}
			e := &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(e)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (a *HTTP) AuthenticateHandler(s Session, r *http.Request) (any, error) {
	var req AuthenticateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	user, err := a.app.CheckUserCredentials(req.Username, req.Password)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusUnauthorized}
	}

	token, err := a.auth.GenerateToken(user.ID)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusUnauthorized}
	}

	return AuthenticateResponse{Token: token}, nil
}
