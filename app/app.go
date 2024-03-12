package app

import (
	"eth-fetcher/database/models"
	"sync"

	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("unauthorized")

var ErrBadRequest = errors.New("bad request")

// App represents the application.
type App struct {
	db  DB
	tg  TransactionGetter
	Log *zap.SugaredLogger
}

// TransactionGetter is an interface for getting transactions.
type TransactionGetter interface {
	GetTransaction(txID string) (*models.Transaction, error)
	Close()
}

// DB is an interface for interacting with the database.
type DB interface {
	SaveTransaction(transaction *models.Transaction) error
	GetTransactionsByHashes(hashes []string) ([]*models.Transaction, error)
	GetAllTransactions() ([]*models.Transaction, error)
	AddUserTransactions(userID string, transactions []*models.Transaction) error
	GetUserTransactions(userID string) ([]*models.Transaction, error)
	GetUserByUsername(username string) (*models.User, error)
	Close() error
}

// NewApp creates a new instance of the App.
func NewApp(db DB, tg TransactionGetter, log *zap.SugaredLogger) *App {
	return &App{
		db:  db,
		tg:  tg,
		Log: log,
	}
}

// Shutdown shuts down the application.
func (a *App) Shutdown() {
	a.Log.Info("shutting down app")
	err := a.db.Close()
	if err != nil {
		a.Log.Errorf("error closing db: %v", err)
	}
	a.tg.Close()
	a.Log.Sync()
}

// GetTransactionsByHashes retrieves transactions by their hashes from the database or
// asynchronously from the eth node.
func (a *App) GetTransactionsByHashes(transactionHashes []string) ([]*models.Transaction, error) {
	if len(transactionHashes) == 0 {
		return nil, ErrBadRequest
	}

	txHashesMap := make(map[string]struct{})
	for _, txHash := range transactionHashes {
		txHashesMap[txHash] = struct{}{}
	}

	txs, err := a.db.GetTransactionsByHashes(transactionHashes)
	if err != nil {
		a.Log.Errorf("error getting transactions from db: %v", err)
	}
	// remove transactions that were found in the db
	for _, tx := range txs {
		delete(txHashesMap, tx.TxHash)
	}

	var wg sync.WaitGroup
	wg.Add(len(txHashesMap))

	tChan := make(chan *models.Transaction, len(txHashesMap))
	ec := make(chan error, len(transactionHashes))

	for transactionHash := range txHashesMap {
		// get transaction from the eth node asynchronously
		go func(transactionHash string, tChan chan *models.Transaction, ec chan error) {
			defer wg.Done()
			transaction, err := a.tg.GetTransaction(transactionHash)
			if err != nil {
				ec <- err
				return
			}
			tChan <- transaction
		}(transactionHash, tChan, ec)
	}

	wg.Wait()
	close(tChan)
	close(ec)

	for tx := range tChan {
		err = a.db.SaveTransaction(tx)

		if err != nil {
			a.Log.Errorf("error saving transaction %s: %v", tx.TxHash, err)
			continue
		}

		txs = append(txs, tx)
	}

	for err := range ec {
		a.Log.Error(err)
	}

	return txs, nil
}

// GetAllTransactions retrieves all transactions.
func (a *App) GetAllTransactions() ([]*models.Transaction, error) {
	return a.db.GetAllTransactions()
}

// GetUserTransactions retrieves transactions for a specific user.
func (a *App) GetUserTransactions(userID string) ([]*models.Transaction, error) {
	if userID == "" {
		return nil, ErrUnauthorized
	}
	return a.db.GetUserTransactions(userID)
}

// AddUserTransactions adds transactions for a specific user.
func (a *App) AddUserTransactions(userID string, transactions []*models.Transaction) error {
	return a.db.AddUserTransactions(userID, transactions)
}

// CheckUserCredentials checks the user credentials.
func (a *App) CheckUserCredentials(username, password string) (*models.User, error) {
	user, err := a.db.GetUserByUsername(username)
	if err != nil {

		return nil, ErrUnauthorized
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}
