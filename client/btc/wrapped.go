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
	"github.com/ethereum/go-ethereum"
	"github.com/uxuycom/indexer/client/xycommon"
	"github.com/uxuycom/indexer/config"
	"github.com/uxuycom/indexer/xylog"
	"math/big"
	"time"
)

type BClient struct {
	btcClient *BtcClient
}

// Dial connects a client to the given chain config.
func Dial(chainCfg *config.ChainConfig) (*BClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := NewBtcClient(chainCfg)
	if err != nil {
		xylog.Logger.Errorf("new btc client failed, err: %v, ctx=%v", err, ctx)
	}
	if err != nil {
		return nil, err
	}
	return &BClient{btcClient: client}, nil
}
func (b BClient) BlockNumber(ctx context.Context) (uint64, error) {
	return b.btcClient.BlockNumber(ctx)
}

func (b BClient) BlockByNumber(ctx context.Context, number *big.Int) (*xycommon.RpcBlock, error) {

	block, err := b.btcClient.BlockByNumber(ctx, number)
	if err != nil {
		xylog.Logger.Errorf("scan call rpc BlockByNumber[%d], err=%s", number, err)
		return nil, err
	}

	blockEvent, errOrd := b.btcClient.OrdClient.BlockByHash(ctx, block.Hash)
	if errOrd != nil {
		return nil, errOrd
	}
	xylog.Logger.Infof("btc BlockNumber blockEvent[%v]", blockEvent)
	if len(blockEvent.Txs) > 0 {
		xylog.Logger.Infof("btc BlockNumber blockEvent txs : [%v]", blockEvent.Txs)
	}
	rpcBlock := xycommon.RpcBlock{
		Number:     number,
		Hash:       block.Hash,
		Time:       uint64(block.Time),
		RpcTxs:     block.Tx,
		ParentHash: block.PreviousHash,
		Txs:        blockEvent.Txs, // get txs from ord
	}
	return &rpcBlock, nil
}

func (b BClient) HeaderByNumber(ctx context.Context, number *big.Int) (*xycommon.RpcHeader, error) {

	block, err := b.BlockByNumber(ctx, number)
	if err != nil {
		xylog.Logger.Errorf("HeaderByNumber error  BlockByNumber[%d], err=%s", number, err)
		return nil, err
	}
	header := xycommon.RpcHeader{
		ParentHash: block.ParentHash,
		Number:     block.Number,
		Time:       block.Time,
		TxHash:     block.TxHash,
	}
	return &header, nil
}

func (b BClient) TransactionSender(ctx context.Context, txHash, blockHash string, txIndex uint) (string, error) {
	return "", nil
}

func (b BClient) TransactionReceipt(ctx context.Context, txHash string) (*xycommon.RpcReceipt, error) {

	return nil, nil
}

func (b BClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]xycommon.RpcLog, error) {
	return nil, nil
}

func (b BClient) GetAddressBalanceByTick(ctx context.Context, address, tick string) (balance *xycommon.RpcOkxBalance, err error) {
	return b.btcClient.OrdClient.GetAddressBalanceByTick(ctx, address, tick)
}

func (b BClient) GetInscription(ctx context.Context, inscriptionId string) (ins *xycommon.RpcOkxInscription, err error) {
	return b.btcClient.OrdClient.GetInscription(ctx, inscriptionId)
}

func (b BClient) GetRunes(ctx context.Context) ([]xycommon.RpcOrdRunes, error) {
	return b.btcClient.OrdinalsClient.GetRunes(ctx)
}
