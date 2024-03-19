// Copyright (c) 2023-2024 The UXUY Developer Team
// License:
// MIT License

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE

package xycommon

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/shopspring/decimal"
)

type OkxAddress struct {
	Address string `json:"address"`
}

type BlockEvent struct {
	Type              string      `json:"type"`
	Tick              string      `json:"tick"`
	InscriptionId     string      `json:"inscriptionId"`
	InscriptionNumber int64       `json:"inscriptionNumber"`
	OldSatpoint       string      `json:"oldSatpoint"`
	NewSatpoint       string      `json:"newSatpoint"`
	Amount            string      `json:"amount"`
	From              *OkxAddress `json:"from"`
	To                *OkxAddress `json:"to"`
	Valid             bool        `json:"valid"`
	Msg               string      `json:"msg"`
}

type Tx struct {
	Txid   string        `json:"txid"`
	Events []*BlockEvent `json:"events"`
}

type RpcOkxBlockResponse struct {
	Txs    []*Tx                 `json:"block"`
	Height int64                 `json:"height"`
	Time   uint64                `json:"time"`
	Hash   string                `json:"hash"`
	RpcTxs []btcjson.TxRawResult `json:"rpc_txs"`
}

type RpcOkxInscription struct {
	Id            string        `json:"id"`
	Number        int64         `json:"number"`
	Content       string        `json:"content"`
	ContentType   string        `json:"contentType"`
	Owner         OkxAddress    `json:"owner"`
	GenesisHeight int64         `json:"genesisHeight"`
	Location      string        `json:"location"`
	Collections   []interface{} `json:"collections"`
	Sat           interface{}   `json:"sat"`
}

type Inscription struct {
	Id           string          `json:"id"`
	Tick         string          `json:"tick"`
	LimitPerMint decimal.Decimal `json:"limit_per_mint"`
	TotalSupply  decimal.Decimal `json:"total_supply"`
	Decimals     int8            `json:"decimals"`
	Owner        string          `json:"owner"`
	Number       int64           `json:"number"`
}

type InscriptionContent struct {
	Protocol     string `json:"p"`
	Operation    string `json:"op"`
	Tick         string `json:"tick"`
	Max          string `json:"max"`
	LimitPerMint string `json:"lim"`
}

type RpcOkxBalance struct {
	Tick                string `json:"tick"`
	AvailableBalance    string `json:"availableBalance"`
	TransferableBalance string `json:"transferableBalance"`
	OverallBalance      string `json:"overallBalance"`
}

type AddressBalance struct {
	Address          string          `json:"address"`
	Tick             string          `json:"tick"`
	AvailableBalance decimal.Decimal `json:"available_balance"`
	OverallBalance   decimal.Decimal `json:"overall_balance"`
}
