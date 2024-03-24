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
	"fmt"
	"github.com/uxuycom/indexer/client/xycommon"
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

func (o *OrdinalsClient) GetTransactionByTxId(ctx context.Context, txId string) (tx xycommon.RpcOrdTxResponse, err error) {

	path := fmt.Sprintf("tx/%s", txId)
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	tx = xycommon.RpcOrdTxResponse{}
	err = o.client.CallContext(ctx, "GET", apiUrl, &tx)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetTransactionByTxId func error,err=%v", err)
		return tx, err
	}
	return tx, nil
}

func (o *OrdinalsClient) GetOutput(ctx context.Context, output string) (out xycommon.RpcOrdOutputResponse, err error) {

	path := fmt.Sprintf("output/%s", output)
	apiUrl := fmt.Sprintf("%s/%s", o.endpoint, strings.TrimLeft(path, "/"))

	out = xycommon.RpcOrdOutputResponse{}
	err = o.client.CallContext(ctx, "GET", apiUrl, &out)
	if err != nil {
		xylog.Logger.Errorf("OrdinalsClient GetOutput func error,err=%v", err)
		return out, err
	}
	return out, nil
}
