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

package state

import (
	"errors"
)

//go:generate msgp

// PrivateData keeps information about receiver and data,
// encrypted by receiver's public key
type PrivateData struct {
	ReceiverAccountID string      `msg:"receiver_account_id" json:"receiver_account_id" mapstructure:"receiver_account_id" bson:"receiver_account_id"`
	Data              interface{} `msg:"data" json:"data" mapstructure:"data" bson:"data"`
}

// Payload struct keeps transaction data related fields.
//   - PublicData keeps open data of any structure;
//   - PrivateData keeps encrypted data set by receiver's public key with ECDH algorithm and represented as base64 string;
//   - CreatedAt is date of object creation in UNIX time (milliseconds).
type Payload struct {
	ID              string         `msg:"_id" json:"_id" mapstructure:"_id" bson:"_id"`
	SenderAccountID string         `msg:"sender_account_id" json:"sender_account_id" mapstructure:"sender_account_id" bson:"sender_account_id"`
	PublicData      interface{}    `msg:"public_data" json:"public_data" mapstructure:"public_data" bson:"public_data"`
	PrivateData     []*PrivateData `msg:"private_data" json:"private_data" mapstructure:"private_data" bson:"private_data"`
	CreatedAt       float64        `msg:"created_at" json:"created_at" mapstructure:"created_at" bson:"created_at"`
}

const payloadsCollection = "data"

// AddPayload method adds new payload to the state if it not exists.
func (s *State) AddPayload(data *Payload) error {
	if s.HasPayload(data.ID) {
		return errors.New("payload exists")
	}
	return s.SetPayload(data)
}

// SetPayload inserts new payload to state without any checks.
func (s *State) SetPayload(data *Payload) error {
	return s.DB.C(payloadsCollection).Insert(data)
}

// HasPayload method checks exists payload in state ot not.
func (s *State) HasPayload(id string) bool {
	if res, _ := s.GetPayload(id); res != nil {
		return true
	}
	return false
}

// GetPayload method gets data from state by it identifier.
func (s *State) GetPayload(id string) (*Payload, error) {
	var result *Payload
	return result, s.DB.C(payloadsCollection).FindId(id).One(&result)
}

// SearchPayloads method finds payloads using mongodb query language.
func (s *State) SearchPayloads(query interface{}, limit, offset int) (result []*Payload, err error) {
	return result, s.DB.C(payloadsCollection).Find(query).Skip(offset).Limit(limit).All(&result)
}
