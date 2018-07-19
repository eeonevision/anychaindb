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

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eeonevision/anychaindb/client"
	"github.com/eeonevision/anychaindb/crypto"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
)

// PrivateData keeps information about receiver and data,
// encrypted by receiver's public key
type PrivateData struct {
	ReceiverAccountID string      `json:"receiver_account_id" mapstructure:"receiver_account_id"`
	Data              interface{} `json:"data" mapstructure:"data"`
}

// Payload struct keeps transaction data related fields.
//   - PublicData keeps open data of any structure;
//   - PrivateData keeps encrypted by affiliate's public key with ECDH algorithm data and represented as base64 string;
//   - CreatedAt is date of object creation in UNIX time (milliseconds).
type Payload struct {
	ID              string         `json:"_id,omitempty" mapstructure:"_id"`
	SenderAccountID string         `json:"sender_account_id,omitempty" mapstructure:"sender_account_id"`
	PublicData      interface{}    `json:"public_data,omitempty" mapstructure:"public_data"`
	PrivateData     []*PrivateData `json:"private_data,omitempty" mapstructure:"private_data"`
	CreatedAt       float64        `json:"created_at,omitempty" mapstructure:"created_at"`
}

// PostPayloadsHandler uses FastAPI for sends new transaction data requests in async mode to blockchain.
func PostPayloadsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse form's JSON data
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		writeResult(http.StatusBadRequest, "request decode error: "+err.Error(), nil, w)
		return
	}

	var data Payload
	if err := mapstructure.Decode(req.Data, &data); err != nil {
		writeResult(http.StatusBadRequest, "payload decode error: "+err.Error(), nil, w)
		return
	}
	defer r.Body.Close()

	// Add payload to blockchain
	key, err := crypto.NewFromStrings(req.PubKey, req.PrivKey)
	if err != nil {
		writeResult(http.StatusUnauthorized, err.Error(), nil, w)
		return
	}
	api := client.NewAPI(endpoint, key, req.AccountID)
	privMrsh, err := json.Marshal(data.PrivateData)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	id, err := api.AddPayload(req.AccountID, data.PublicData, privMrsh)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	writeResult(http.StatusAccepted, "payload added", Payload{ID: id}, w)
	return
}

// GetPayloadsHandler uses BaseAPI for search and list transaction data.
// Query parameters: Query, Limit, Offset can be optional.
// Query - MongoDB query string.
// Limit - maximum 500 items.
// Offset - default 0.
func GetPayloadsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var query interface{}
	var limit int
	var offset int
	var err error

	// Get GET query params
	if q := r.URL.Query().Get("query"); q != "" {
		err := json.Unmarshal([]byte(q), &query)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"cannot parse query parameter: "+err.Error(), nil, w)
			return
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"cannot parse limit parameter: "+err.Error(), nil, w)
			return
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"cannot parse offset parameter: "+err.Error(), nil, w)
			return
		}
	}
	// Get basic auth data: receiver's account id and private key
	re, pk, _ := r.BasicAuth()

	// Check limits
	if limit > 500 || limit <= 0 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	api := client.NewAPI(endpoint, nil, "")
	searchReq := mongoQuery{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	searchReqStr, _ := json.Marshal(searchReq)
	cnv, err := api.SearchPayloads(searchReqStr, re, pk)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	writeResult(http.StatusOK, "OK", cnv, w)
	return
}

// GetPayloadDetailsHandler uses BaseAPI for get payload details by it id.
// Query parameters ID is required.
func GetPayloadDetailsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("id")
	if id == "" {
		writeResult(http.StatusBadRequest,
			"id should not be empty", nil, w)
		return
	}
	// Get basic auth data: receiver's account id and private key
	re, pk, _ := r.BasicAuth()

	api := client.NewAPI(endpoint, nil, "")
	cnv, err := api.GetPayload(id, re, pk)

	// Check special case when account not found
	// Temporary solution in case of introduce more right way of error handling
	if err.Error() == errNotFound.Error() {
		writeResult(http.StatusNotFound, err.Error(), nil, w)
		return
	}
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	writeResult(http.StatusOK, "OK", cnv, w)
	return
}
