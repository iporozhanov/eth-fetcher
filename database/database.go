package database

import (
	"eth-fetcher/database/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

// NewClient creates a new database client with the provided DSN.
func NewClient(dsn string) (*Client, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Transaction{}, &models.User{})
	c := &Client{db}
	c.insertDefaultUsers()

	return c, nil
}

func NewTestClient() (*Client, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return &Client{gormDB}, mock, nil
}

func (c *Client) SaveTransaction(transaction *models.Transaction) error {
	return c.db.Create(transaction).Error
}

func (c *Client) GetTransactionsByHashes(hashes []string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := c.db.Where("tx_hash IN ?", hashes).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *Client) GetAllTransactions() ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := c.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *Client) AddUserTransactions(userID string, transactions []*models.Transaction) error {
	return c.db.Model(&models.User{ID: userID}).Association("ViewedTransactions").Append(transactions)
}

func (c *Client) GetUserTransactions(userID string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := c.db.Model(&models.User{ID: userID}).Association("ViewedTransactions").Find(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := c.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

var defaultUserNames = []string{"alice", "bob", "carol", "dave"}

func (c *Client) insertDefaultUsers() {
	for _, name := range defaultUserNames {
		u := &models.User{}
		result := c.db.Find(u, "username = ?", name)
		if result.RowsAffected > 0 {
			continue
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
		if err != nil {
			log.Error("generating password for %s:%w", name, err)

		}
		user := &models.User{Username: name, Password: string(pass)}
		err = c.db.Create(user).Error
		if err != nil {
			log.Error("crating default user %s:%w", name, err)
		}
	}
}
