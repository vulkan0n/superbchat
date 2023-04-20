package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// fullstack.cash
var apiURL = "https://api.fullstack.cash/v5/electrumx"
var transactionsMethod = "/transactions/"
var transactionDetailsMethod = "/tx/data/"

type transactionsResponse struct {
	Success      bool `json:"success"`
	Transactions []struct {
		Height  int    `json:"height"`
		Tx_Hash string `json:"tx_hash"`
	}
}

type transactionsDetailsResponse struct {
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

func getTxsDetailsResponse(txsDetailsResp *transactionsDetailsResponse, txHashes []string) {
	txs := strings.Join(txHashes, `","`)
	payloadTxt := `{ "txids" : ["` + txs + `"], "verbose": false }`
	payload := strings.NewReader(payloadTxt)
	reqTxDet, _ := http.NewRequest("POST", apiURL+transactionDetailsMethod, payload)
	reqTxDet.Header.Set("Content-Type", "application/json")
	respTxDet, _ := http.DefaultClient.Do(reqTxDet)
	if err := json.NewDecoder(respTxDet.Body).Decode(txsDetailsResp); err != nil {
		fmt.Println(err.Error())
	}
}

func getTXs(bchAddress string, txHashes *[]string) {
	res, err := http.Get(apiURL + transactionsMethod + bchAddress)
	if err == nil {
		txResp := &transactionsResponse{}
		if err := json.NewDecoder(res.Body).Decode(txResp); err != nil {
			fmt.Println(err.Error())
		}

		for _, tx := range txResp.Transactions {
			*txHashes = append(*txHashes, tx.Tx_Hash)
		}
	}
}
