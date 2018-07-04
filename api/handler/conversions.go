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

	"github.com/julienschmidt/httprouter"
	"github.com/leadschain/leadschain/client"
	"github.com/leadschain/leadschain/crypto"
	"github.com/mitchellh/mapstructure"
)

// Conversion struct describes conversion related fields
//   - PrivateData keeps cpa_uid, client_id, goal_id, comment, status and some other relevant to postback private data.
//     Encrypted by affiliate's public key by ECDH algorithm and represented as base64 string.
//   - PublicData keeps offer_id, stream_id, status, advertiser_account_id and affiliate's public key
//     to provide possibility for transaction proving by affiliate. Encrypted by BLAKE2B 256bit hash.
type Conversion struct {
	ID                 string      `msg:"_id" json:"_id" mapstructure:"_id" bson:"_id"`
	AffiliateAccountID string      `msg:"affiliate_account_id" json:"affiliate_account_id" mapstructure:"affiliate_account_id" bson:"affiliate_account_id"`
	PrivateData        interface{} `msg:"private_data" json:"private_data" mapstructure:"private_data" bson:"private_data"`
	PublicData         interface{} `msg:"public_data" json:"public_data" mapstructure:"public_data" bson:"public_data"`
	CreatedAt          float64     `msg:"created_at" json:"created_at" mapstructure:"created_at" bson:"created_at"`
}

// PostConversionsHandler uses FastAPI for sends new conversions requests in async mode to blockchain
func PostConversionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse form's JSON data
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		writeResult(http.StatusBadRequest, "Request decode error: "+err.Error(), nil, w)
		return
	}

	var data Conversion
	if err := mapstructure.Decode(req.Data, &data); err != nil {
		writeResult(http.StatusBadRequest, "Conversion decode error: "+err.Error(), nil, w)
		return
	}
	defer r.Body.Close()

	// Add conversion to blockchain
	key, err := crypto.NewFromStrings(req.PubKey, req.PrivKey)
	if err != nil {
		writeResult(http.StatusUnauthorized, err.Error(), nil, w)
		return
	}
	api := client.NewAPI(endpoint, key, req.AccountID)
	advMrsh, err := json.Marshal(data.PrivateData)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	pubMrsh, err := json.Marshal(data.PublicData)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	id, err := api.AddConversion(data.AffiliateAccountID, advMrsh, pubMrsh)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	// OK
	writeResult(http.StatusAccepted, "Conversion added", Conversion{ID: id}, w)
}

// GetConversionsHandler uses BaseAPI for search and list conversions.
// Query parameters: Query, Limit, Offset can be optional.
// Query - MongoDB query string.
// Limit - maximum 500 items.
// Offset - default 0.
func GetConversionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var query interface{}
	var limit int
	var offset int
	var err error

	if q := r.URL.Query().Get("query"); q != "" {
		err := json.Unmarshal([]byte(q), &query)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"Cannot parse query parameter: "+err.Error(), nil, w)
			return
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"Cannot parse limit parameter: "+err.Error(), nil, w)
			return
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			writeResult(http.StatusBadRequest,
				"Cannot parse offset parameter: "+err.Error(), nil, w)
			return
		}
	}
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
	cnv, err := api.SearchConversions(searchReqStr)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	// Empty conversions list
	if cnv == nil {
		writeResult(http.StatusNotFound, "Empty list", nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", cnv, w)
	return
}

// GetConversionDetailsHandler uses BaseAPI for get conversion details by it id.
// Query parameters ID is required.
func GetConversionDetailsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("id")
	if id == "" {
		writeResult(http.StatusBadRequest,
			"ID should not be empty", nil, w)
		return
	}
	api := client.NewAPI(endpoint, nil, "")
	cnv, err := api.GetConversion(id)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	// Conversion not found
	if cnv == nil {
		writeResult(http.StatusNotFound, "Not Found", nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", cnv, w)
	return
}
