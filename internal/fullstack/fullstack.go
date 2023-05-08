package fullstack

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// fullstack.cash
var apiURL = "https://api.fullstack.cash/v5/electrumx"
var transactionsMethod = "/transactions/"
var transactionDetailsMethod = "/tx/data/"

var ErrInvalidAddrFormat = errors.New("fullstack: Unsupported address format")

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

func GetTxsDetailsResponse(txHashes []string, infoLog *log.Logger) (*TransactionsDetailsResponse, error) {
	txsDetailsResp := &TransactionsDetailsResponse{}
	txs := strings.Join(txHashes, `","`)
	payloadTxt := `{ "txids" : ["` + txs + `"], "verbose": false }`
	payload := strings.NewReader(payloadTxt)
	infoLog.Printf("FULLSTACK - POST %s", transactionDetailsMethod)
	reqTxDet, _ := http.NewRequest("POST", apiURL+transactionDetailsMethod, payload)
	reqTxDet.Header.Set("Content-Type", "application/json")
	respTxDet, _ := http.DefaultClient.Do(reqTxDet)
	if err := json.NewDecoder(respTxDet.Body).Decode(txsDetailsResp); err != nil {
		return nil, err
	}
	return txsDetailsResp, nil
}

type TransactionsResponse struct {
	Success      bool `json:"success"`
	Transactions []struct {
		Height  int    `json:"height"`
		Tx_Hash string `json:"tx_hash"`
	}
	Error string `json:"error"`
}

func GetTXs(bchAddress string, infoLog *log.Logger) ([]string, error) {
	var txsWallet []string
	infoLog.Printf("FULLSTACK - GET %s", transactionsMethod)
	res, err := http.Get(apiURL + transactionsMethod + bchAddress)
	if err != nil {
		return nil, err
	}
	txResp := &TransactionsResponse{}
	if err := json.NewDecoder(res.Body).Decode(txResp); err != nil {
		return nil, err
	}
	if !txResp.Success {
		if strings.Contains(txResp.Error, "Unsupported address format") {
			return nil, ErrInvalidAddrFormat
		} else {
			return nil, errors.New(txResp.Error)
		}
	}
	for _, tx := range txResp.Transactions {
		txsWallet = append(txsWallet, tx.Tx_Hash)
	}

	if len(txsWallet) > 20 {
		txsWallet = txsWallet[:20]
	}

	return txsWallet, nil
}
