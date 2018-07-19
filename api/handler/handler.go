/*
 * Copyright (C) 2018 eeonevision
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

var endpoint string

// errNotFound error returned when a document could not be found
var errNotFound = errors.New("not found")

// Request struct represents request related fields.
//
// AccountID - unique identifier of request-maker in blockchain
// PrivKey - private key of account
// PubKey - public key of account
// Data - arbitrary passed data (maybe any supported in system)
type Request struct {
	AccountID string      `json:"account_id,omitempty"`
	PrivKey   string      `json:"private_key,omitempty"`
	PubKey    string      `json:"public_key,omitempty"`
	Data      interface{} `json:"data"`
}

// Result struct represents response from Anychaindb API.
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// MongoQuery is a struct for parse search query from a user.
type mongoQuery struct {
	Query  interface{} `json:"query"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
}

func writeResult(code int, message string, data interface{}, w http.ResponseWriter) {
	w.WriteHeader(code)
	trs, _ := json.Marshal(Result{
		Code: code,
		Msg:  message,
		Data: data,
	})
	w.Write(trs)
}

// SetEndpoint method defines validator GRPC address.
func SetEndpoint(addr string) {
	endpoint = addr
}
