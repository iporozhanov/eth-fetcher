package nodeconnect_test

import (
	"encoding/json"
	node "eth-fetcher/nodeconnect"
	"eth-fetcher/nodeconnect/mocks"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getTransactionResponse() *types.Transaction {
	s := `{
		"accessList":[],
		"blockHash":"0x0155db99111f10086bad292d3bd0be9472aff9cf0f33d7d35f2db4814ffad0f6",
		"blockNumber":"0x112418d",
		"chainId":"0x5",
		"from":"0xe2a467bfe1e1bedcdf1343d3a45f60c50e988696",
		"gas":"0x3c546",
		"gasPrice":"0x20706def53",
		"hash":"0xce0aadd04968e21f569167570011abc8bc17de49d4ae3aed9476de9e03facff9",
		"input":"0xb6f9de9500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000e2a467bfe1e1bedcdf1343d3a45f60c50e9886960000000000000000000000000000000000000000000000000000000064e54a3b0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000de15b9919539113a1930d3eed5088cd10338abb5",
		"maxFeePerGas":"0x22b05d8efd",
		"maxPriorityFeePerGas":"0x1bf08eb000",
		"nonce":"0x12c",
		"r":"0xa07fd6c16e169f0e54b394235b3a8201101bb9d0eba9c8ae52dbdf556a363388",
		"s":"0x36f5da9310b87fefbe9260c3c05ec6cbefc426f1ff3b3a41ea21b5533a787dfc",
		"to":"0x7a250d5630b4cf539739df2c5dacb4c659f2488d",
		"transactionIndex":"0x0",
		"type":"0x2",
		"v":"0x1",
		"value":"0x2c68af0bb140000",
		"yParity":"0x1"
	 }`
	txn := &types.Transaction{}
	json.Unmarshal([]byte(s), txn)
	return txn
}

func getTransactionReceiptResponse() *types.Receipt {
	jsr := `{
		"blockHash":"0x0155db99111f10086bad292d3bd0be9472aff9cf0f33d7d35f2db4814ffad0f6",
		"blockNumber":"0x112418d",
		"contractAddress": 	null,
		"cumulativeGasUsed":"0xc5f3e7",
		"effectiveGasPrice":"0xa45b9a444",
		"from":"0xe2a467bfe1e1bedcdf1343d3a45f60c50e988696",
		"gasUsed":"0x565f",
		"logs": [
		{
			"address":"0x388c818ca8b9251b393131c08a736a67ccb19297",
			"blockHash":"0x0155db99111f10086bad292d3bd0be9472aff9cf0f33d7d35f2db4814ffad0f6",
			"blockNumber":"0x112418d",
			"data":"0x00000000000000000000000000000000000000000000000011b6b79503fb875d",
			"logIndex":"0x187",
			"removed":false,
			"topics": [
			"0x27f12abfe35860a9a927b465bb3d4a9c23c8428174b83f278fe45ed7b4da2662"
			],
			"transactionHash":"0x7114b4da1a6ed391d5d781447ed443733dcf2b508c515b81c17379dea8a3c9af",
			"transactionIndex":"0x76"
		}
		],
		"logsBloom":"0x00000000000000000000000000000000000100004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000",
		"status":"0x1",
		"to":"0x388c818ca8b9251b393131c08a736a67ccb19297",
		"transactionHash":"0x7114b4da1a6ed391d5d781447ed443733dcf2b508c515b81c17379dea8a3c9af",
		"transactionIndex":"0x76",
		"type":"0x2"}`
	txr := &types.Receipt{}
	json.Unmarshal([]byte(jsr), txr)
	return txr
}

func TestNode_GetTransaction_Success(t *testing.T) {
	txn := getTransactionResponse()
	txr := getTransactionReceiptResponse()
	client := mocks.NewClient(t)

	client.EXPECT().TransactionByHash(mock.Anything, mock.Anything).
		Return(txn, true, nil)
	client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(txr, nil)

	node := &node.Node{Client: client}
	tx, err := node.GetTransaction("0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2")
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", tx.TxHash)
	assert.Equal(t, int(1), tx.TxStatus)
	assert.Equal(t, "0x0155db99111f10086bad292d3bd0be9472aff9cf0f33d7d35f2db4814ffad0f6", tx.BlockHash)
	assert.Equal(t, int64(17973645), tx.BlockNumber)
	assert.Equal(t, "0x425Db51efE6971d86512e892BeABA90Bc920Cdda", tx.From)
	assert.Equal(t, "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", tx.To.String)
	assert.Equal(t, false, tx.ContractAddress.Valid)
	assert.Equal(t, 1, tx.LogsCount)
	assert.Equal(t, "0xb6f9de9500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000e2a467bfe1e1bedcdf1343d3a45f60c50e9886960000000000000000000000000000000000000000000000000000000064e54a3b0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000de15b9919539113a1930d3eed5088cd10338abb5", tx.Input)
	assert.Equal(t, int64(200000000000000000), tx.Value)
}

func TestNode_GetTransaction_SuccessValidContractAddress(t *testing.T) {
	txn := getTransactionResponse()
	txr := getTransactionReceiptResponse()
	txr.ContractAddress = common.HexToAddress("0x3664F6c1178E19Bb775b597d6584CaA3B88a1C35")
	client := mocks.NewClient(t)

	client.EXPECT().TransactionByHash(mock.Anything, mock.Anything).
		Return(txn, true, nil)
	client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(txr, nil)
	node := &node.Node{Client: client}
	tx, err := node.GetTransaction("0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2")
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", tx.TxHash)
	assert.Equal(t, int(1), tx.TxStatus)
	assert.Equal(t, "0x0155db99111f10086bad292d3bd0be9472aff9cf0f33d7d35f2db4814ffad0f6", tx.BlockHash)
	assert.Equal(t, int64(17973645), tx.BlockNumber)
	assert.Equal(t, "0x425Db51efE6971d86512e892BeABA90Bc920Cdda", tx.From)
	assert.Equal(t, "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", tx.To.String)
	assert.Equal(t, "0x3664F6c1178E19Bb775b597d6584CaA3B88a1C35", tx.ContractAddress.String)
	assert.Equal(t, 1, tx.LogsCount)
	assert.Equal(t, "0xb6f9de9500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000e2a467bfe1e1bedcdf1343d3a45f60c50e9886960000000000000000000000000000000000000000000000000000000064e54a3b0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000de15b9919539113a1930d3eed5088cd10338abb5", tx.Input)
	assert.Equal(t, int64(200000000000000000), tx.Value)
}

func TestNode_GetTransaction_FailOnTransactionByHash(t *testing.T) {
	client := mocks.NewClient(t)
	txr := getTransactionReceiptResponse()
	expected := fmt.Errorf("error on by hash")
	client.EXPECT().TransactionByHash(mock.Anything, mock.Anything).
		Return(nil, false, expected)
	client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(txr, nil)

	node := &node.Node{Client: client}
	tx, err := node.GetTransaction("0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2")
	assert.Error(t, err)
	assert.ErrorContains(t, err, expected.Error())
	assert.Nil(t, tx)
}

func TestNode_GetTransaction_FailOnTransactionReceipt(t *testing.T) {
	txr := getTransactionResponse()
	client := mocks.NewClient(t)
	expected := fmt.Errorf("error on receipt")
	client.EXPECT().TransactionByHash(mock.Anything, mock.Anything).
		Return(txr, true, nil)
	client.EXPECT().TransactionReceipt(mock.Anything, mock.Anything).Return(nil, expected)

	node := &node.Node{Client: client}
	tx, err := node.GetTransaction("0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2")
	assert.Error(t, err)
	assert.ErrorContains(t, err, expected.Error())
	assert.Nil(t, tx)
}
