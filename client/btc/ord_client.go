package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uxuycom/indexer/client/xycommon"
	"github.com/uxuycom/indexer/xylog"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

type OrdClient struct {
	endpoint     string
	client       *http.Client
	blockTimeMap *sync.Map
}

func NewOrdClient(endpoint string) *OrdClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	return &OrdClient{
		endpoint:     strings.TrimRight(strings.TrimSpace(endpoint), "/"),
		client:       client,
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
	Version    string     `json:""version`
	CommitHash string     `json:"commitHash"`
	BuildTime  string     `json:"buildTime"`
	ChainInfo  *ChainInfo `json:"chainInfo"`
}

type ChainInfo struct {
	Network     string      `json:"network"`
	OrdHeight   int64       `json:"ordHeight"`
	ChainHeight interface{} `json:"chainHeight"`
}

func (c *OrdClient) doCallContext(ctx context.Context, path string, out interface{}) error {
	startTs := time.Now()
	defer func() {
		xylog.Logger.Debugf("call ord api[%s] cost[%v]", path, time.Since(startTs))
	}()

	// check out whether is a pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("out should be a pointer")
	}

	uri := fmt.Sprintf("%s/%s", c.endpoint, strings.TrimLeft(path, "/"))
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		xylog.Logger.Debugf("111-call ord api[%s] data:[%s] err[%s]", path, "====", err)
		return fmt.Errorf("error creating request: %v", err)
	}

	// set headers
	req.Header.Set("Accept", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	// check if out is a []byte
	if reflect.TypeOf(out).Elem().Kind() == reflect.Slice {
		if reflect.TypeOf(out).Elem().Elem().Kind() == reflect.Uint8 {
			reflect.ValueOf(out).Elem().SetBytes(data)
			return nil
		}
	}

	// check if out is a string
	if reflect.TypeOf(out).Elem().Kind() == reflect.String {
		reflect.ValueOf(out).Elem().SetString(string(data))
		return nil
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("error parsing response body[%s], err[%v]", string(data), err)
	}

	return nil
}

func (c *OrdClient) callContext(ctx context.Context, path string, out interface{}) (err error) {
	ts := time.Millisecond * 100
	for retry := 0; retry < 5; retry++ {
		err = c.doCallContext(ctx, path, out)
		if err == nil {
			return nil
		}
		<-time.After(ts * time.Duration(retry))
	}
	return err
}

func (c *OrdClient) BlockNumber(ctx context.Context) (number int64, err error) {
	path := fmt.Sprintf("api/v1/node/info")
	result := &NodeInfoResponse{}
	err = c.callContext(ctx, path, &result)
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
	result := &BlockEventResponse{}
	err = c.callContext(ctx, path, &result)
	if err == nil && result.Data != nil {
		ret = result.Data
		ret.Hash = blockHash
	}
	return
}

func (c *OrdClient) GetInscription(ctx context.Context, inscriptionId string) (ins *xycommon.RpcOkxInscription, err error) {
	path := fmt.Sprintf("api/v1/ord/id/%s/inscription", inscriptionId)
	result := &InscriptionResponse{}
	err = c.callContext(ctx, path, &result)
	if err == nil {
		ins = result.Data
	}
	return
}

func (c *OrdClient) GetAddressBalanceByTick(ctx context.Context, address, tick string) (balance *xycommon.RpcOkxBalance, err error) {
	path := fmt.Sprintf("api/v1/brc20/tick/%s/address/%s/balance", url.QueryEscape(tick), address)
	result := &AddressBalanceResponse{}
	err = c.callContext(ctx, path, &result)
	if err == nil && result.Data != nil {
		balance = result.Data
	}
	return
}
