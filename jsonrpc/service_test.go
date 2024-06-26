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

package jsonrpc

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func Test_TxHash(t *testing.T) {
	hash := []byte("0x7ffc56b2bf20f4f3474c1fd503fc3f1fb9066c8b0665d6da11185cac892108a5")
	t.Logf("tx_hash=%v", common.Bytes2Hex(hash))
}

func Test_json(t *testing.T) {

	params := make(map[string]interface{})
	params["method"] = "inds_chainStat"
	params["id"] = "1"
	params["jsonrpc"] = "2.0"
	params["params"] = [][]string{[]string{}}

	data, _ := json.Marshal(params)
	t.Logf("data=%v", string(data))

}
