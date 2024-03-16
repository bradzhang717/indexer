package explorer

import (
	"github.com/shopspring/decimal"
	"github.com/uxuycom/indexer/dcache"
)

func (e *Explorer) updateDeployCache(tick string, limit, total decimal.Decimal, decimals int8) {
	t := &dcache.Tick{
		LimitPerMint: limit,
		TotalSupply:  total,
		Decimals:     decimals,
	}
	e.dCache.Inscription.Create(defaultProtocol, tick, t)

	//Add new tick stats
	ts := &dcache.InsStats{
		TxCnt: 1,
	}
	e.dCache.InscriptionStats.Create(defaultProtocol, tick, ts)
}

func (e *Explorer) updateMintCache(tick string, amount decimal.Decimal, address string) {
	//Update mint stats
	e.dCache.InscriptionStats.Mint(defaultProtocol, tick, amount)
	e.dCache.InscriptionStats.TxCnt(defaultProtocol, tick, 1)

	//Update minter balances
	ok, balance := e.dCache.Balance.Get(defaultProtocol, tick, address)
	if !ok {
		//fmt.Printf("=============updateMintCache !ok tick[%s] address[%s]\n", tick, address)
		e.dCache.Balance.Create(defaultProtocol, tick, address, &dcache.BalanceItem{
			Available: amount,
			Overall:   amount,
		})
		//fmt.Printf("=============updateMintCache !ok tick[%s] address[%s] ,a[%v]\n", tick, address, a)
		e.dCache.InscriptionStats.Holders(defaultProtocol, tick, 1)
	} else {
		if balance.Overall.LessThanOrEqual(decimal.Zero) {
			e.dCache.InscriptionStats.Holders(defaultProtocol, tick, 1)
		}
		//fmt.Printf("=============updateMintCache ok tick[%s] address[%s] ,a[%v]\n", tick, address, balance)
		available := balance.Available.Add(amount)
		overall := balance.Overall.Add(amount)
		e.dCache.Balance.Update(defaultProtocol, tick, address, &dcache.BalanceItem{
			Available: available,
			Overall:   overall,
		})
	}
}

func (e *Explorer) updateTransferCache(tick, from, to string, amount decimal.Decimal) {
	//Update transfer stats
	e.dCache.InscriptionStats.TxCnt(defaultProtocol, tick, 1)

	//Update sender balances

	holders := int64(0)
	ok, senderBalance := e.dCache.Balance.Get(defaultProtocol, tick, from)
	if !ok {
		return
	}
	senderAmount := senderBalance.Overall.Sub(amount)
	if senderAmount.LessThanOrEqual(decimal.Zero) {
		holders--
	}
	e.dCache.Balance.Update(defaultProtocol, tick, from, &dcache.BalanceItem{
		Available: senderBalance.Available,
		Overall:   senderAmount,
	})

	// to
	ok, receiveBalance := e.dCache.Balance.Get(defaultProtocol, tick, to)
	if !ok {
		holders++

		receiveAmount := amount
		e.dCache.Balance.Create(defaultProtocol, tick, to, &dcache.BalanceItem{
			Available: receiveAmount,
			Overall:   receiveAmount,
		})

	} else {
		if receiveBalance.Overall.LessThanOrEqual(decimal.Zero) {
			holders++
		}

		availableAmount := receiveBalance.Available.Add(amount)
		overallAmount := receiveBalance.Overall.Add(amount)
		e.dCache.Balance.Update(defaultProtocol, tick, to, &dcache.BalanceItem{
			Available: availableAmount,
			Overall:   overallAmount,
		})
	}

	if holders == 0 {
		return
	}
	//e.dCache.InscriptionStats.Holders(defaultProtocol, tick, holders)
	//e.updateHolders(defaultProtocol, tick)
}

func (e *Explorer) updateInscribeTransferCache(tick, address, txhash, inscriptionID string, amount decimal.Decimal) {
	// update available balance
	ok, balance := e.dCache.Balance.Get(defaultProtocol, tick, address)
	if !ok {
		e.dCache.Balance.Create(defaultProtocol, tick, address, &dcache.BalanceItem{
			Available: decimal.Zero,
			Overall:   amount,
		})
		return
	}
	available := balance.Available.Sub(amount)
	e.dCache.Balance.Update(defaultProtocol, tick, address, &dcache.BalanceItem{
		Available: available,
		Overall:   balance.Overall,
	})
	// add utxo record
	e.dCache.UTXO.Add(defaultProtocol, tick, txhash, address, amount, inscriptionID)
}
