package database_test

import (
	"eth-fetcher/database"
	"eth-fetcher/database/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func TestClient_SaveTransaction(t *testing.T) {
	client, mock, err := database.NewTestClient()
	assert.NoError(t, err)
	defer client.Close()

	transaction := &models.Transaction{

		TxHash:          "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2",
		TxStatus:        1,
		BlockHash:       "0x92557f7e29c39cae6be013ffc817620fcd5233b68405cdfc6e0b5528261e81e5",
		BlockNumber:     7976373,
		From:            "0xF29A6c0f8eE500dC87d0d4EB8B26a6faC7A76767",
		To:              null.NewString("0xb0428bF0D49eB5c2239A815B43E59E124b84E303", true),
		ContractAddress: null.NewString("", false),
		LogsCount:       1,
		Input:           "0x",
		Value:           50000000000000000,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"transactions\" (.+)").
		WithArgs(transaction.TxHash,
			transaction.TxStatus,
			transaction.BlockHash,
			transaction.BlockNumber,
			transaction.From,
			transaction.To,
			transaction.ContractAddress,
			transaction.LogsCount,
			transaction.Input,
			transaction.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = client.SaveTransaction(transaction)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_GetTransactionsByHashes(t *testing.T) {
	client, mock, err := database.NewTestClient()
	assert.NoError(t, err)
	defer client.Close()

	hashes := []string{"0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", "0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057"}

	rows := sqlmock.NewRows([]string{
		"transactionHash",
		"transactionStatus",
		"blockHash",
		"blockNumber",
		"from",
		"to",
		"contractAddress",
		"logsCount",
		"input",
		"value"}).
		AddRow(
			"0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2",
			1,
			"0x92557f7e29c39cae6be013ffc817620fcd5233b68405cdfc6e0b5528261e81e5",
			7976373,
			"0xF29A6c0f8eE500dC87d0d4EB8B26a6faC7A76767",
			null.NewString("0xb0428bF0D49eB5c2239A815B43E59E124b84E303", true),
			null.NewString("0x0000000000000000000000000000000000000000", true),
			1,
			"0x",
			50000000000000000).
		AddRow(
			"0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057",
			1,
			"0x32edca7a39d0b1fc3d19fd1487c3c69beadad7cdcd5e5f1c9e815e7d1c460a0d",
			7957369,
			"0x22Ba753CA065d65D4d0b9f4FAC7a669746175199",
			null.NewString("0x14cB06e8dE2222912138F9a062E5a4d9F4821409", true),
			null.NewString("0x0000000000000000000000000000000000000000", true),
			3,
			"0x95297e24a8d61b73377cdc07fcf0cdd5473a1c81d541d3bcbbac29dd02d9f680af901d705591dba920dde33e5f1043c1b84106b0f223e7b954b17bde9ffe62206b583b2d00000000000000000000000000000000000000000000000000000000000021d300eea48906338871c59f0d12348b85c66461cd8c9e80faa4d3e63b134279595a159d3bfa16088686ea5b2406f82109b60a5792b77bc173106c01f0fbfed6598905d63c4ca36d0740e427d53ea8d1cc707b15a846c854a14b2e3b2e30ce129b8721a54e650d0e077cf260d8c3c84a431b287bd35ffe4c03c27a19d9a0d3320ae905a76cdd5a8bfeffa1c837279d67654e053a8e80cf2e581968a93bf827c3cf702d8c881054165ebd6d1ebea052f9af3a10338c9314ed99609735b8b76fe274c411d32840d8a1f85b51ee84bd2b0d70fe5725362406ac200a1186ea82ae39731a05d84408b5eca5130fa799aa898bbb2132054dcd8890ff004ac855f57c813fc6",
			0)

	mock.ExpectQuery("SELECT (.+) FROM \"transactions\" WHERE tx_hash IN (.+)").
		WithArgs(hashes[0], hashes[1]).
		WillReturnRows(rows)

	transactions, err := client.GetTransactionsByHashes(hashes)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_GetAllTransactions(t *testing.T) {
	client, mock, err := database.NewTestClient()
	assert.NoError(t, err)
	defer client.Close()

	rows := sqlmock.NewRows([]string{
		"transactionHash",
		"transactionStatus",
		"blockHash",
		"blockNumber",
		"from",
		"to",
		"contractAddress",
		"logsCount",
		"input",
		"value"}).
		AddRow(
			"0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2",
			1,
			"0x92557f7e29c39cae6be013ffc817620fcd5233b68405cdfc6e0b5528261e81e5",
			7976373,
			"0xF29A6c0f8eE500dC87d0d4EB8B26a6faC7A76767",
			null.NewString("0xb0428bF0D49eB5c2239A815B43E59E124b84E303", true),
			null.NewString("0x0000000000000000000000000000000000000000", true),
			1,
			"0x",
			50000000000000000).
		AddRow(
			"0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057",
			1,
			"0x32edca7a39d0b1fc3d19fd1487c3c69beadad7cdcd5e5f1c9e815e7d1c460a0d",
			7957369,
			"0x22Ba753CA065d65D4d0b9f4FAC7a669746175199",
			null.NewString("0x14cB06e8dE2222912138F9a062E5a4d9F4821409", true),
			null.NewString("0x0000000000000000000000000000000000000000", true),
			3,
			"0x95297e24a8d61b73377cdc07fcf0cdd5473a1c81d541d3bcbbac29dd02d9f680af901d705591dba920dde33e5f1043c1b84106b0f223e7b954b17bde9ffe62206b583b2d00000000000000000000000000000000000000000000000000000000000021d300eea48906338871c59f0d12348b85c66461cd8c9e80faa4d3e63b134279595a159d3bfa16088686ea5b2406f82109b60a5792b77bc173106c01f0fbfed6598905d63c4ca36d0740e427d53ea8d1cc707b15a846c854a14b2e3b2e30ce129b8721a54e650d0e077cf260d8c3c84a431b287bd35ffe4c03c27a19d9a0d3320ae905a76cdd5a8bfeffa1c837279d67654e053a8e80cf2e581968a93bf827c3cf702d8c881054165ebd6d1ebea052f9af3a10338c9314ed99609735b8b76fe274c411d32840d8a1f85b51ee84bd2b0d70fe5725362406ac200a1186ea82ae39731a05d84408b5eca5130fa799aa898bbb2132054dcd8890ff004ac855f57c813fc6",
			0)

	mock.ExpectQuery("SELECT (.+) FROM \"transactions\" ").
		WillReturnRows(rows)

	transactions, err := client.GetAllTransactions()
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_GetUserByUsername(t *testing.T) {
	client, mock, err := database.NewTestClient()
	assert.NoError(t, err)
	defer client.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"username",
		"password",
	}).
		AddRow(
			"1",
			"alice",
			"$2a$10$")

	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE username = (.+)").
		WithArgs("alice").
		WillReturnRows(rows)
	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE username = (.+)").
		WithArgs("bob").
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := client.GetUserByUsername("alice")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = client.GetUserByUsername("bob")

	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClient_AddUserTransactions(t *testing.T) {
	client, mock, err := database.NewTestClient()
	assert.NoError(t, err)
	defer client.Close()

	transactions := []*models.Transaction{
		{
			TxHash:          "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2",
			TxStatus:        1,
			BlockHash:       "0x92557f7e29c39cae6be013ffc817620fcd5233b68405cdfc6e0b5528261e81e5",
			BlockNumber:     7976373,
			From:            "0xF29A6c0f8eE500dC87d0d4EB8B26a6faC7A76767",
			To:              null.NewString("0xb0428bF0D49eB5c2239A815B43E59E124b84E303", true),
			ContractAddress: null.NewString("", false),
			LogsCount:       1,
			Input:           "0x",
			Value:           50000000000000000,
		},
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"transactions\" (.+)").
		WithArgs(transactions[0].TxHash,
			transactions[0].TxStatus,
			transactions[0].BlockHash,
			transactions[0].BlockNumber,
			transactions[0].From,
			transactions[0].To,
			transactions[0].ContractAddress,
			transactions[0].LogsCount,
			transactions[0].Input,
			transactions[0].Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO \"user_viewed_transactions\" (.+)").
		WithArgs("1", transactions[0].TxHash).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = client.AddUserTransactions("1", transactions)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
