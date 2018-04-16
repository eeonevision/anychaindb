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

package app

import (
	"errors"

	"github.com/leadschain/leadschain/state"
	"github.com/leadschain/leadschain/transaction"
)

func checkTransitionAddTransaction(tx *transaction.Transaction, s *state.State) error {
	data := &state.Transition{}
	_, err := data.UnmarshalMsg(tx.Data)
	if err != nil {
		return err
	}
	if s.HasTransition(data.ID) {
		return errors.New("Transition exists")
	}
	k, err := s.GetAccountPubKey(data.AdvertiserAccountID)
	if err != nil {
		return errors.New("Pubkey can't be loaded: " + err.Error())
	}
	if err := tx.Verify(k); err != nil {
		return errors.New("TX can't be verified: " + err.Error())
	}
	return nil
}

func deliverTransitionAddTransaction(tx *transaction.Transaction, s *state.State) error {
	data := &state.Transition{}
	_, err := data.UnmarshalMsg(tx.Data)
	if err != nil {
		return err
	}
	return s.AddTransition(data)
}
