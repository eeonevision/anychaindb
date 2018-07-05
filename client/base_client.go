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
	"encoding/json"
	"fmt"

	"github.com/leadschain/anychaindb/crypto"
	"github.com/leadschain/anychaindb/state"
	"github.com/leadschain/anychaindb/transaction"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

type BaseClient struct {
	Key       *crypto.Key
	AccountID string
	tm        client.Client
}

func NewHTTPClient(endpoint string, key *crypto.Key, accountID string) *BaseClient {
	tm := client.NewHTTP(endpoint, "/websocket")
	return &BaseClient{key, accountID, tm}
}

func (c *BaseClient) AddAccount(acc *state.Account) error {
	txBytes, err := acc.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.AccountAdd, c.AccountID, txBytes)
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxAsync(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("%v: %s", res.Code, res.Log)
	}
	return nil
}

func (c *BaseClient) GetAccount(id string) (*state.Account, error) {
	resp, err := c.tm.ABCIQuery("accounts", []byte(id))
	if err != nil {
		return nil, err
	}
	acc := &state.Account{}
	if err := json.Unmarshal(resp.Response.GetValue(), &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *BaseClient) ListAccounts() ([]state.Account, error) {
	resp, err := c.tm.ABCIQuery("accounts", nil)
	if err != nil {
		return nil, err
	}
	acc := []state.Account{}
	if err := json.Unmarshal(resp.Response.GetValue(), &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *BaseClient) SearchAccounts(searchQuery []byte) ([]state.Account, error) {
	resp, err := c.tm.ABCIQuery("accounts/search", searchQuery)
	if err != nil {
		return nil, err
	}
	acc := []state.Account{}
	if err := json.Unmarshal(resp.Response.GetValue(), &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *BaseClient) AddPayload(cv *state.Payload) error {
	txBytes, err := cv.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.PayloadAdd, c.AccountID, txBytes)
	if err := tx.Sign(c.Key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()
	res, err := c.tm.BroadcastTxAsync(types.Tx(bs))
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("%v: %s", res.Code, res.Log)
	}
	return nil
}

func (c *BaseClient) GetPayload(id string) (*state.Payload, error) {
	resp, err := c.tm.ABCIQuery("payloads", []byte(id))
	if err != nil {
		return nil, err
	}
	res := &state.Payload{}
	if err := json.Unmarshal(resp.Response.GetValue(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *BaseClient) ListPayloads() ([]state.Payload, error) {
	resp, err := c.tm.ABCIQuery("payloads", nil)
	if err != nil {
		return nil, err
	}
	res := []state.Payload{}
	if err := json.Unmarshal(resp.Response.GetValue(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *BaseClient) SearchPayloads(searchQuery []byte) ([]state.Payload, error) {
	resp, err := c.tm.ABCIQuery("payloads/search", searchQuery)
	if err != nil {
		return nil, err
	}
	res := []state.Payload{}
	if err := json.Unmarshal(resp.Response.GetValue(), &res); err != nil {
		return nil, err
	}
	return res, nil
}
