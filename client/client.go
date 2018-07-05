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

	"github.com/globalsign/mgo/bson"
	"github.com/leadschain/anychaindb/crypto"
	"github.com/leadschain/anychaindb/state"
)

// API is the high level interface for Anychaindb client applications
type API interface {
	AccountAPI
	PayloadAPI
}

// AccountAPI describes all account related functions
type AccountAPI interface {
	CreateAccount() (id, pub, priv string, err error)
	GetAccount(id string) (*state.Account, error)
	ListAccounts() ([]state.Account, error)
	SearchAccounts(query []byte) ([]state.Account, error)
}

// PayloadAPI interface provides all transaction data related methods
type PayloadAPI interface {
	AddPayload(senderAccountID, receiverAccountID string, publicData, privateData []byte) (ID string, err error)
	GetPayload(ID string) (*state.Payload, error)
	ListPayloads() ([]state.Payload, error)
	SearchPayloads(query []byte) ([]state.Payload, error)
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

func (api *apiClient) AddPayload(senderAccountID, receiverAccountID string, publicData, privateData []byte) (ID string, err error) {
	id := bson.NewObjectId().Hex()
	now := time.Now()
	createdAt := now.UTC().UnixNano()

	// Get receiver's public key
	receiver, err := api.GetAccount(receiverAccountID)
	if err != nil {
		return "", err
	}

	// ECDH encrypted private data with public key of receiver
	receiverPubKey, err := crypto.NewFromStrings(receiver.PubKey, "")
	if err != nil {
		return "", err
	}
	privateDataEnc, err := receiverPubKey.Encrypt(privateData)
	if err != nil {
		return "", err
	}

	if err = api.fast.AddPayload(&state.Payload{
		ID:                id,
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
		PublicData:        string(publicData),
		PrivateData:       base64.StdEncoding.EncodeToString(privateDataEnc),
		CreatedAt:         float64(createdAt),
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (api *apiClient) GetPayload(id string) (*state.Payload, error) {
	return api.base.GetPayload(id)
}

func (api *apiClient) ListPayloads() ([]state.Payload, error) {
	return api.base.ListPayloads()
}

func (api *apiClient) SearchPayloads(query []byte) ([]state.Payload, error) {
	return api.base.SearchPayloads(query)
}
