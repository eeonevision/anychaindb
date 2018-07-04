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

// Conversion struct keep conversion related fields.
//   - PrivateData keeps cpa_uid, client_id, goal_id, comment and some other relevant to postback private data.
//     Encrypted by affiliate's public key by ECDH algorithm and represented as base64 string.
//   - PublicData keeps offer_id, stream_id, advertiser_account_id and affiliate's public key
//     to provide possibility for transaction proving by affiliate. Encrypted by BLAKE2B 256bit hash,
//     represented as base64 string.
type Conversion struct {
	ID                 string  `msg:"_id" json:"_id" mapstructure:"_id" bson:"_id"`
	AffiliateAccountID string  `msg:"affiliate_account_id" json:"affiliate_account_id" mapstructure:"affiliate_account_id" bson:"affiliate_account_id"`
	PrivateData        string  `msg:"private_data" json:"private_data" mapstructure:"private_data" bson:"private_data"`
	PublicData         string  `msg:"public_data" json:"public_data" mapstructure:"public_data" bson:"public_data"`
	CreatedAt          float64 `msg:"created_at" json:"created_at" mapstructure:"created_at" bson:"created_at"`
}

const conversionsCollection = "conversions"

// AddConversion method adds new conversion to the state if it not exists.
func (s *State) AddConversion(conversion *Conversion) error {
	if s.HasConversion(conversion.ID) {
		return errors.New("Conversion exists")
	}
	return s.SetConversion(conversion)
}

// SetConversion inserts new conversion to state without any checks.
func (s *State) SetConversion(conversion *Conversion) error {
	return s.DB.C(conversionsCollection).Insert(conversion)
}

// HasConversion method checks exists conversion in state ot not.
func (s *State) HasConversion(id string) bool {
	if res, _ := s.GetConversion(id); res != nil {
		return true
	}
	return false
}

// GetConversion method gets conversion from state by it identifier.
func (s *State) GetConversion(id string) (*Conversion, error) {
	var result *Conversion
	return result, s.DB.C(conversionsCollection).FindId(id).One(&result)
}

// ListConversions method returns list of all conversions in state.
func (s *State) ListConversions() (result []*Conversion, err error) {
	return result, s.DB.C(conversionsCollection).Find(nil).All(&result)
}

// SearchConversions method finds conversions using mongodb query language.
func (s *State) SearchConversions(query interface{}, limit, offset int) (result []*Conversion, err error) {
	return result, s.DB.C(conversionsCollection).Find(query).Skip(offset).Limit(limit).All(&result)
}
