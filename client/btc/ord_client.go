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
	"net/url"
	"strings"
	"sync"
	"time"
)

type OrdClient struct {
	endpoint     string
	client       *utils.HttpClient
	blockTimeMap *sync.Map
}

func NewOrdClient(endpoint string) *OrdClient {

	return &OrdClient{
		endpoint:     strings.TrimRight(strings.TrimSpace(endpoint), "/"),
		client:       utils.NewHttpClient(),
		blockTimeMap: &sync.Map{},
	}
}

type BlockInscriptions struct {
	Inscriptions []string `json:"inscriptions"`
	More         bool     `json:"more"`
	PageIndex    int      `json:"page_index"`
}

type BlockEventResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *xycommon.RpcOkxBlockResponse
}

type InscriptionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *xycommon.RpcOkxInscription
}

type AddressBalanceResponse struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *xycommon.RpcOkxBalance `json:"data"`
}

type NodeInfoResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *NodeInfo `json:"data"`
}

type NodeInfo struct {
	Version    string     `json:"version"`
	CommitHash string     `json:"commitHash"`
	BuildTime  string     `json:"buildTime"`
	ChainInfo  *ChainInfo `json:"chainInfo"`
}

type ChainInfo struct {
	Network     string      `json:"network"`
	OrdHeight   int64       `json:"ordHeight"`
	ChainHeight interface{} `json:"chainHeight"`
}

func (c *OrdClient) BlockNumber(ctx context.Context) (number int64, err error) {
	path := fmt.Sprintf("api/v1/node/info")
	result := &NodeInfoResponse{}

	apiUrl := fmt.Sprintf("%s/%s", c.endpoint, strings.TrimLeft(path, "/"))
	err = c.client.CallContext(ctx, "GET", apiUrl, &result)
	if result.Data != nil && result.Data.ChainInfo != nil {
		number = result.Data.ChainInfo.OrdHeight
	}

	blockNum, _ := c.checkBlockNumber(number)
	xylog.Logger.Infof("BlockNumber block checkBlockNumber return  number[%v]", blockNum)
	number = blockNum
	if blockNum < 0 {
		number = 0
	}

	xylog.Logger.Infof("ord block height[%d]", number)
	return number, nil
}

func (c *OrdClient) checkBlockNumber(number int64) (int64, int64) {
	bNumber, ok := c.blockTimeMap.Load("blockNumber")
	bTime, tok := c.blockTimeMap.Load("blockTime")
	xylog.Logger.Infof("BlockNumber block checkBlockNumber bNumber[%v] bTime[%v] number[%v]", bNumber, bTime, number)
	if !ok || !tok {
		c.blockTimeMap.Store("blockNumber", number)
		c.blockTimeMap.Store("blockTime", time.Now().Unix())
		return number - 1, 0
	}

	blockNum := bNumber.(int64)
	blockTime := bTime.(int64)
	if number < blockNum {
		xylog.Logger.Infof("BlockNumber block checkBlockNumber bNumber[%v] bTime[%v] number[%v] number < blockNum", bNumber, bTime, number)
		return number, time.Now().Unix()
	} else if number > blockNum {
		c.blockTimeMap.Store("blockNumber", number)
		c.blockTimeMap.Store("blockTime", time.Now().Unix())
		xylog.Logger.Infof("BlockNumber block checkBlockNumber bNumber[%v] bTime[%v] number[%v] number > blockNum", bNumber, bTime, number)
		return number - 1, 0
	} else if blockNum == number && time.Now().Unix()-blockTime < 30 {
		xylog.Logger.Infof("BlockNumber block checkBlockNumber bNumber[%v] bTime[%v] number[%v] blockNum == number && time.Now().Unix()-blockTime < 30 (number - 1)[%v]", bNumber, bTime, number, number-1)
		return number - 1, 0
	}

	return number, blockTime
}

func (c *OrdClient) BlockByHash(ctx context.Context, blockHash string) (ret *xycommon.RpcOkxBlockResponse, err error) {
	path := fmt.Sprintf("api/v1/brc20/block/%s/events", blockHash)
	apiUrl := fmt.Sprintf("%s/%s", c.endpoint, strings.TrimLeft(path, "/"))
	result := &BlockEventResponse{}
	err = c.client.CallContext(ctx, "GET", apiUrl, &result)
	if err == nil && result.Data != nil {
		ret = result.Data
		ret.Hash = blockHash
	}
	return
}

func (c *OrdClient) GetInscription(ctx context.Context, inscriptionId string) (ins *xycommon.RpcOkxInscription, err error) {
	path := fmt.Sprintf("api/v1/ord/id/%s/inscription", inscriptionId)
	apiUrl := fmt.Sprintf("%s/%s", c.endpoint, strings.TrimLeft(path, "/"))
	result := &InscriptionResponse{}
	err = c.client.CallContext(ctx, "GET", apiUrl, &result)
	if err == nil {
		ins = result.Data
	}
	return
}

func (c *OrdClient) GetAddressBalanceByTick(ctx context.Context, address, tick string) (balance *xycommon.RpcOkxBalance, err error) {
	path := fmt.Sprintf("api/v1/brc20/tick/%s/address/%s/balance", url.QueryEscape(tick), address)

	apiUrl := fmt.Sprintf("%s/%s", c.endpoint, strings.TrimLeft(path, "/"))
	result := &AddressBalanceResponse{}
	err = c.client.CallContext(ctx, "GET", apiUrl, &result)
	if err == nil && result.Data != nil {
		balance = result.Data
	}
	return
}
