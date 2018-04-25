/*
 * Copyright (C) 2018 Leads Studio
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

	"github.com/leadschain/leadschain/crypto"
	"github.com/leadschain/leadschain/state"
	"github.com/leadschain/leadschain/transaction"
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
	tx := transaction.New(transaction.AccountAdd, txBytes)
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

func (c *BaseClient) AddTransition(tr *state.Transition) error {
	txBytes, err := tr.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.TransitionAdd, txBytes)
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

func (c *BaseClient) GetTransition(id string) (*state.Transition, error) {
	resp, err := c.tm.ABCIQuery("transitions", []byte(id))
	if err != nil {
		return nil, err
	}
	tr := &state.Transition{}
	if err := json.Unmarshal(resp.Response.GetValue(), &tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *BaseClient) ListTransitions() ([]state.Transition, error) {
	resp, err := c.tm.ABCIQuery("transitions", nil)
	if err != nil {
		return nil, err
	}
	tr := []state.Transition{}
	if err := json.Unmarshal(resp.Response.GetValue(), &tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *BaseClient) SearchTransitions(searchQuery []byte) ([]state.Transition, error) {
	resp, err := c.tm.ABCIQuery("transitions/search", searchQuery)
	if err != nil {
		return nil, err
	}
	tr := []state.Transition{}
	if err := json.Unmarshal(resp.Response.GetValue(), &tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *BaseClient) AddConversion(cv *state.Conversion) error {
	txBytes, err := cv.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.ConversionAdd, txBytes)
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

func (c *BaseClient) GetConversion(id string) (*state.Conversion, error) {
	resp, err := c.tm.ABCIQuery("conversions", []byte(id))
	if err != nil {
		return nil, err
	}
	cv := &state.Conversion{}
	if err := json.Unmarshal(resp.Response.GetValue(), &cv); err != nil {
		return nil, err
	}
	fmt.Println(cv)
	return cv, nil
}

func (c *BaseClient) ListConversions() ([]state.Conversion, error) {
	resp, err := c.tm.ABCIQuery("conversions", nil)
	if err != nil {
		return nil, err
	}
	cv := []state.Conversion{}
	if err := json.Unmarshal(resp.Response.GetValue(), &cv); err != nil {
		return nil, err
	}
	return cv, nil
}

func (c *BaseClient) SearchConversions(searchQuery []byte) ([]state.Conversion, error) {
	resp, err := c.tm.ABCIQuery("conversions/search", searchQuery)
	if err != nil {
		return nil, err
	}
	cv := []state.Conversion{}
	if err := json.Unmarshal(resp.Response.GetValue(), &cv); err != nil {
		return nil, err
	}
	return cv, nil
}
