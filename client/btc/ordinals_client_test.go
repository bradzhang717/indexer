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
// SOFTWARE

package btc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/uxuycom/indexer/xylog"
	"testing"
)

var endpoint = "http://47.128.183.199:8081"

func Test_BlockNumber(t *testing.T) {

	if lv, err := logrus.ParseLevel("debug"); err == nil {
		xylog.InitLog(lv, "")
	}
	ordinalsClient := NewOrdinalsClient(endpoint)
	number, err := ordinalsClient.BlockNumber(context.Background())
	if err != nil {
		t.Logf("test Test_BlockNumber error =%v", err)
	}
	t.Logf("Test_BlockNumber number =%v", number)

}

func Test_GetTransactionByTxId(t *testing.T) {

	if lv, err := logrus.ParseLevel("debug"); err == nil {
		xylog.InitLog(lv, "")
	}
	ordinalsClient := NewOrdinalsClient(endpoint)

	txId := "2b2fa77ca54ad5a58d506d933dcaa2ed1efb41bf5d8b0ce82afbace399a014b1"
	rsp, err := ordinalsClient.GetTransactionByTxId(context.Background(), txId)
	if err != nil {
		t.Logf("test Test_BlockNumber error =%v", err)
	}
	t.Logf("Test_BlockNumber number =%v", rsp)

}

func Test_GetOutput(t *testing.T) {

	if lv, err := logrus.ParseLevel("debug"); err == nil {
		xylog.InitLog(lv, "")
	}
	ordinalsClient := NewOrdinalsClient(endpoint)
	output := "2b2fa77ca54ad5a58d506d933dcaa2ed1efb41bf5d8b0ce82afbace399a014b1:0"
	rsp, err := ordinalsClient.GetOutput(context.Background(), output)
	if err != nil {
		t.Logf("test Test_BlockNumber error =%v", err)
	}
	t.Logf("Test_BlockNumber number =%v", rsp)
}

func Test_GetRunes(t *testing.T) {

	if lv, err := logrus.ParseLevel("debug"); err == nil {
		xylog.InitLog(lv, "")
	}

	ordinalsClient := NewOrdinalsClient(endpoint)
	rsp, err := ordinalsClient.GetRunes(context.Background())
	if err != nil {
		t.Logf("test Test_GetRunes error =%v", err)
	}
	t.Logf("Test_GetRunes number =%v", rsp)

}

func Test_GetRune(t *testing.T) {

	if lv, err := logrus.ParseLevel("debug"); err == nil {
		xylog.InitLog(lv, "")
	}
	runes := "IIIIJJGFDGGB"
	ordinalsClient := NewOrdinalsClient(endpoint)
	rsp, err := ordinalsClient.GetRune(context.Background(), runes)
	if err != nil {
		t.Logf("test Test_GetRunes error =%v", err)
	}
	t.Logf("Test_GetRunes number =%v", rsp)

}
