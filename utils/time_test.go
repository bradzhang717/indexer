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

package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"testing"
	"time"
)

func Test_Hour(t *testing.T) {
	t.Logf("hour: %v", Hour(time.Now()))
}

func Test_Yesterday(t *testing.T) {
	t.Logf("hour: %v", YesterdayHour())
}

func Test_BeforeYesterdayHour(t *testing.T) {
	t.Logf("hour: %v", BeforeYesterdayHour())
}

func Test_TimeHour(t *testing.T) {
	now := Hour(time.Now())
	yesterday := YesterdayHour()
	beforeYesterday := BeforeYesterdayHour()
	t.Logf("now: %v", TimeHourInt(now))
	t.Logf("yesterday: %v", TimeHourInt(yesterday))
	t.Logf("beforeYesterday: %v", TimeHourInt(beforeYesterday))
}

func Test_TimeLineFormat(t *testing.T) {

	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	t.Logf("time %v", TimeLineFormat(time.Now()))
	t.Logf("time %v", formattedTime)

}

func Test_TimeFormatHours(t *testing.T) {
	day := 20240314
	begin := TimeFormatDayHours(day)
	t.Logf("begin: %v", begin)

}

func Test_TimeFormatHourBeginAndEnd(t *testing.T) {
	day := 20240314
	hours := TimeFormatDayHours(day)

	for _, value := range hours {
		s, e := TimeFormatHourBeginAndEnd(int(value))
		t.Logf("begin: %v, end:=%v", s, e)
	}

}

func Test_TimeHash(t *testing.T) {

	// Example input
	input := "7397969152f26e7dc46c6236a048cdb8ed105910b4540114325c4b984a091b91"
	reversed, err := reverseBytes(input)
	if err != nil {
		fmt.Println("Error reversing bytes:", err)
		return
	}
	fmt.Printf("Input: N = 0x%s\nOutput: 0x%s\n", input, reversed)

	// Your integer value
	var myInt int8 = 3
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, myInt)
	// Convert the bytes.Buffer to a []byte
	byteSlice := buf.Bytes()
	// Print the resulting []byte
	//fmt.Println(byteSlice)

	builder := txscript.NewScriptBuilder()
	builder.AddData(byteSlice)
	ret, _ := builder.Script()
	fmt.Println(hex.EncodeToString(ret))

	fmt.Println(string([]byte("␃")))
}

// reverseBytes takes a hex string representing a large number and returns its byte-reversed hex string
func reverseBytes(hexStr string) (string, error) {
	// Decode the hex string to bytes
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	//fmt.Printf("Middile:0x%s\n", string(data))

	// Reverse the bytes
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	// Encode the bytes back to hex string and return
	return hex.EncodeToString(data), nil

}
