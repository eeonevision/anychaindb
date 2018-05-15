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

package state

import (
	"errors"
)

//go:generate msgp

// Conversion struct keep conversion related fields.
type Conversion struct {
	ID                  string  `msg:"_id" json:"_id" mapstructure:"_id" bson:"_id"`
	AdvertiserAccountID string  `msg:"advertiser_account_id" json:"advertiser_account_id" mapstructure:"advertiser_account_id" bson:"advertiser_account_id"`
	AffiliateAccountID  string  `msg:"affiliate_account_id" json:"affiliate_account_id" mapstructure:"affiliate_account_id" bson:"affiliate_account_id"`
	ClickID             string  `msg:"click_id" json:"click_id" mapstructure:"click_id" bson:"click_id"`
	StreamID            string  `msg:"stream_id" json:"stream_id" mapstructure:"stream_id" bson:"stream_id"`
	OfferID             string  `msg:"offer_id" json:"offer_id" mapstructure:"offer_id" bson:"offer_id"`
	ClientID            string  `msg:"client_id" json:"client_id" mapstructure:"client_id" bson:"client_id"`
	GoalID              string  `msg:"goal_id" json:"goal_id" mapstructure:"goal_id" bson:"goal_id"`
	CreatedAt           float64 `msg:"created_at" json:"created_at" mapstructure:"created_at" bson:"created_at"`
	Comment             string  `msg:"comment" json:"comment" mapstructure:"comment" bson:"comment"`
	Status              string  `msg:"status" json:"status" mapstructure:"status" bson:"status"`
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
