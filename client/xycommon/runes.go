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
	"math/big"
)

type RpcOrdTxResponse struct {
	Chain            string    `json:"chain"`
	Etching          string    `json:"etching"`
	InscriptionCount int64     `json:"inscription_count"`
	Transaction      *RpcOrdTx `json:"transaction"`
	TxId             string    `json:"txid"`
}

type RpcOrdOutputResponse struct {
	Address      string   `json:"address"`
	Indexed      bool     `json:"indexed"`
	Inscriptions []string `json:"inscriptions"`
	Runes        []string `json:"runes"`
	SatRanges    string   `json:"sat_ranges"`
	PubKeyScript string   `json:"script_pubkey"`
	Spent        int64    `json:"spent"`
	TxId         string   `json:"transaction"`
	Value        int64    `json:"value"`
}

type RpcOrdRuneResponse struct {
	Entries [][]interface{} `json:"entries"`
}

type RpcOrdRunes struct {
	BlockHeight  int64      `json:"block_height"`
	Index        int64      `json:"index"`
	Burned       *big.Int   `json:"burned"`
	Divisibility int        `json:"divisibility"`
	Etching      string     `json:"etching"`
	Mint         *MintEntry `json:"mint"`
	Mints        int64      `json:"mints"`
	Number       int64      `json:"number"`
	Rune         string     `json:"rune"`
	Spacers      int64      `json:"spacers"`
	Supply       *big.Int   `json:"supply"`
	Symbol       string     `json:"symbol"`
	Timestamp    int64      `json:"timestamp"`
}

type MintEntry struct {
	Deadline int64 `json:"deadline"`
	End      int64 `json:"end"`
	Limit    int64 `json:"limit"`
}

type OrdTx struct {
	TxId   string           `json:"txid"`
	Events []*OrdBlockEvent `json:"events"`
}

type OrdBlockEvent struct {
	Type   string `json:"type"`
	Tick   string `json:"tick"`
	RuneId string `json:"rune_id"`
	Amount string `json:"amount"`
	From   string `json:"from"`
	To     string `json:"to"`
	Valid  bool   `json:"valid"`
	Msg    string `json:"msg"`
}

type RpcOrdTx struct {
	Version  int32    `json:"version"`
	TxIn     []*TxIn  `json:"input"`
	TxOut    []*TxOut `json:"output"`
	LockTime uint32   `json:"lock_time"`
}

type TxIn struct {
	PreviousOutPoint string   `json:"previous_output"`
	SignatureScript  string   `json:"script_sig"`
	Witness          []string `json:"witness"`
	Sequence         uint32   `json:"sequence"`
}

type TxOut struct {
	Value        int64  `json:"value"`
	PubKeyScript string `json:"script_pubkey"`
}
