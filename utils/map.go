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

import "sync"

// SafeOrderedMap is a safe order map
type SafeOrderedMap struct {
	sync.Mutex
	keys []string
	m    map[string]interface{}
}

// NewSafeOrderedMap New SafeOrderedMap instance
func NewSafeOrderedMap() *SafeOrderedMap {
	return &SafeOrderedMap{
		keys: make([]string, 0),
		m:    make(map[string]interface{}),
	}
}

// Set key and value
func (som *SafeOrderedMap) Set(key string, value interface{}) {
	som.Lock()
	defer som.Unlock()

	if _, exists := som.m[key]; !exists {
		som.keys = append(som.keys, key)
	}
	som.m[key] = value
}

// Get the value by key
func (som *SafeOrderedMap) Get(key string) (interface{}, bool) {
	som.Lock()
	defer som.Unlock()

	val, exists := som.m[key]
	return val, exists
}

// Delete delete key
func (som *SafeOrderedMap) Delete(key string) {
	som.Lock()
	defer som.Unlock()

	if _, exists := som.m[key]; exists {
		delete(som.m, key)
		// Remove the key from the keys slice
		for i, k := range som.keys {
			if k == key {
				som.keys = append(som.keys[:i], som.keys[i+1:]...)
				break
			}
		}
	}
}

// Keys get all keys from safe ordered map
func (som *SafeOrderedMap) Keys() []string {
	som.Lock()
	defer som.Unlock()
	return som.keys
}
