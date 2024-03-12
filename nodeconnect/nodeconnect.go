package nodeconnect

import (
	"context"
	"encoding/hex"
	"eth-fetcher/database/models"
	"fmt"

	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v4"
)

type Node struct {
	Client Client
	log    *zap.SugaredLogger
}

type Client interface {
	ethereum.TransactionReader
	Close()
}

const zeroAddress = "0x0000000000000000000000000000000000000000"

func (n *Node) GetTransaction(txID string) (*models.Transaction, error) {
	txHash := common.HexToHash(txID)
	tx := &models.Transaction{
		TxHash: txID,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	ec := make(chan error, 2)
	m := new(sync.Mutex)

	go n.getTxReceiptData(txHash, tx, &wg, ec, m)

	go n.getTxData(txHash, tx, &wg, ec, m)

	wg.Wait()
	close(ec)

	for err := range ec {
		return nil, fmt.Errorf("error getting transaction  %s:%w", txID, err)
	}

	return tx, nil

}

func (n *Node) getTxReceiptData(txHash common.Hash, tx *models.Transaction, wg *sync.WaitGroup, ec chan error, m *sync.Mutex) {
	defer wg.Done()

	txReceipt, err := n.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		ec <- err
		return
	}

	m.Lock()
	tx.BlockHash = txReceipt.BlockHash.Hex()
	tx.TxStatus = int(txReceipt.Status)
	tx.BlockNumber = txReceipt.BlockNumber.Int64()
	txReceipt.ContractAddress = common.BytesToAddress(common.TrimLeftZeroes(txReceipt.ContractAddress[:]))
	if contractAddress := txReceipt.ContractAddress.Hex(); contractAddress != "" && contractAddress != zeroAddress {
		tx.ContractAddress = null.NewString(contractAddress, true)
	}
	tx.LogsCount = len(txReceipt.Logs)
	m.Unlock()
}

func (n *Node) getTxData(txHash common.Hash, tx *models.Transaction, wg *sync.WaitGroup, ec chan error, m *sync.Mutex) {
	defer wg.Done()

	t, _, err := n.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		ec <- err
		return
	}

	from, err := types.Sender(types.LatestSignerForChainID(t.ChainId()), t)
	if err != nil {
		ec <- err
		return
	}

	m.Lock()
	tx.From = from.Hex()
	if t.To() != nil {
		tx.To = null.NewString(t.To().Hex(), true)
	}
	sDec := make([]byte, hex.EncodedLen(len(t.Data())))
	_ = hex.Encode(sDec, t.Data())
	tx.Input = fmt.Sprintf("0x%s", string(sDec))
	tx.Value = t.Value().Int64()
	m.Unlock()
}

func (n *Node) Close() {
	n.Client.Close()
}

func NewNode(url string, log *zap.SugaredLogger) *Node {
	client, err := ethclient.Dial(url)
	network, _ := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Oops! There was a problem:%v", err)
	} else {
		log.Infof("Success! you are connected to networkID %d", network)
	}

	return &Node{
		client,
		log,
	}
}
