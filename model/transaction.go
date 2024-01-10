package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TxEvent int8

const (
	TransactionEventDeploy   TxEvent = 1
	TransactionEventMint     TxEvent = 2
	TransactionEventTransfer TxEvent = 3
	TransactionEventList     TxEvent = 4
	TransactionEventExchange TxEvent = 5
)

type AddressTxs struct {
	ID       uint64          `gorm:"primaryKey" json:"id"`
	Event    TxEvent         `json:"event" gorm:"column:event"`
	TxHash   string          `json:"tx_hash" gorm:"column:tx_hash"`
	Address  string          `json:"address" gorm:"column:address"`
	Amount   decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Tick     string          `json:"tick" gorm:"column:tick"`
	Protocol string          `json:"protocol" gorm:"column:protocol"`
	Operate  string          `json:"operate" gorm:"column:operate"`
	//Desc      string          `json:"desc" gorm:"column:desc"`
	Chain     string    `json:"chain" gorm:"column:chain"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (AddressTxs) TableName() string {
	return "address_txs"
}

type BalanceTxn struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`
	Chain     string          `json:"chain" gorm:"column:chain"`
	Protocol  string          `json:"protocol" gorm:"column:protocol"`
	Event     TxEvent         `json:"event" gorm:"column:event"`
	Address   string          `json:"address" gorm:"column:address"`
	Tick      string          `json:"tick" gorm:"column:tick"`
	Amount    decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Available decimal.Decimal `json:"available" gorm:"column:available;type:decimal(36,18)"`
	Balance   decimal.Decimal `json:"balance" gorm:"column:balance;type:decimal(36,18)"`
	TxHash    string          `json:"tx_hash" gorm:"column:tx_hash"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
}

func (BalanceTxn) TableName() string {
	return "balance_txn"
}

type Transaction struct {
	ID              uint64          `gorm:"primaryKey" json:"id"`
	Chain           string          `json:"chain" gorm:"column:chain"`                         // chain name
	Protocol        string          `json:"protocol" gorm:"column:protocol"`                   // protocol name
	BlockHeight     uint64          `json:"block_height" gorm:"column:block_height"`           // block height
	PositionInBlock uint64          `json:"position_in_block" gorm:"column:position_in_block"` // Position in Block
	BlockTime       time.Time       `json:"block_time" gorm:"column:block_time"`               // block time
	TxHash          string          `json:"tx_hash" gorm:"column:tx_hash"`                     // tx hash
	From            string          `json:"from" gorm:"column:from"`                           // from address
	To              string          `json:"to" gorm:"column:to"`                               // to address
	Op              string          `json:"op" gorm:"column:op"`                               // op code
	Tick            string          `json:"tick" gorm:"column:tick"`                           // inscription code
	Amount          decimal.Decimal `json:"amt" gorm:"column:amt;type:decimal(36,18)"`         // balance
	Gas             int64           `json:"gas" gorm:"column:gas"`                             // gas
	GasPrice        int64           `json:"gas_price" gorm:"column:gas_price"`                 // gas price
	Input           string          `json:"input" gorm:"column:input"`                         // tx content, json string
	Status          int8            `json:"status" gorm:"column:status"`                       // tx status
	CreatedAt       time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at"`
}

func (Transaction) TableName() string {
	return "txs"
}

type AddressTransaction struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`
	Event     int8            `json:"event" gorm:"column:event"`
	TxHash    string          `json:"tx_hash" gorm:"column:tx_hash"`
	Address   string          `json:"address" gorm:"column:address"`
	Amount    decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(36,18)"`
	Tick      string          `json:"tick" gorm:"column:tick"`
	Protocol  string          `json:"protocol" gorm:"column:protocol"`
	Operate   string          `json:"operate" gorm:"column:operate"`
	Chain     string          `json:"chain" gorm:"column:chain"`
	Status    int8            `json:"status" gorm:"column:status"` // tx status
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
}