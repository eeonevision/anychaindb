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
	"errors"

	"github.com/tendermint/tendermint/rpc/core/types"

	"github.com/eeonevision/anychaindb/crypto"
	"github.com/eeonevision/anychaindb/state"
	"github.com/eeonevision/anychaindb/transaction"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

// baseClient struct contains config
// parameters for performing requests.
type baseClient struct {
	key       *crypto.Key
	endpoint  string
	mode      string
	accountID string
	tm        client.Client
}

// newHTTPClient initializes new base client instance.
func newHTTPClient(endpoint, mode string, key *crypto.Key, accountID string) *baseClient {
	// Set default mode
	switch mode {
	case "sync":
		mode = "sync"
		break
	case "async":
		mode = "async"
		break
	case "commit":
		mode = "commit"
		break
	default:
		mode = "sync"
	}
	tm := client.NewHTTP(endpoint, "/websocket")
	return &baseClient{key, endpoint, mode, accountID, tm}
}

func (c *baseClient) addAccount(acc *state.Account) error {
	var err error

	txBytes, err := acc.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.AccountAdd, c.accountID, txBytes)
	bs, _ := tx.ToBytes()

	return c.broadcastTx(bs)
}

func (c *baseClient) getAccount(id string) (*state.Account, error) {
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

func (c *baseClient) searchAccounts(searchQuery []byte) ([]state.Account, error) {
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

func (c *baseClient) addPayload(cv *state.Payload) error {
	txBytes, err := cv.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.PayloadAdd, c.accountID, txBytes)
	if err := tx.Sign(c.key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()

	return c.broadcastTx(bs)
}

func (c *baseClient) getPayload(id string) (*state.Payload, error) {
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

func (c *baseClient) searchPayloads(searchQuery []byte) ([]state.Payload, error) {
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

func (c *baseClient) broadcastTx(bs []byte) error {
	var res interface{}
	var err error

	switch c.mode {
	case "async":
		res, err = c.tm.BroadcastTxAsync(types.Tx(bs))
	case "sync":
		res, err = c.tm.BroadcastTxSync(types.Tx(bs))
	case "commit":
		res, err = c.tm.BroadcastTxCommit(types.Tx(bs))
	}
	// Check transport errors
	if err != nil {
		return err
	}
	// Check special empty case
	if res == nil {
		return errors.New("empty response")
	}
	// Check for async/sync response
	if r, ok := res.(*core_types.ResultBroadcastTx); ok && r.Code != 0 {
		return errors.New(r.Log)
	}
	// Check for commit response
	if r, ok := res.(*core_types.ResultBroadcastTxCommit); ok && (r.CheckTx.Code != 0 || r.DeliverTx.Code != 0) {
		return errors.New("check tx error: " + r.CheckTx.Log + "; deliver tx error: " + r.DeliverTx.Log)
	}

	return nil
}
