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
	"encoding/base64"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/globalsign/mgo/bson"
	"github.com/leadschain/leadschain/crypto"
	"github.com/leadschain/leadschain/state"
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
	ListAccounts() ([]state.Account, error)
	SearchAccounts(query []byte) ([]state.Account, error)
}

// TransitionAPI interface provides all transition related methods
type TransitionAPI interface {
	AddTransition(affiliateID, advertiserID, clickID, streamID, offerID string, expiresIn int64) (ID string, err error)
	GetTransition(ID string) (*state.Transition, error)
	ListTransitions() ([]state.Transition, error)
	SearchTransitions(query []byte) ([]state.Transition, error)
}

// ConversionAPI interface provides all conversion related methods
type ConversionAPI interface {
	AddConversion(affiliateID string, advertiserData, publicData []byte, status string) (ID string, err error)
	GetConversion(ID string) (*state.Conversion, error)
	ListConversions() ([]state.Conversion, error)
	SearchConversions(query []byte) ([]state.Conversion, error)
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

func (api *apiClient) ListAccounts() ([]state.Account, error) {
	return api.base.ListAccounts()
}

func (api *apiClient) SearchAccounts(query []byte) ([]state.Account, error) {
	return api.base.SearchAccounts(query)
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

func (api *apiClient) ListTransitions() ([]state.Transition, error) {
	return api.base.ListTransitions()
}

func (api *apiClient) SearchTransitions(query []byte) ([]state.Transition, error) {
	return api.base.SearchTransitions(query)
}

func (api *apiClient) AddConversion(affiliateID string, advertiserData, publicData []byte, status string) (ID string, err error) {
	id := bson.NewObjectId().Hex()
	now := time.Now()
	createdAt := now.UTC().UnixNano()

	// Get affiliate's public key
	affiliate, err := api.GetAccount(affiliateID)
	if err != nil {
		return "", err
	}

	// ECDH encrypted advertiser's data with public key of affiliate
	affiliatePubKey, err := crypto.NewFromStrings(affiliate.PubKey, "")
	if err != nil {
		return "", err
	}
	advertiserDataEnc, err := affiliatePubKey.Encrypt(advertiserData)
	if err != nil {
		return "", err
	}

	// BLAKE2B 256-bit hashed public data
	blake2BHash, _ := blake2b.New256(nil)
	_, err = blake2BHash.Write(publicData)
	if err != nil {
		return "", err
	}

	if err = api.fast.AddConversion(&state.Conversion{
		ID:                 id,
		AffiliateAccountID: affiliateID,
		AdvertiserData:     base64.StdEncoding.EncodeToString(advertiserDataEnc),
		PublicData:         base64.StdEncoding.EncodeToString(blake2BHash.Sum(nil)),
		CreatedAt:          float64(createdAt),
		Status:             status,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (api *apiClient) GetConversion(id string) (*state.Conversion, error) {
	return api.base.GetConversion(id)
}

func (api *apiClient) ListConversions() ([]state.Conversion, error) {
	return api.base.ListConversions()
}

func (api *apiClient) SearchConversions(query []byte) ([]state.Conversion, error) {
	return api.base.SearchConversions(query)
}
