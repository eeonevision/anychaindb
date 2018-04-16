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
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/leadschain/leadschain/crypto"
	"gitlab.com/leadschain/leadschain/state"
)

// API is the high level interface for leadschain client applications
type API interface {
	AccountAPI
	TransitionAPI
	ConversionAPI
}

// AccountAPI describes all account related functions
type AccountAPI interface {
	CreateAccount() (id, pub, priv string, err error)
	GetAccount(id string) (*state.Account, error)
	DeleteAccount(id string) error
	ListAccounts() ([]*state.Account, error)
}

// TransitionAPI interface provides all transition related methods
type TransitionAPI interface {
	AddTransition(affiliateID, advertiserID, clickID, streamID, offerID string, expiresIn int64) (ID string, err error)
	GetTransition(ID string) (*state.Transition, error)
	ListTransitions() ([]*state.Transition, error)
	SearchTransitions(query []byte) ([]*state.Transition, error)
}

// ConversionAPI interface provides all conversion related methods
type ConversionAPI interface {
	AddConversion(affiliateID, advertiserID, clickID, streamID, clientID, goalID, offerID, status, comment string) (ID string, err error)
	GetConversion(ID string) (*state.Conversion, error)
	ListConversions() ([]*state.Conversion, error)
	SearchConversions(query []byte) ([]*state.Conversion, error)
}

// NewAPI constructs a new API instances based on an http transport
func NewAPI(endpoint string, key *crypto.Key, accountID string) API {
	base := NewHTTPClient(endpoint, key, accountID)
	fast := NewFastClient(endpoint, key, accountID)
	return &apiClient{endpoint, base, fast}
}

type apiClient struct {
	endpoint string
	base     *BaseClient
	fast     *FastClient
}

func (api *apiClient) CreateAccount() (id, pub, priv string, err error) {
	key, err := crypto.CreateKeyPair()
	if err != nil {
		return "", "", "", err
	}
	api.base.Key = key
	id = bson.NewObjectId().Hex()
	if err := api.fast.AddAccount(&state.Account{ID: id, PubKey: key.GetPubString()}); err != nil {
		return "", "", "", err
	}
	return id, key.GetPubString(), key.GetPrivString(), nil
}

func (api *apiClient) GetAccount(id string) (*state.Account, error) {
	return api.base.GetAccount(id)
}

func (api *apiClient) DeleteAccount(id string) error {
	return api.base.DelAccount(id)
}

func (api *apiClient) ListAccounts() ([]*state.Account, error) {
	return api.base.ListAccounts()
}

func (api *apiClient) AddTransition(affiliateID, advertiserID, clickID, streamID, offerID string, expiresIn int64) (ID string, err error) {
	id := bson.NewObjectId().Hex()
	now := time.Now()
	createdAt := now.UTC().UnixNano()
	if err := api.fast.AddTransition(&state.Transition{
		ID:                  id,
		AdvertiserAccountID: advertiserID,
		AffiliateAccountID:  affiliateID,
		ClickID:             clickID,
		StreamID:            streamID,
		OfferID:             offerID,
		CreatedAt:           float64(createdAt),
		ExpiresIn:           expiresIn,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (api *apiClient) GetTransition(id string) (*state.Transition, error) {
	return api.base.GetTransition(id)
}

func (api *apiClient) ListTransitions() ([]*state.Transition, error) {
	return api.base.ListTransitions()
}

func (api *apiClient) SearchTransitions(query []byte) ([]*state.Transition, error) {
	return api.base.SearchTransitions(query)
}

func (api *apiClient) AddConversion(affiliateID, advertiserID, clickID, streamID, clientID, goalID, offerID, status, comment string) (ID string, err error) {
	id := bson.NewObjectId().Hex()
	now := time.Now()
	createdAt := now.UTC().UnixNano()
	if err := api.fast.AddConversion(&state.Conversion{
		ID:                  id,
		AdvertiserAccountID: advertiserID,
		AffiliateAccountID:  affiliateID,
		ClickID:             clickID,
		StreamID:            streamID,
		ClientID:            clientID,
		OfferID:             offerID,
		GoalID:              goalID,
		CreatedAt:           float64(createdAt),
		Status:              status,
		Comment:             comment,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (api *apiClient) GetConversion(id string) (*state.Conversion, error) {
	return api.base.GetConversion(id)
}

func (api *apiClient) ListConversions() ([]*state.Conversion, error) {
	return api.base.ListConversions()
}

func (api *apiClient) SearchConversions(query []byte) ([]*state.Conversion, error) {
	return api.base.SearchConversions(query)
}
