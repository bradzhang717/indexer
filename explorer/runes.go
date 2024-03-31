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

package explorer

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/uxuycom/indexer/client/xycommon"
	"github.com/uxuycom/indexer/devents"
	"github.com/uxuycom/indexer/xyerrors"
	"github.com/uxuycom/indexer/xylog"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

// handleRunesTxs  dispose btc txs and save to db
func (e *Explorer) handleRunesTxs(block *xycommon.RpcBlock) *xyerrors.InsError {

	startTime := time.Now()
	defer func() {
		xylog.Logger.Infof("handle handleRunesTxs use time [%v] block[%v]", time.Since(startTime), block.Number)
	}()
	xylog.Logger.Infof("handleRunesTxs begin RpcTxs[%v] len[%v] block[%v]", block.RpcTxs, len(block.RpcTxs), block.Number)

	var blockTxs []*devents.DBModelEvent
	txIdx, btcTxs := e.convertTxIdx(block.RpcTxs)

	addressBalances, err := e.getRunesAddressBalanceFromBlock(block.Txs, txIdx, block.Number.Int64())
	if err != nil {
		return xyerrors.ErrInternal
	}
	xylog.Logger.Infof("getAddressBalanceFromBlock use time[%v] block[%v]", time.Since(startTime), block.Number)

	startRangTxTime := time.Now()
	for _, tx := range block.OrdTx {

		idx, ok := txIdx[tx.TxId]
		btcTx, _ := btcTxs[tx.TxId]
		if !ok {
			xylog.Logger.Errorf("the transaction does not exist on the block for the original node. txid[%s]", tx.TxId)
			return xyerrors.ErrInternal
		}
		for _, event := range tx.Events {
			if !event.Valid {
				continue
			}
			if event.From == "" || event.To == "" {
				xylog.Logger.Infof("the from or to address is empty and the transaction is skipped. txid[%s] blockHash[%s]", tx.TxId, block.Hash)
				continue
			}
			xylog.Logger.Infof("binding transaction data. txid[%s] block[%d]", tx.TxId, block.Number)
			dm, err := e.buildRunesModel(tx.TxId, event, idx, btcTx, block, addressBalances)
			if err != nil {
				xylog.Logger.Errorf("handleTxs error. err[%s]", err)
				return xyerrors.ErrInternal
			}
			blockTxs = append(blockTxs, dm)
			xylog.Logger.Infof("binding transaction data end. txid[%s] block[%d]", tx.TxId, block.Number)
		}
	}

	xylog.Logger.Infof("handleTxs  end. block[%d] use time[%v]", block.Number, time.Since(startRangTxTime))
	e.writeDBAsync(block, blockTxs)
	return nil
}

func (e *Explorer) buildRunesModel(txid string, event *xycommon.OrdBlockEvent, idx int, btcTx btcjson.TxRawResult, block *xycommon.RpcBlock, addressBalances map[string]*xycommon.AddressBalance) (*devents.DBModelEvent, error) {

	//fromOK, _ := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.From)
	//toOK, _ := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.To)
	//amount, _ := e.convertAmount(event.Amount, event.Tick)
	//xylog.Logger.Infof("buildModel txid[%s]", txid)
	//
	//startTime := time.Now()
	//if event.Type == "deploy" {
	//
	//	xylog.Logger.Infof("buildModel - deploy- GetInscription txid[%s]", txid)
	//	ins, err := e.GetInscription(event.InscriptionId)
	//	xylog.Logger.Infof("buildModel GetInscription end txid[%s]", txid)
	//	if err != nil {
	//		return nil, err
	//	}
	//	e.updateDeployCache(event.Tick, ins.LimitPerMint, ins.TotalSupply, defaultTickDecimals)
	//	xylog.Logger.Infof("buildModel- deploy- updateDeployCache end txid[%s]", txid)
	//} else if event.Type == "mint" {
	//
	//	xylog.Logger.Infof("buildModel- mint- updateMintCache end txid[%s]", txid)
	//	e.updateMintCache(event.Tick, amount, event.To.Address)
	//	xylog.Logger.Infof("buildModel- mint- updateMintCache end txid[%s]", txid)
	//
	//} else if event.Type == "transfer" {
	//	xylog.Logger.Infof("buildModel- transfer- updateTransferCache end txid[%s]", txid)
	//	e.updateTransferCache(event.Tick, event.From.Address, event.To.Address, amount)
	//	xylog.Logger.Infof("buildModel- transfer- updateTransferCache end txid[%s]", txid)
	//
	//} else if event.Type == "inscribeTransfer" {
	//	xylog.Logger.Infof("buildModel- inscribeTransfer- updateInscribeTransferCache end txid[%s]", txid)
	//	e.updateInscribeTransferCache(event.Tick, event.To.Address, txid, event.InscriptionId, amount)
	//	xylog.Logger.Infof("buildModel- inscribeTransfer- updateInscribeTransferCache end txid[%s]", txid)
	//
	//}
	//xylog.Logger.Infof("buildModel- buildTx- begin txid[%s] updateCacheUseTime[%v] block[%v]", txid, time.Since(startTime), block.Number)
	//buildTxStartTime := time.Now()
	//tx, err := e.buildTx(txid, event, idx, btcTx, block)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildTx- err txid[%s], err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//buildInscriptionStartTime := time.Now()
	//xylog.Logger.Infof("buildModel- buildInscription- begin txid[%s] buildTxUseTime[%v]", txid, time.Since(buildTxStartTime))
	//inscription, err := e.buildInscription(txid, event, block)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildInscription- err txid[%s],err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//buildInscriptionStatsStartTime := time.Now()
	//xylog.Logger.Infof("buildModel- buildInscriptionStats- begin txid[%s] buildInscriptionStartTime[%v] block[%v]", txid, time.Since(buildInscriptionStartTime), block.Number)
	//inscriptionStats, err := e.buildInscriptionStats(event, block)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildInscriptionStats- err txid[%s] err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//xylog.Logger.Infof("buildModel- buildBalance- begin txid[%s] buildInscriptionStatsStartTime[%v] block[%v]", txid, time.Since(buildInscriptionStatsStartTime), block.Number)
	//txns, balances, err := e.buildBalance(txid, event, block, fromOK, toOK, addressBalances)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildBalance- err txid[%s], err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//xylog.Logger.Infof("buildModel- buildAddressTx- begin txid[%s]", txid)
	//addressTxs, err := e.buildAddressTx(txid, event, block)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildAddressTx- err txid[%s],err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//xylog.Logger.Infof("buildModel- buildUTXO- begin txid[%s]", txid)
	//utxos, err := e.buildUTXO(event, block)
	//if err != nil {
	//	xylog.Logger.Infof("buildModel- buildUTXO- err txid[%s] err[%v]", txid, err)
	//	return nil, err
	//}
	//
	//xylog.Logger.Infof("buildModel- end txid[%s], err[%v] useTime[%v] block[%v]", txid, err, time.Since(startTime), block.Number)
	//return &devents.DBModelEvent{
	//	Tx:               tx,
	//	Inscriptions:     inscription,
	//	InscriptionStats: inscriptionStats,
	//	BalanceTxs:       txns,
	//	Balances:         balances,
	//	AddressTxs:       addressTxs,
	//	UTXOs:            utxos,
	//}, nil
	return nil, nil
}

func (e *Explorer) getRunesAddressBalanceFromBlock(txs []*xycommon.Tx, txIdx map[string]int, blockNumber int64) (map[string]*xycommon.AddressBalance, error) {

	addrBalancesMap := e.getAddressFromTx(txs, txIdx)
	balanceMap := &sync.Map{}
	g, _ := errgroup.WithContext(context.Background())

	addrChan := make(chan int, 20)
	for _, v := range addrBalancesMap {
		addr := v
		key := fmt.Sprintf("%s_%s", addr.Tick, addr.Address)
		addrChan <- 1
		g.Go(func() error {
			defer func() {
				<-addrChan
			}()
			ordBalance, availableBalance, err := e.GetAddressBalance(addr.Address, addr.Tick)
			xylog.Logger.Infof("address[%s] tick[%s] ordBalance[%v]  availableBalance[%v] err[%v]", addr.Address, addr.Tick, ordBalance, availableBalance, err)
			if err != nil {
				xylog.Logger.Errorf("address[%s] tick[%s] ordBalance[%v]  availableBalance[%v] err[%v]", addr.Address, addr.Tick, ordBalance, availableBalance, err)
				return err
			}
			newAddrBalance := &xycommon.AddressBalance{
				Address:          addr.Address,
				OverallBalance:   ordBalance,
				AvailableBalance: availableBalance,
				Tick:             addr.Tick,
			}

			balanceMap.Store(key, newAddrBalance)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("concurrent block scanning failed. err=%s", err)
	}

	xylog.Logger.Infof("getAddressBalanceFromBlock len:[%v]", len(addrBalancesMap))
	newBalanceMap := make(map[string]*xycommon.AddressBalance)

	for _, v := range addrBalancesMap {
		key := fmt.Sprintf("%s_%s", v.Tick, v.Address)
		xylog.Logger.Infof("getAddressBalanceFromBlock range addressBalances key:[%v]", key)
		if v, ok := balanceMap.Load(key); ok {
			xylog.Logger.Infof("getAddressBalanceFromBlock range addressBalances key:[%v] load", key)
			if vv, ok := v.(*xycommon.AddressBalance); ok {
				xylog.Logger.Infof("getAddressBalanceFromBlock range addressBalances key:[%v] load address[%s] tick[%s] balance[%v] aval[%v]", key, vv.Address, vv.Tick, vv.OverallBalance, vv.AvailableBalance)
				newBalanceMap[key] = vv
			}
		}
	}

	return newBalanceMap, nil
}

func (e *Explorer) getRunesAddressFromTx(txs []btcjson.TxRawResult, txIdx map[string]int) map[string]*xycommon.AddressBalance {

	addrBalances := make(map[string]*xycommon.AddressBalance)
	//for _, tx := range txs {
	//	if _, ok := txIdx[tx.Txid]; !ok {
	//		continue
	//	}
	//
	//	for _, event := range tx.Vin {
	//		if !event.Valid {
	//			xylog.Logger.Infof(" !event.Valid getAddressBalanceFromBlock----- from[%s] to[%s]", event.From.Address, event.To.Address)
	//			continue
	//		}
	//
	//		if event.From.Address == "" || event.To.Address == "" {
	//			xylog.Logger.Infof("the from or to address is empty and the transaction is skipped. txid[%s]", tx.Txid)
	//			continue
	//		}
	//		// deploy balance not change don't need dispose
	//		if event.Type == "deploy" {
	//			continue
	//		}
	//		fromKey := fmt.Sprintf("%s_%s", event.Tick, event.From.Address)
	//		toKey := fmt.Sprintf("%s_%s", event.Tick, event.To.Address)
	//		addrBalances[fromKey] = &xycommon.AddressBalance{
	//			Address: event.From.Address,
	//			Tick:    event.Tick,
	//		}
	//		addrBalances[toKey] = &xycommon.AddressBalance{
	//			Address: event.To.Address,
	//			Tick:    event.Tick,
	//		}
	//	}
	//}
	return addrBalances
}

//
//func (e *Explorer) buildRunesTx(txId string, idx int, tx btcjson.TxRawResult, block *xycommon.RpcBlock) (*model.Transaction, error) {
//
//	amount, err := e.convertAmount(event.Amount, event.Tick)
//	if err != nil {
//		return nil, err
//	}
//	ins, err := e.GetInscription(event.InscriptionId)
//	if err != nil {
//		return nil, err
//	}
//
//	var content string
//	if ins != nil {
//		content = ins.Content
//	}
//	return &model.Transaction{
//		Chain:           e.config.Chain.ChainName,
//		Protocol:        "brc-20",
//		BlockHeight:     block.Number.Uint64(),
//		PositionInBlock: uint64(idx),
//		BlockTime:       time.Unix(int64(block.Time), 0),
//		TxHash:          common.FromHex(txId),
//		From:            event.From.Address,
//		To:              event.To.Address,
//		Op:              event.Type,
//		Tick:            strings.ToLower(event.Tick),
//		Amount:          amount,
//		Gas:             0,
//		GasPrice:        0,
//		CreatedAt:       time.Unix(int64(block.Time), 0),
//		Content:         content,
//	}, nil
//}
