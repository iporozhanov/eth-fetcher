package handlers

import "eth-fetcher/database/models"

type GetTransactionsResponse struct {
	Transactions []*models.Transaction `json:"transactions"`
}

// type Transaction struct {
// 	TxHash          string      `json:"transactionHash"`
// 	TxStatus        int         `json:"transactionStatus"`
// 	BlockHash       string      `json:"blockHash"`
// 	BlockNumber     int64       `json:"blockNumber"`
// 	From            string      `json:"from"`
// 	To              null.String `json:"to"`
// 	ContractAddress null.String `json:"contractAddress"`
// 	LogsCount       int         `json:"logsCount"`
// 	Input           string      `json:"input"`
// 	Value           int64       `json:"value"`
// }

type ErrorResponse struct {
	Msg  string `json:"error"`
	Code int    `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return e.Msg
}

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}
