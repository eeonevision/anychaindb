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

package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/eeonevision/anychaindb/crypto"
	"github.com/eeonevision/anychaindb/state"
	"github.com/eeonevision/anychaindb/transaction"
)

type FastClient struct {
	Endpoint  string
	Key       *crypto.Key
	AccountID string
	client    *http.Client
}

func NewFastClient(endpoint string, key *crypto.Key, accountID string) *FastClient {
	tm := &http.Client{Timeout: 30 * time.Second}
	return &FastClient{endpoint, key, accountID, tm}
}

func (c *FastClient) BroadcastTxAsync(tx []byte) (*http.Response, error) {
	ba := base64.StdEncoding.EncodeToString(tx)
	txData := fmt.Sprintf(`{"jsonrpc":"2.0","id":"anything","method":"broadcast_tx_async","params": {"tx": "%s"}}`, ba)
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer([]byte(txData)))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (c *FastClient) AddAccount(acc *state.Account) error {
	txBytes, err := acc.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.AccountAdd, c.AccountID, txBytes)
	bs, _ := tx.ToBytes()
	res, err := c.BroadcastTxAsync(bs)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%v: %s", res.StatusCode, res.Status)
	}
	return nil
}

func (c *FastClient) AddPayload(data *state.Payload) error {
	txBytes, err := data.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.PayloadAdd, c.AccountID, txBytes)
	if err := tx.Sign(c.Key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()
	res, err := c.BroadcastTxAsync(bs)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%v: %s", res.StatusCode, res.Status)
	}
	return nil
}
