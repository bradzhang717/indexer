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

package btc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uxuycom/indexer/client/xycommon"
	"github.com/uxuycom/indexer/model"
	"github.com/uxuycom/indexer/utils"
	"github.com/uxuycom/indexer/xylog"
	"strconv"
	"strings"
	"sync"
)

type OrdinalsClient struct {
	endpoint     string
	client       *utils.HttpClient
	blockTimeMap *sync.Map
}

func NewOrdinalsClient(endpoint string) *OrdinalsClient {
	return &OrdinalsClient{
		endpoint:     strings.TrimRight(strings.TrimSpace(endpoint), "/"),
		client:       utils.NewHttpClient(),
		blockTimeMap: &sync.Map{},
	}
}
func (o *OrdinalsClient) BlockNumber(ctx context.Context) (number int64, err error) {

	path := fmt.Sprintf("blockheight")
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	var blockHeight string
	err = o.client.CallContext(ctx, "GET", apiUrl, &blockHeight)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient BlockNumber func error,err=%v", err)
		return int64(0), err
	}
	return strconv.ParseInt(blockHeight, 10, 64)
}

func (o *OrdinalsClient) GetTransactionByTxId(ctx context.Context, txId string) (*xycommon.RpcOrdTxResponse, error) {

	path := fmt.Sprintf("tx/%s", txId)
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	tx := xycommon.RpcOrdTxResponse{}
	err := o.client.CallContext(ctx, "GET", apiUrl, &tx)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetTransactionByTxId func error,err=%v", err)
		return &tx, err
	}
	return &tx, nil
}

func (o *OrdinalsClient) GetOutput(ctx context.Context, output string) (*xycommon.RpcOrdOutputResponse, error) {

	path := fmt.Sprintf("output/%s", output)
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	out := xycommon.RpcOrdOutputResponse{}
	err := o.client.CallContext(ctx, "GET", apiUrl, &out)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v", err)
		return &out, err
	}
	return &out, nil
}

func (o *OrdinalsClient) GetRune(ctx context.Context, rune string) (*xycommon.RpcOrdRunes, error) {

	path := fmt.Sprintf("rune/%s", rune)
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	var rsp *xycommon.RpcOrdRuneResponse
	err := o.client.CallContext(ctx, "GET", apiUrl, &rsp)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v", err)
		return nil, err
	}
	splitId := strings.Split(rsp.Id, ":")
	blockHeight := splitId[0]
	index := splitId[1]
	if num, err := strconv.ParseInt(blockHeight, 10, 64); err == nil {
		rsp.Entry.BlockHeight = num
	}
	if idx, err := strconv.ParseInt(index, 10, 64); err == nil {
		rsp.Entry.Index = idx
	}
	return rsp.Entry, nil
}

func (o *OrdinalsClient) GetRunes(ctx context.Context) ([]xycommon.RpcOrdRunes, error) {

	path := fmt.Sprintf("runes")
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	var rsp *xycommon.RpcOrdRunesResponse
	err := o.client.CallContext(ctx, "GET", apiUrl, &rsp)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v", err)
		return nil, err
	}

	allRunes := make([]xycommon.RpcOrdRunes, 0)
	for _, entry := range rsp.Entries {
		runId := entry[0].(string)
		if len(runId) < 0 {
			continue
		}

		splitId := strings.Split(runId, ":")
		blockHeight := splitId[0]
		index := splitId[1]

		entryData, _ := json.Marshal(entry[1])
		var runes xycommon.RpcOrdRunes

		if err := json.Unmarshal(entryData, &runes); err != nil {
			xylog.Logger.Errorf("Error parsing entry detail=%v", err)
			continue
		}

		if num, err := strconv.ParseInt(blockHeight, 10, 64); err == nil {
			runes.BlockHeight = num
		}
		if idx, err := strconv.ParseInt(index, 10, 64); err == nil {
			runes.Index = idx
		}
		allRunes = append(allRunes, runes)
	}
	return allRunes, nil
}

func (o *OrdinalsClient) SaveRunesToDb(ctx context.Context) error {

	runes, err := o.GetRunes(ctx)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetRunes func error,err=%v", err)
		return err
	}
	dbRunes := make([]model.Runes, 0)
	for _, r := range runes {

		ru := model.Runes{
			BlockHeight:  r.BlockHeight,
			Index:        r.Index,
			Burned:       r.Burned,
			Divisibility: r.Divisibility,
			Etching:      r.Etching,
			Mints:        r.Mints,
			Number:       r.Number,
			Rune:         r.Rune,
			Spacers:      r.Spacers,
			Supply:       r.Supply,
			Symbol:       r.Symbol,
			Timestamp:    r.Timestamp,
		}
		if r.Mint != nil {
			ru.Deadline = r.Mint.Deadline
			ru.End = r.Mint.End
			ru.Limit = r.Mint.Limit
		}
		dbRunes = append(dbRunes, ru)

	}
	//dbClient.SaveRunes(dbRunes) // TODO
	return nil
}

func (o *OrdinalsClient) GetRunesBalances(ctx context.Context) (*xycommon.RpcOrdRuneResponse, error) {

	path := fmt.Sprintf("runes/balances")
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	rsp := xycommon.RpcOrdRuneResponse{}
	err := o.client.CallContext(ctx, "GET", apiUrl, &rsp)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v", err)
		return &rsp, err
	}
	return &rsp, nil
}

func (o *OrdinalsClient) GetTxEvents(ctx context.Context, txIds []string) ([]*xycommon.OrdTx, error) {

	if len(txIds) <= 0 {
		return nil, errors.New("txIds is empty")
	}

	txs := make([]*xycommon.OrdTx, 0)
	for _, txId := range txIds {

		txRsp, err := o.GetTransactionByTxId(ctx, txId)
		if err != nil {
			xylog.Logger.Errorf("OrdinalsClient GetTxEvents func error,err=%v, txId =%v", err, txId)
			continue
		}
		if txRsp == nil {
			continue
		}

		previousOutput := txRsp.Transaction.TxIn[0].PreviousOutPoint
		outFrom, err := o.GetOutput(ctx, previousOutput)
		if err != nil {
			xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v, previousOutput =%v", err, previousOutput)
			continue
		}
		addressFrom := outFrom.Address

		outTo, err := o.GetOutput(ctx, txId+":0")
		if err != nil {
			xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v, output =%v", err, txId+":0")
		}
		runes := outTo.Runes

		events := make([]*xycommon.OrdBlockEvent, 0)
		if len(runes) > 0 {
			for _, runeId := range runes {

				event := &xycommon.OrdBlockEvent{
					RuneId: runeId,
					From:   addressFrom,
					To:     outTo.Address,
					Type:   "",    // TODO
					Tick:   "",    // TODO
					Amount: "",    // TODO
					Valid:  false, // TODO
					Msg:    "",    // TODO
				}
				events = append(events, event)
			}
		}
		tx := &xycommon.OrdTx{
			TxId:   txId,
			Events: events,
		}
		txs = append(txs, tx)
	}

	return txs, nil
}
