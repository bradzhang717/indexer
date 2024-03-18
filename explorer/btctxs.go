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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/uxuycom/indexer/client/btc"
	"github.com/uxuycom/indexer/client/xycommon"
	"github.com/uxuycom/indexer/dcache"
	"github.com/uxuycom/indexer/devents"
	"github.com/uxuycom/indexer/model"
	"github.com/uxuycom/indexer/xyerrors"
	"github.com/uxuycom/indexer/xylog"
	"golang.org/x/sync/errgroup"
	"math"
	"strings"
	"sync"
	"time"
)

var defaultProtocol = "brc-20"
var defaultTickDecimals int8 = 18
var statsDecimals int8 = 8
var decimals = 8

// handleBtcTxs  dispose btc txs and save to db
func (e *Explorer) handleBtcTxs(block *xycommon.RpcBlock) *xyerrors.InsError {

	startTime := time.Now()
	defer func() {
		xylog.Logger.Infof("handle tx use time [%v] block[%v]", time.Since(startTime), block.Number)
	}()
	xylog.Logger.Infof("handleTxs begin txs[%v] len[%v] block[%v]", block.Txs, len(block.Txs), block.Number)

	var blockTxs []*devents.DBModelEvent
	txIdx, btcTxs := e.convertTxIdx(block.RpcTxs)
	addressBalances, err := e.getAddressBalanceFromBlock(block.Txs, txIdx, block.Number.Int64())
	if err != nil {
		return xyerrors.ErrInternal
	}
	xylog.Logger.Infof("getAddressBalanceFromBlock use time[%v] block[%v]", time.Since(startTime), block.Number)

	startRangTxTime := time.Now()
	for _, tx := range block.Txs {

		idx, ok := txIdx[tx.Txid]
		btcTx, _ := btcTxs[tx.Txid]
		if !ok {
			xylog.Logger.Errorf("the transaction does not exist on the block for the original node. txid[%s]", tx.Txid)
			return xyerrors.ErrInternal
		}
		for _, event := range tx.Events {
			if !event.Valid {
				continue
			}
			if event.From.Address == "" || event.To.Address == "" {
				xylog.Logger.Infof("the from or to address is empty and the transaction is skipped. txid[%s] blockHash[%s]", tx.Txid, block.Hash)
				continue
			}
			xylog.Logger.Infof("binding transaction data. txid[%s] block[%d]", tx.Txid, block.Number)
			dm, err := e.buildModel(tx.Txid, event, idx, btcTx, block, addressBalances)
			if err != nil {
				xylog.Logger.Errorf("handleTxs error. err[%s]", err)
				return xyerrors.ErrInternal
			}
			blockTxs = append(blockTxs, dm)
			xylog.Logger.Infof("binding transaction data end. txid[%s] block[%d]", tx.Txid, block.Number)
		}
	}

	xylog.Logger.Infof("handleTxs  end. block[%d] use time[%v]", block.Number, time.Since(startRangTxTime))
	e.writeDBAsync(block, blockTxs)
	return nil
}

func (e *Explorer) convertTxIdx(txs []btcjson.TxRawResult) (map[string]int, map[string]btcjson.TxRawResult) {
	idxs := make(map[string]int, len(txs))
	newTxs := make(map[string]btcjson.TxRawResult, len(txs))
	for idx, v := range txs {
		idxs[v.Txid] = idx
		newTxs[v.Txid] = v
	}
	return idxs, newTxs
}

func (e *Explorer) getAddressBalanceFromBlock(txs []*xycommon.Tx, txIdx map[string]int, blockNumber int64) (map[string]*xycommon.AddressBalance, error) {

	addrBals := e.getAddressFromTx(txs, txIdx)
	balanceMap := &sync.Map{}
	g, _ := errgroup.WithContext(context.Background())

	addrChan := make(chan int, 20)
	//newBlockNumber, err1 := btcClient.BlockNumber(context.Background())
	//if err1 != nil {
	//	return nil, err1
	//}
	for _, v := range addrBals {
		addr := v
		key := fmt.Sprintf("%s_%s", addr.Tick, addr.Address)
		addrChan <- 1
		//logs.Logger.Infof("------- address[%s] tick[%s]", addr.Address, addr.Tick)
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

			////// balance not change send lark
			//ok, cBalance := e.dCache.Balance.Get(defaultProtocol, addr.Tick, addr.Address)
			//if ok && cBalance.Overall.Equal(ordBalance) && cBalance.Available.Equal(availableBalance) && newBlockNumber == blockNumber {
			//	go e.larkNotification(addr, cBalance, ordBalance, availableBalance, blockNumber)
			//	xylog.Logger.Infof("the balance is not as expected. address[%s] tick[%s] ordBalance[%v]  availableBalance[%v] cacheBalance[%v] cacheAvailableBalance[%v]", addr.Address, addr.Tick, ordBalance, availableBalance, cBalance.Overall, cBalance.Available)
			//}

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

	xylog.Logger.Infof("getAddressBalanceFromBlock len:[%v]", len(addrBals))
	newBalanceMap := make(map[string]*xycommon.AddressBalance)

	for _, v := range addrBals {
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

func (e *Explorer) getAddressFromTx(txs []*xycommon.Tx, txIdx map[string]int) map[string]*xycommon.AddressBalance {

	addrBalances := make(map[string]*xycommon.AddressBalance)
	for _, tx := range txs {
		if _, ok := txIdx[tx.Txid]; !ok {
			continue
		}

		for _, event := range tx.Events {
			if !event.Valid {
				xylog.Logger.Infof(" !event.Valid getAddressBalanceFromBlock----- from[%s] to[%s]", event.From.Address, event.To.Address)
				continue
			}

			if event.From.Address == "" || event.To.Address == "" {
				xylog.Logger.Infof("the from or to address is empty and the transaction is skipped. txid[%s]", tx.Txid)
				continue
			}
			// deploy balance not change don't need dispose
			if event.Type == "deploy" {
				continue
			}
			fromKey := fmt.Sprintf("%s_%s", event.Tick, event.From.Address)
			toKey := fmt.Sprintf("%s_%s", event.Tick, event.To.Address)
			addrBalances[fromKey] = &xycommon.AddressBalance{
				Address: event.From.Address,
				Tick:    event.Tick,
			}
			addrBalances[toKey] = &xycommon.AddressBalance{
				Address: event.To.Address,
				Tick:    event.Tick,
			}
		}
	}
	return addrBalances
}

func (e *Explorer) GetAddressBalance(address string, tick string) (decimal.Decimal, decimal.Decimal, error) {

	bClient := e.node.(*btc.BClient)
	ordBalance, err := bClient.GetAddressBalanceByTick(context.Background(), address, tick)

	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	if ordBalance == nil {
		xylog.Logger.Errorf("GetAddressBalance is nil. address[%s] tick[%s]", address, tick)
		return decimal.Zero, decimal.Zero, nil
	}
	balance, err := e.convertAmount(ordBalance.OverallBalance, tick)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	availableBalance, err := e.convertAmount(ordBalance.AvailableBalance, tick)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	return balance, availableBalance, nil
}

func (e *Explorer) convertBitcoinStats(value decimal.Decimal) decimal.Decimal {
	if value.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	return value.Mul(decimal.NewFromFloat(math.Pow10(decimals))).Round(int32(decimals))
}

func (e *Explorer) convertAmount(amount string, tick string) (decimal.Decimal, error) {
	if amount == "" {
		return decimal.Zero, nil
	}
	value, err := decimal.NewFromString(amount)
	if err != nil {
		xylog.Logger.Errorf("convertAmount error buildUTXO. err[%s] amount[%v] tick[%s]", err, amount, tick)
		return decimal.Zero, err
	}
	//decimals := 18
	a := value.Shift(int32(-defaultTickDecimals))
	return a, nil
}

func (e *Explorer) buildModel(txid string, event *xycommon.BlockEvent, idx int, btcTx btcjson.TxRawResult, block *xycommon.RpcBlock, addressBalances map[string]*xycommon.AddressBalance) (*devents.DBModelEvent, error) {

	fromOK, _ := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.From.Address)
	toOK, _ := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.To.Address)
	amount, _ := e.convertAmount(event.Amount, event.Tick)
	xylog.Logger.Infof("buildModel txid[%s]", txid)

	startTime := time.Now()
	if event.Type == "deploy" {

		xylog.Logger.Infof("buildModel - deploy- GetInscription txid[%s]", txid)
		ins, err := e.GetInscription(event.InscriptionId)
		xylog.Logger.Infof("buildModel GetInscription end txid[%s]", txid)
		if err != nil {
			return nil, err
		}
		e.updateDeployCache(event.Tick, ins.LimitPerMint, ins.TotalSupply, defaultTickDecimals)
		xylog.Logger.Infof("buildModel- deploy- updateDeployCache end txid[%s]", txid)
	} else if event.Type == "mint" {

		xylog.Logger.Infof("buildModel- mint- updateMintCache end txid[%s]", txid)
		e.updateMintCache(event.Tick, amount, event.To.Address)
		xylog.Logger.Infof("buildModel- mint- updateMintCache end txid[%s]", txid)

	} else if event.Type == "transfer" {
		xylog.Logger.Infof("buildModel- transfer- updateTransferCache end txid[%s]", txid)
		e.updateTransferCache(event.Tick, event.From.Address, event.To.Address, amount)
		xylog.Logger.Infof("buildModel- transfer- updateTransferCache end txid[%s]", txid)

	} else if event.Type == "inscribeTransfer" {
		xylog.Logger.Infof("buildModel- inscribeTransfer- updateInscribeTransferCache end txid[%s]", txid)
		e.updateInscribeTransferCache(event.Tick, event.To.Address, txid, event.InscriptionId, amount)
		xylog.Logger.Infof("buildModel- inscribeTransfer- updateInscribeTransferCache end txid[%s]", txid)

	}
	xylog.Logger.Infof("buildModel- buildTx- begin txid[%s] updateCacheUseTime[%v] block[%v]", txid, time.Since(startTime), block.Number)
	buildTxStartTime := time.Now()
	tx, err := e.buildTx(txid, event, idx, btcTx, block)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildTx- err txid[%s], err[%v]", txid, err)
		return nil, err
	}

	buildInscriptionStartTime := time.Now()
	xylog.Logger.Infof("buildModel- buildInscription- begin txid[%s] buildTxUseTime[%v]", txid, time.Since(buildTxStartTime))
	inscription, err := e.buildInscription(txid, event, block)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildInscription- err txid[%s],err[%v]", txid, err)
		return nil, err
	}

	buildInscriptionStatsStartTime := time.Now()
	xylog.Logger.Infof("buildModel- buildInscriptionStats- begin txid[%s] buildInscriptionStartTime[%v] block[%v]", txid, time.Since(buildInscriptionStartTime), block.Number)
	inscriptionStats, err := e.buildInscriptionStats(event, block)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildInscriptionStats- err txid[%s] err[%v]", txid, err)
		return nil, err
	}

	xylog.Logger.Infof("buildModel- buildBalance- begin txid[%s] buildInscriptionStatsStartTime[%v] block[%v]", txid, time.Since(buildInscriptionStatsStartTime), block.Number)
	txns, balances, err := e.buildBalance(txid, event, block, fromOK, toOK, addressBalances)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildBalance- err txid[%s], err[%v]", txid, err)
		return nil, err
	}

	xylog.Logger.Infof("buildModel- buildAddressTx- begin txid[%s]", txid)
	addressTxs, err := e.buildAddressTx(txid, event, block)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildAddressTx- err txid[%s],err[%v]", txid, err)
		return nil, err
	}

	xylog.Logger.Infof("buildModel- buildUTXO- begin txid[%s]", txid)
	utxos, err := e.buildUTXO(event, block)
	if err != nil {
		xylog.Logger.Infof("buildModel- buildUTXO- err txid[%s] err[%v]", txid, err)
		return nil, err
	}

	xylog.Logger.Infof("buildModel- end txid[%s], err[%v] useTime[%v] block[%v]", txid, err, time.Since(startTime), block.Number)
	return &devents.DBModelEvent{
		Tx:               tx,
		Inscriptions:     inscription,
		InscriptionStats: inscriptionStats,
		BalanceTxs:       txns,
		Balances:         balances,
		AddressTxs:       addressTxs,
		UTXOs:            utxos,
	}, nil
}

func (e *Explorer) buildUTXO(event *xycommon.BlockEvent, block *xycommon.RpcBlock) (map[devents.DBAction]*model.UTXO, error) {
	amount, err := e.convertAmount(event.Amount, event.Tick)
	if err != nil {
		xylog.Logger.Errorf("convertAmount error buildUTXO. err[%s]", err)
		return nil, err
	}
	if event.Type == "inscribeTransfer" {
		return map[devents.DBAction]*model.UTXO{
			devents.DBActionCreate: {
				Chain:    e.config.Chain.ChainName,
				Protocol: defaultProtocol,
				Tick:     strings.ToLower(event.Tick),
				Sn:       event.InscriptionId,
				Status:   model.UTXOStatusUnspent,
				RootHash: block.TxHash,
				Address:  event.To.Address,
				Amount:   amount,
			},
		}, nil
	}

	if event.Type == "transfer" {
		return map[devents.DBAction]*model.UTXO{
			devents.DBActionUpdate: {
				Chain:   e.config.Chain.ChainName,
				Sn:      event.InscriptionId,
				Address: event.From.Address,
				Status:  model.UTXOStatusSpent,
			},
		}, nil
	}

	return nil, nil
}

func (e *Explorer) buildAddressTx(txId string, event *xycommon.BlockEvent, block *xycommon.RpcBlock) ([]*model.AddressTxs, error) {
	var txs []*model.AddressTxs
	amount, err := e.convertAmount(event.Amount, event.Tick)
	if err != nil {
		return nil, err
	}
	txs = append(txs, &model.AddressTxs{
		Event:     e.getEventByOperate(event.Type),
		Address:   event.To.Address,
		Amount:    amount,
		TxHash:    common.FromHex(txId),
		Tick:      strings.ToLower(event.Tick),
		Protocol:  defaultProtocol,
		Operate:   event.Type,
		Chain:     e.config.Chain.ChainName,
		CreatedAt: time.Unix(int64(block.Time), 0),
	})
	if !strings.EqualFold(event.From.Address, event.To.Address) {
		txs = append(txs, &model.AddressTxs{
			Event:     e.getEventByOperate(event.Type),
			Address:   event.From.Address,
			Amount:    amount,
			TxHash:    common.FromHex(txId),
			Tick:      strings.ToLower(event.Tick),
			Protocol:  defaultProtocol,
			Operate:   event.Type,
			Chain:     e.config.Chain.ChainName,
			CreatedAt: time.Unix(int64(block.Time), 0),
		})
	}

	return txs, nil
}

func (e *Explorer) buildBalance(txid string, event *xycommon.BlockEvent, block *xycommon.RpcBlock, fromOK, toOK bool, addressBalances map[string]*xycommon.AddressBalance) ([]*model.BalanceTxn, map[devents.DBAction][]*model.Balances, error) {
	var txns []*model.BalanceTxn
	//var balances map[DBAction][]*model.Balances
	balances := make(map[devents.DBAction][]*model.Balances, 2)

	if event.Type == "mint" {
		amount, err := e.convertAmount(event.Amount, event.Tick)
		if err != nil {
			return nil, nil, err
		}
		ok, toBalance := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.To.Address)
		if !ok {
			toBalance = e.dCache.Balance.Create(defaultProtocol, event.Tick, event.To.Address, &dcache.BalanceItem{
				Available: amount,
				Overall:   amount,
			})
			fmt.Printf("txid[%s] event.Tick[%s] event.To.Address[%s]\n", txid, event.Tick, event.To.Address)
		}
		//overallBalance := toBalance.Overall.Add(amount)
		//availableBalance := toBalance.Available.Add(amount)
		txns = append(txns, &model.BalanceTxn{
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Event:     e.getEventByOperate(event.Type),
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Amount:    amount,
			Balance:   toBalance.Overall,
			Available: toBalance.Available,
			TxHash:    common.FromHex(txid),
			CreatedAt: time.Unix(int64(block.Time), 0),
		})

		var ordBalance decimal.Decimal
		var availableBalance decimal.Decimal
		//if event.To.Address == "tb1pxc0h9ccs62el098z8ckjts0hm4xc3ndn959hj2tsnwhllqup446safkerz" {
		//	fmt.Println("tb1pxc0h9ccs62el098z8ckjts0hm4xc3ndn959hj2tsnwhllqup446safkerz")
		//}
		toKey := fmt.Sprintf("%s_%s", event.Tick, event.To.Address)
		//logs.Logger.Infof("%s: %v", toKey, addressBalances)
		addrBalance, ok := addressBalances[toKey]
		if ok {
			ordBalance = addrBalance.OverallBalance
			availableBalance = addrBalance.AvailableBalance
			xylog.Logger.Infof("ok  address[%s] ordBalance[%v] availableBalance[%v]", event.To.Address, ordBalance, availableBalance)
		} else {
			var err1 error
			ordBalance, availableBalance, err1 = e.GetAddressBalance(event.To.Address, event.Tick)
			xylog.Logger.Infof("!ok  address[%s] ordBalance[%v] availableBalance[%v], err1[%s]", event.To.Address, ordBalance, availableBalance, err1)
			if err1 != nil {
				return nil, nil, err
			}
		}

		action := devents.DBActionUpdate
		if !toOK {
			action = devents.DBActionCreate
		}
		if _, ok := balances[action]; !ok {
			balances[action] = make([]*model.Balances, 0, 1)
		}

		balances[action] = setBalances(balances[action], &model.Balances{
			SID:       toBalance.SID,
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Balance:   ordBalance,
			Available: availableBalance,
		}, action)

	} else if event.Type == "transfer" {
		amount, err := e.convertAmount(event.Amount, event.Tick)
		if err != nil {
			return nil, nil, err
		}
		tok, toBalance := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.To.Address)
		toOverall := decimal.Zero
		toAvailable := decimal.Zero
		if tok {
			toOverall = toBalance.Overall
			toAvailable = toBalance.Available
		} else {
			toBalance = e.dCache.Balance.Create(defaultProtocol, event.Tick, event.To.Address, &dcache.BalanceItem{
				Available: amount,
				Overall:   amount,
			})
		}
		txns = append(txns, &model.BalanceTxn{
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Event:     e.getEventByOperate(event.Type),
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Amount:    amount,
			Balance:   toOverall,
			Available: toAvailable,
			TxHash:    common.FromHex(txid),
			CreatedAt: time.Unix(int64(block.Time), 0),
		})

		action := devents.DBActionUpdate
		if !toOK {
			action = devents.DBActionCreate
		}
		if _, ok := balances[action]; !ok {
			balances[action] = make([]*model.Balances, 0, 1)
		}

		//ordBalance, ordAvailableBalance, err := e.GetAddressBalance(event.To.Address, event.Tick)
		//if err != nil {
		//	return nil, nil, err
		//}
		var ordBalance decimal.Decimal
		var availableBalance decimal.Decimal
		toKey := fmt.Sprintf("%s_%s", event.Tick, event.To.Address)
		addrBalance, ok := addressBalances[toKey]
		if ok {
			ordBalance = addrBalance.OverallBalance
			availableBalance = addrBalance.AvailableBalance
		} else {
			var err1 error
			ordBalance, availableBalance, err1 = e.GetAddressBalance(event.To.Address, event.Tick)
			if err1 != nil {
				return nil, nil, err
			}
		}
		balances[action] = setBalances(balances[action], &model.Balances{
			SID:       toBalance.SID,
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Balance:   ordBalance,
			Available: availableBalance,
		}, action)

		ok, senderBalance := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.From.Address)
		if !ok {
			senderBalance = e.dCache.Balance.Create(defaultProtocol, event.Tick, event.From.Address, &dcache.BalanceItem{
				Available: amount,
				Overall:   amount,
			})
		}
		//senderOverallBalance := senderBalance.Overall.Sub(amount)
		//senderAvailableBalance := senderBalance.Available.Sub(amount)
		txns = append(txns, &model.BalanceTxn{
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Event:     e.getEventByOperate(event.Type),
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Amount:    amount.Neg(),
			Balance:   senderBalance.Overall,
			Available: senderBalance.Available,
			TxHash:    common.FromHex(txid),
			CreatedAt: time.Unix(int64(block.Time), 0),
		})

		senderAction := devents.DBActionUpdate
		if !fromOK {
			senderAction = devents.DBActionCreate
		}
		if _, ok := balances[senderAction]; !ok {
			balances[senderAction] = make([]*model.Balances, 0, 1)
		}
		//senderOrdBalance, senderAvailableBalance, err := e.GetAddressBalance(event.From.Address, event.Tick)
		//if err != nil {
		//	return nil, nil, err
		//}
		var senderOrdBalance decimal.Decimal
		var senderAvailableBalance decimal.Decimal
		fromKey := fmt.Sprintf("%s_%s", event.Tick, event.From.Address)
		fromAddrBalance, ok := addressBalances[fromKey]
		if ok {
			senderOrdBalance = fromAddrBalance.OverallBalance
			senderAvailableBalance = fromAddrBalance.AvailableBalance
		} else {
			var err1 error
			senderOrdBalance, senderAvailableBalance, err1 = e.GetAddressBalance(event.To.Address, event.Tick)
			if err1 != nil {
				return nil, nil, err
			}
		}

		balances[senderAction] = setBalances(balances[senderAction], &model.Balances{
			SID:       senderBalance.SID,
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Balance:   senderOrdBalance,
			Available: senderAvailableBalance,
		}, senderAction)

	} else if event.Type == "inscribeTransfer" {
		amount, err := e.convertAmount(event.Amount, event.Tick)
		if err != nil {
			return nil, nil, err
		}
		var overallBalance decimal.Decimal
		var avaBalance decimal.Decimal
		toOK, toBalance := e.dCache.Balance.Get(defaultProtocol, event.Tick, event.To.Address)
		action := devents.DBActionUpdate
		if !toOK {
			action = devents.DBActionCreate
			avaBalance = decimal.Zero
			overallBalance = amount
			return nil, nil, fmt.Errorf("abnormal user balance. address[%s] amount[%s]", event.To.Address, amount.String())
		} else {
			overallBalance = toBalance.Overall
			avaBalance = toBalance.Available
		}
		//overallBalance := toBalance.Overall.Add(amount)
		//availableBalance := toBalance.Available.Sub(amount)
		txns = append(txns, &model.BalanceTxn{
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Event:     e.getEventByOperate(event.Type),
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Amount:    amount,
			Balance:   overallBalance,
			Available: avaBalance,
			TxHash:    common.FromHex(txid),
			CreatedAt: time.Unix(int64(block.Time), 0),
		})

		if _, ok := balances[action]; !ok {
			balances[action] = make([]*model.Balances, 0, 1)
		}
		//ordBalance, availableBalance, err := e.GetAddressBalance(event.To.Address, event.Tick)
		//if err != nil {
		//	return nil, nil, err
		//}
		var ordBalance decimal.Decimal
		var availableBalance decimal.Decimal
		toKey := fmt.Sprintf("%s_%s", event.Tick, event.To.Address)
		addrBalance, ok := addressBalances[toKey]
		if ok {
			ordBalance = addrBalance.OverallBalance
			availableBalance = addrBalance.AvailableBalance
		} else {
			var err1 error
			ordBalance, availableBalance, err1 = e.GetAddressBalance(event.To.Address, event.Tick)
			if err1 != nil {
				return nil, nil, err
			}
		}
		balances[action] = setBalances(balances[action], &model.Balances{
			SID:       toBalance.SID,
			Chain:     e.config.Chain.ChainName,
			Protocol:  defaultProtocol,
			Address:   event.To.Address,
			Tick:      strings.ToLower(event.Tick),
			Balance:   ordBalance,
			Available: availableBalance,
		}, action)
	}

	return txns, balances, nil
}

func setBalances(balances []*model.Balances, balance *model.Balances, action devents.DBAction) []*model.Balances {

	if action == devents.DBActionCreate {
		var isOk = false
		for _, item := range balances {
			if item.Tick == balance.Tick &&
				item.Address == balance.Address &&
				item.Protocol == balance.Protocol &&
				item.Chain == balance.Chain {
				item = balance
				isOk = true
			}
		}
		if !isOk {
			balances = append(balances, balance)
		}
	} else {
		balances = append(balances, balance)
	}
	return balances
}

func (e *Explorer) buildInscriptionStats(event *xycommon.BlockEvent, block *xycommon.RpcBlock) (map[devents.DBAction]*model.InscriptionsStats, error) {

	ok, d := e.dCache.InscriptionStats.Get(defaultProtocol, event.Tick)
	if !ok {
		xylog.Logger.Infof("inscription stats not found. inscriptionId[%s] tick[%s]", event.InscriptionId, event.Tick)
		return nil, nil
	}
	data := &model.InscriptionsStats{
		SID:      d.SID,
		Chain:    e.config.Chain.ChainName,
		Protocol: defaultProtocol,
		Tick:     strings.ToLower(event.Tick),
		Minted:   d.Minted,
		Holders:  uint64(d.Holders),
		TxCnt:    d.TxCnt,
	}

	if event.Type == "mint" {
		amount, err := e.convertAmount(event.Amount, event.Tick)
		if err != nil {
			return nil, err
		}
		if d.Minted.Equal(amount) {
			data.MintFirstBlock = block.Number.Uint64()
		}

		// final mint block record
		ok, inscription := e.dCache.Inscription.Get(defaultProtocol, event.Tick)
		if !ok {
			xylog.Logger.Infof("inscription not found. inscriptionId[%s] tick[%s]", event.InscriptionId, event.Tick)
		}
		if inscription.TotalSupply.LessThanOrEqual(d.Minted) {
			data.MintLastBlock = block.Number.Uint64()

			ts := time.Unix(int64(block.Time), 0)
			data.MintCompletedTime = &ts
		}
	}

	if event.Type == "deploy" {
		return map[devents.DBAction]*model.InscriptionsStats{
			devents.DBActionCreate: data,
		}, nil
	} else {
		return map[devents.DBAction]*model.InscriptionsStats{
			devents.DBActionUpdate: data,
		}, nil
	}
}

func (e *Explorer) buildInscription(txid string, event *xycommon.BlockEvent, block *xycommon.RpcBlock) (map[devents.DBAction]*model.Inscriptions, error) {
	if event.Type != "deploy" {
		return nil, nil
	}
	ok, d := e.dCache.Inscription.Get(defaultProtocol, event.Tick)
	if !ok {
		xylog.Logger.Infof("inscription not found. txid[%s] tick[%s]", txid, event.Tick)
		return nil, fmt.Errorf("inscription not found. txid[%s] tick[%s]", txid, event.Tick)
	}
	ret := make(map[devents.DBAction]*model.Inscriptions, 1)
	ins, err := e.GetInscription(event.InscriptionId)
	if err != nil {
		return nil, err
	}
	inscription := &model.Inscriptions{
		SID:          d.SID,
		Chain:        e.config.Chain.ChainName,
		Protocol:     defaultProtocol,
		Tick:         strings.ToLower(event.Tick),
		Name:         ins.Tick,
		LimitPerMint: ins.LimitPerMint,
		TotalSupply:  ins.TotalSupply,
		DeployBy:     ins.Owner,
		DeployHash:   txid,
		DeployTime:   time.Unix(int64(block.Time), 0),
		Decimals:     ins.Decimals,
	}
	ret[devents.DBActionCreate] = inscription

	return ret, nil
}

func (e *Explorer) GetInscription(id string) (*xycommon.Inscription, error) {

	bClient := e.node.(*btc.BClient)
	ins, err := bClient.GetInscription(context.Background(), id)
	if err != nil {
		return nil, err
	}
	content, err := hex.DecodeString(ins.Content)
	if err != nil {
		return nil, err
	}
	insContent := &xycommon.InscriptionContent{}
	err1 := json.Unmarshal(content, &insContent)
	if err1 != nil {
		return nil, err1
	}

	limit, _ := decimal.NewFromString(insContent.LimitPerMint)
	totalSupply, _ := decimal.NewFromString(insContent.Max)

	return &xycommon.Inscription{
		Id:           id,
		Tick:         insContent.Tick,
		LimitPerMint: limit,
		TotalSupply:  totalSupply,
		Decimals:     defaultTickDecimals,
		Owner:        ins.Owner.Address,
	}, nil
}

func (e *Explorer) buildTx(txId string, event *xycommon.BlockEvent, idx int, tx btcjson.TxRawResult, block *xycommon.RpcBlock) (*model.Transaction, error) {

	amount, err := e.convertAmount(event.Amount, event.Tick)
	if err != nil {
		return nil, err
	}
	return &model.Transaction{
		Chain:           e.config.Chain.ChainName,
		Protocol:        "brc-20",
		BlockHeight:     block.Number.Uint64(),
		PositionInBlock: uint64(idx),
		BlockTime:       time.Unix(int64(block.Time), 0),
		TxHash:          common.FromHex(txId),
		From:            event.From.Address,
		To:              event.To.Address,
		Op:              event.Type,
		Tick:            strings.ToLower(event.Tick),
		Amount:          amount,
		Gas:             0,
		GasPrice:        0,
		CreatedAt:       time.Unix(int64(block.Time), 0),
	}, nil
}

func (e *Explorer) getEventByOperate(operate string) model.TxEvent {
	switch operate {
	case OperateDeploy:
		return model.TransactionEventDeploy
	case OperateMint:
		return model.TransactionEventMint
	case OperateTransfer:
		return model.TransactionEventTransfer
	case OperateList:
		return model.TransactionEventList
	case OperateDelist:
		return model.TransactionEventDelist
	case OperateExchange:
		return model.TransactionEventExchange
	case OperateInscribeTransfer:
		return model.TransactionEventInscribeTransfer
	}
	return model.TxEvent(0)
}

const (
	OperateDeploy           string = "deploy"
	OperateMint             string = "mint"
	OperateTransfer         string = "transfer"
	OperateInscribeTransfer string = "inscribeTransfer"
	OperateList             string = "list"
	OperateDelist           string = "delist"
	OperateExchange         string = "exchange"
)
