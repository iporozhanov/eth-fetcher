package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v4"

	"eth-fetcher/app"
	"eth-fetcher/app/mocks"
	"eth-fetcher/database/models"
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

var tx3 = &models.Transaction{
	TxHash:      "hash3",
	TxStatus:    1,
	BlockHash:   "blockHash3",
	BlockNumber: 7976373,
	From:        "from3",
	To:          null.NewString("to3", true),
	LogsCount:   1,
	Input:       "0x3",
	Value:       50000000000000000,
}

func Setup(t *testing.T) (*mocks.DB, *mocks.TransactionGetter, *app.App) {
	db := mocks.NewDB(t)
	tg := mocks.NewTransactionGetter(t)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sLog := logger.Sugar()
	return db, tg, app.NewApp(db, tg, sLog)
}

func TestApp_GetTransactionsByHashes(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, tg, app := Setup(t)

	// Set up test data
	transactionHashes := []string{"hash1", "hash2", "hash3"}

	// Mock the database's GetTransactionsByHashes method
	db.EXPECT().GetTransactionsByHashes(transactionHashes).Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)

	tg.EXPECT().GetTransaction("hash3").Return(tx3, nil)

	db.EXPECT().SaveTransaction(tx3).Return(nil)

	// Call the GetTransactionsByHashes method
	transactions, err := app.GetTransactionsByHashes(transactionHashes)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 3)
	assert.Equal(t, "hash1", transactions[0].TxHash)
	assert.Equal(t, "hash2", transactions[1].TxHash)
	assert.Equal(t, "hash3", transactions[2].TxHash)

	// Verify that the mock database's GetTransactionsByHashes method was called
	db.AssertExpectations(t)

	// Verify that the mock transaction generator's GetTransaction method was called
	tg.AssertExpectations(t)
}

func TestApp_GetTransactionsByHashes_EmptyHashes(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, tg, app := Setup(t)

	// Set up test data
	transactionHashes := []string{}

	// Call the GetTransactionsByHashes method
	transactions, err := app.GetTransactionsByHashes(transactionHashes)
	assert.Error(t, err)
	assert.Nil(t, transactions)

	// Verify that the mock database's GetTransactionsByHashes method was not called
	db.AssertNotCalled(t, "GetTransactionsByHashes", transactionHashes)

	// Verify that the mock transaction generator's GetTransaction method was not called
	tg.AssertNotCalled(t, "GetTransaction")
}

func TestApp_GetTransactionsByHashes_GetTransactionsByHashesError(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, tg, app := Setup(t)

	// Set up test data
	transactionHashes := []string{"hash1"}

	// Mock the database's GetTransactionsByHashes method
	expected := assert.AnError
	db.EXPECT().GetTransactionsByHashes(transactionHashes).Return(nil, expected)
	tg.EXPECT().GetTransaction("hash1").Return(tx1, nil)
	db.EXPECT().SaveTransaction(tx1).Return(nil)

	// Call the GetTransactionsByHashes method
	transactions, err := app.GetTransactionsByHashes(transactionHashes)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	// Verify that the mock database's GetTransactionsByHashes method was called
	db.AssertExpectations(t)

	// Verify that the mock transaction generator's GetTransaction method was not called
	tg.AssertExpectations(t)
}

func TestApp_CheckUserCredentials(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)

	// Set up test data
	username := "user1"
	password := "password1"

	// Mock the database's GetUserByUsername method
	user := &models.User{
		Username: username,
		Password: "$2a$10$RmIhxSs.xMrqUT0xU4v/wuAdH97Kmb.l50AcSQVfkg/nPVyGxp1cu",
	}
	db.EXPECT().GetUserByUsername(username).Return(user, nil)
	// Call the CheckUserCredentials method
	u, err := app.CheckUserCredentials(username, password)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, username, u.Username)
	assert.Equal(t, "$2a$10$RmIhxSs.xMrqUT0xU4v/wuAdH97Kmb.l50AcSQVfkg/nPVyGxp1cu", u.Password)

	// Verify that the mock database's GetUserByUsername method was called
	db.AssertExpectations(t)
}

func TestApp_CheckUserCredentials_GetUserByUsernameError(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)

	// Set up test data
	username := "user1"
	password := "password1"

	// Mock the database's GetUserByUsername method
	expected := assert.AnError
	db.EXPECT().GetUserByUsername(username).Return(nil, expected)

	// Call the CheckUserCredentials method
	u, err := app.CheckUserCredentials(username, password)
	assert.Error(t, err)
	assert.Nil(t, u)

	// Verify that the mock database's GetUserByUsername method was called
	db.AssertExpectations(t)
}

func TestApp_GetUserTransactions(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)
	// Set up test data
	userID := "user1"

	// Mock the database's GetUserTransactions method
	db.EXPECT().GetUserTransactions(userID).Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)

	// Call the GetUserTransactions method
	transactions, err := app.GetUserTransactions(userID)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "hash1", transactions[0].TxHash)
	assert.Equal(t, "hash2", transactions[1].TxHash)

	// Verify that the mock database's GetUserTransactions method was called
	db.AssertExpectations(t)
}

func TestApp_GetUserTransactions_EmptyUserID(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)

	// Set up test data
	userID := ""

	// Call the GetUserTransactions method
	transactions, err := app.GetUserTransactions(userID)
	assert.Error(t, err)
	assert.Nil(t, transactions)

	// Verify that the mock database's GetUserTransactions method was not called
	db.AssertNotCalled(t, "GetUserTransactions", userID)
}

func TestApp_GetAllTransactions(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)

	// Mock the database's GetAllTransactions method
	db.EXPECT().GetAllTransactions().Return([]*models.Transaction{
		tx1,
		tx2,
	}, nil)

	// Call the GetAllTransactions method
	transactions, err := app.GetAllTransactions()
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "hash1", transactions[0].TxHash)
	assert.Equal(t, "hash2", transactions[1].TxHash)

	// Verify that the mock database's GetAllTransactions method was called
	db.AssertExpectations(t)
}

func TestApp_AddUserTransactions(t *testing.T) {
	// Create mock instances of the database and transaction generator
	db, _, app := Setup(t)

	// Set up test data
	userID := "user1"
	transactions := []*models.Transaction{
		tx1,
		tx2,
	}

	// Mock the database's AddUserTransactions method
	db.EXPECT().AddUserTransactions(userID, transactions).Return(nil)

	// Call the AddUserTransactions method
	err := app.AddUserTransactions(userID, transactions)
	assert.NoError(t, err)

	// Verify that the mock database's AddUserTransactions method was called
	db.AssertExpectations(t)
}
