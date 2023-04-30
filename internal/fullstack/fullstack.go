package fullstack

import (
	"encoding/json"
	"net/http"
	"strings"
)

// fullstack.cash
var apiURL = "https://api.fullstack.cash/v5/electrumx"
var transactionsMethod = "/transactions/"
var transactionDetailsMethod = "/tx/data/"

type TransactionsResponse struct {
	Success      bool `json:"success"`
	Transactions []struct {
		Height  int    `json:"height"`
		Tx_Hash string `json:"tx_hash"`
	}
}

type TransactionsDetailsResponse struct {
	Success      bool `json:"success"`
	Transactions []struct {
		Details struct {
			Vout []struct {
				Value float64 `json:"value"`
			}
		}
		TxId string `json:"txid"`
	}
}

func GetTxsDetailsResponse(txHashes []string) (*TransactionsDetailsResponse, error) {
	txsDetailsResp := &TransactionsDetailsResponse{}
	txs := strings.Join(txHashes, `","`)
	payloadTxt := `{ "txids" : ["` + txs + `"], "verbose": false }`
	payload := strings.NewReader(payloadTxt)
	reqTxDet, _ := http.NewRequest("POST", apiURL+transactionDetailsMethod, payload)
	reqTxDet.Header.Set("Content-Type", "application/json")
	respTxDet, _ := http.DefaultClient.Do(reqTxDet)
	if err := json.NewDecoder(respTxDet.Body).Decode(txsDetailsResp); err != nil {
		return nil, err
	}
	return txsDetailsResp, nil
}

func GetTXs(bchAddress string) ([]string, error) {
	var txsWallet []string
	res, err := http.Get(apiURL + transactionsMethod + bchAddress)
	if err != nil {
		return nil, err
	}
	txResp := &TransactionsResponse{}
	if err := json.NewDecoder(res.Body).Decode(txResp); err != nil {
		return nil, err
	}

	for _, tx := range txResp.Transactions {
		txsWallet = append(txsWallet, tx.Tx_Hash)
	}
	return txsWallet, nil
}
