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

package transaction

import (
	"math/rand"
	"time"

	"github.com/tinylib/msgp/msgp"

	"github.com/leadschain/leadschain/crypto"

	"golang.org/x/crypto/sha3"
)

//go:generate msgp

type Transaction struct {
	Type      TransactionType `msg:"type" json:"type"`
	Timestamp int64           `msg:"timestamp" json:"timestamp"`
	Signer    string          `msg:"signer" json:"signer"`
	Signature string          `msg:"signature" json:"signature"`
	Nonce     uint32          `msg:"nonce" json:"nonce"`
	Data      []byte          `msg:"data" json:"data"`
}

type TransactionType string

const (
	AccountAdd    TransactionType = "add-account"
	TransitionAdd TransactionType = "add-transition"
	ConversionAdd TransactionType = "add-conversion"
)

func (t *Transaction) FromBytes(bs []byte) error {
	_, err := t.UnmarshalMsg(bs)
	return err
}

func (t *Transaction) ToBytes() ([]byte, error) {
	return t.MarshalMsg(nil)
}

func (t *Transaction) Hash() []byte {
	hash := sha3.New512()
	w := msgp.NewWriter(hash)
	w.WriteString(string(t.Type))
	w.WriteInt64(t.Timestamp)
	w.WriteString(t.Signer)
	w.WriteUint32(t.Nonce)
	w.WriteBytes(t.Data)
	w.Flush()
	return hash.Sum(nil)
}

func (t *Transaction) Sign(key *crypto.Key) error {
	hash := t.Hash()
	signature, err := key.Sign(hash)
	if err != nil {
		return err
	}
	t.Signature = signature
	return nil
}

func (t *Transaction) Verify(key *crypto.Key) error {
	hash := t.Hash()
	return key.Verify(hash, t.Signature)
}

func New(t TransactionType, signer string, data []byte) *Transaction {
	now := time.Now().UnixNano()
	rand.Seed(now)
	return &Transaction{t, now, signer, "", rand.Uint32(), data}
}
