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
	"encoding/json"
	"errors"
	"time"

	"github.com/eeonevision/anychaindb/crypto"
	"github.com/eeonevision/anychaindb/state"
	"github.com/globalsign/mgo/bson"
)

// API is the high level interface for Anychaindb client applications.
type API interface {
	AccountAPI
	PayloadAPI
}

// AccountAPI describes all account related functions.
type AccountAPI interface {
	CreateAccount() (id, pub, priv string, err error)
	GetAccount(id string) (*state.Account, error)
	SearchAccounts(query []byte) ([]state.Account, error)
}

// PayloadAPI interface provides all transaction data related methods.
type PayloadAPI interface {
	AddPayload(senderAccountID string, publicData interface{}, privateData []byte) (ID string, err error)
	GetPayload(ID, receiverID, privKey string) (*state.Payload, error)
	SearchPayloads(query []byte, receiverID, privKey string) ([]state.Payload, error)
}

// NewAPI constructs a new API instances based on an http transport.
func NewAPI(endpoint, mode string, key *crypto.Key, accountID string) API {
	base := newHTTPClient(endpoint, mode, key, accountID)
	fast := newFastClient(endpoint, mode, key, accountID)
	return &apiClient{endpoint, mode, base, fast}
}

type apiClient struct {
	endpoint string
	mode     string
	base     *baseClient
	fast     *fastClient
}

func (api *apiClient) CreateAccount() (id, pub, priv string, err error) {
	key, err := crypto.CreateKeyPair()
	if err != nil {
		return "", "", "", err
	}
	api.base.key = key
	id = bson.NewObjectId().Hex()
	err = api.base.addAccount(&state.Account{ID: id, PubKey: key.GetPubString()})
	if err != nil {
		return "", "", "", err
	}
	return id, key.GetPubString(), key.GetPrivString(), nil
}

func (api *apiClient) GetAccount(id string) (*state.Account, error) {
	return api.fast.getAccount(id)
}

func (api *apiClient) SearchAccounts(query []byte) ([]state.Account, error) {
	return api.base.searchAccounts(query)
}

func (api *apiClient) AddPayload(senderAccountID string, publicData interface{}, privateData []byte) (ID string, err error) {
	id := bson.NewObjectId().Hex()

	// Unmarshal private data
	var privData []*state.PrivateData
	err = json.Unmarshal(privateData, &privData)
	if err != nil {
		return "", errors.New("error in unmarshalling private data: " + err.Error())
	}

	for _, data := range privData {
		// Get receiver's public key
		receiver, err := api.fast.getAccount(data.ReceiverAccountID)
		if err != nil {
			return "", errors.New(
				"error in getting receiver's account " + data.ReceiverAccountID + ": " + err.Error(),
			)
		}
		// Check special case if account is empty
		if receiver == nil {
			return "", errors.New(
				"error in getting receiver's account " + data.ReceiverAccountID + ": " + "empty result",
			)
		}
		receiverPubKey, err := crypto.NewFromStrings(receiver.PubKey, "")
		if err != nil {
			return "", errors.New("error in processing receiver's public key: " + err.Error())
		}
		// Marshal private data of receiver
		privMrsh, err := json.Marshal(data.Data)
		if err != nil {
			return "", errors.New("error in marshalling private data: " + err.Error())
		}
		// ECDH encrypted private data with public key of receiver
		privateDataEnc, err := receiverPubKey.Encrypt(privMrsh)
		if err != nil {
			return "", errors.New("error in encrypting private data: " + err.Error())
		}
		// Reassign data from raw to encrypted and base64 encoded string
		data.Data = base64.StdEncoding.EncodeToString(privateDataEnc)
	}
	err = api.fast.addPayload(&state.Payload{
		ID:              id,
		SenderAccountID: senderAccountID,
		PublicData:      publicData,
		PrivateData:     privData,
		CreatedAt:       float64(time.Now().UnixNano() / 1000000),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (api *apiClient) GetPayload(id, receiverID, privKey string) (*state.Payload, error) {
	payload, err := api.base.getPayload(id)
	if err != nil {
		return payload, err
	}
	// Check if decoding not needed
	if receiverID == "" && privKey == "" {
		return payload, nil
	}
	// Check if payload result is empty
	if payload == nil {
		return payload, nil
	}
	res, err := api.decryptPrivateData(receiverID, privKey, []state.Payload{*payload})
	return &res[0], err
}

func (api *apiClient) SearchPayloads(query []byte, receiverID, privKey string) ([]state.Payload, error) {
	payloads, err := api.base.searchPayloads(query)
	if err != nil {
		return payloads, err
	}
	// Check if decoding not needed
	if receiverID == "" && privKey == "" {
		return payloads, nil
	}
	// Check if payload result is empty
	if len(payloads) == 0 {
		return payloads, nil
	}
	// Decrypt private data
	return api.decryptPrivateData(receiverID, privKey, payloads)
}

func (api *apiClient) decryptPrivateData(receiverID, privKey string, payloads []state.Payload) ([]state.Payload, error) {
	// Get account's public key
	acc, err := api.fast.getAccount(receiverID)
	if err != nil {
		return payloads, errors.New("invalid authorization account id: " + err.Error())
	}
	// Set private key structure for account
	pK, err := crypto.NewFromStrings(acc.PubKey, privKey)
	if err != nil {
		return payloads, errors.New("invalid authorization private key: " + err.Error())
	}
	// Decrypt all data signed by public key of receiver
	for _, payload := range payloads {
		for _, p := range payload.PrivateData {
			if p.ReceiverAccountID == receiverID {
				decodedBin, _ := base64.StdEncoding.DecodeString(p.Data.(string))
				decryptedBin, err := pK.Decrypt(decodedBin)
				if err != nil {
					return payloads, errors.New(
						"cannot decrypt private data for receiver's account id " + p.ReceiverAccountID + ": " + err.Error(),
					)
				}
				var decrypted interface{}
				err = json.Unmarshal(decryptedBin, &decrypted)
				if err != nil {
					return payloads, err
				}
				p.Data = decrypted
			}
		}
	}
	return payloads, nil
}
