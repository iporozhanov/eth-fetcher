package models

import (
	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Transaction struct {
	TxHash          string      `gorm:"primaryKey" json:"transactionHash"`
	TxStatus        int         `json:"transactionStatus"`
	BlockHash       string      `json:"blockHash"`
	BlockNumber     int64       `json:"blockNumber"`
	From            string      `json:"from"`
	To              null.String `json:"to"`
	ContractAddress null.String `json:"contractAddress"`
	LogsCount       int         `json:"logsCount"`
	Input           string      `json:"input"`
	Value           int64       `json:"value"`
}

type User struct {
	ID                 string `gorm:"primaryKey"`
	Username           string `gorm:"unique"`
	Password           string
	ViewedTransactions []Transaction `gorm:"many2many:user_viewed_transactions;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = ksuid.New().String()
	return
}
