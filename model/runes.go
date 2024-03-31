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

package model

import (
	"math/big"
	"time"
)

type Runes struct {
	Id           int64     `gorm:"primaryKey" json:"id"` // ID
	Burned       *big.Int  `gorm:"burned" json:"burned"`
	Divisibility int       `gorm:"divisibility" json:"divisibility"`
	Etching      string    `gorm:"etching" json:"etching"`
	Mints        int64     `gorm:"mints" json:"mints"`
	Number       int64     `gorm:"number" json:"number"`
	Rune         string    `gorm:"rune" json:"rune"`
	Spacers      int64     `gorm:"spacers" json:"spacers"`
	Supply       *big.Int  `gorm:"supply" json:"supply"`
	Symbol       string    `gorm:"symbol" json:"symbol"`
	Deadline     int64     `gorm:"deadline" json:"deadline"` // mint deadline
	End          int64     `gorm:"end" json:"end"`           // mint end
	Limit        int64     `gorm:"limit" json:"limit"`       // mint limit
	BlockHeight  int64     `gorm:"block_height" json:"block_height"`
	Index        int64     `gorm:"index" json:"index"`
	Timestamp    time.Time `gorm:"burned" json:"timestamp"`
}

func (Runes) TableName() string {
	return "runes"
}
