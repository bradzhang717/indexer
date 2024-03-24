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
	"context"
	"encoding/json"
	"fmt"
	"github.com/uxuycom/indexer/xylog"
	"io"
	"net/http"
	"reflect"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	return &HttpClient{
		client: client,
	}
}

func (h *HttpClient) doCallContext(ctx context.Context, method string, url string, out interface{}) error {
	startTs := time.Now()
	defer func() {
		xylog.Logger.Debugf("call api[%s] cost[%v]", url, time.Since(startTs))
	}()

	// check out whether is a pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("out should be a pointer")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		xylog.Logger.Debugf("call http.NewRequestWithContext api[%s] data:[%s] err[%s]", url, "", err)
		return fmt.Errorf("error creating request: %v", err)
	}

	// set headers
	req.Header.Set("Accept", "application/json")

	response, err := h.client.Do(req)
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
		return fmt.Errorf("http client error parsing response body[%s], err[%v]", string(data), err)
	}

	return nil
}

func (h *HttpClient) CallContext(ctx context.Context, method string, url string, out interface{}) (err error) {
	ts := time.Millisecond * 100
	for retry := 0; retry < 5; retry++ {
		err = h.doCallContext(ctx, method, url, out)
		if err == nil {
			return nil
		}
		<-time.After(ts * time.Duration(retry))
	}
	return err
}
