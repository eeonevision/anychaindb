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

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/leadschain/leadschain/client"
	"gitlab.com/leadschain/leadschain/crypto"
)

// Transition struct describes transition related fields
// All fields, except AdvertiserAccountID and AffiliateAccountID are optional.
//
// AdvertiserAccountID - unique identifier of advertiser account
// AffiliateAccountID - unique identifier of affiliate account
// ClickID - unique identifier of user's click (usually transmitted from affiliate network)
// StreamID (a.k.a. PlatformID) - uniques identifier of webmaster's platform in afiliate network
// OfferID - unique identifier of offer in affiliate network
// ExpiresIn - Unix timestamp of transition expiration (in seconds)
type Transition struct {
	AdvertiserAccountID string `json:"advertiser_account_id" mapstructure:"advertiser_account_id"`
	AffiliateAccountID  string `json:"affiliate_account_id" mapstructure:"affiliate_account_id"`
	ClickID             string `json:"click_id" mapstructure:"click_id"`
	StreamID            string `json:"stream_id" mapstructure:"stream_id"`
	OfferID             string `json:"offer_id" mapstructure:"offer_id"`
	ExpiresIn           int64  `json:"expires_in" mapstructure:"expires_in"`
}

// PostTransitionsHandler uses FastAPI for sends new transitions requests in async mode to blockchain
func PostTransitionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse form's JSON data
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		writeResult(http.StatusBadRequest, "Request decode error: "+err.Error(), nil, w)
		return
	}

	var data Transition
	if err := mapstructure.Decode(req.Data, &data); err != nil {
		writeResult(http.StatusBadRequest, "Transition decode error: "+err.Error(), nil, w)
		return
	}
	defer r.Body.Close()

	// Add transition to blockchain
	key, err := crypto.NewFromStrings(req.PubKey, req.PrivKey)
	if err != nil {
		writeResult(http.StatusUnauthorized, err.Error(), nil, w)
		return
	}
	api := client.NewAPI(endpoint, key, req.AccountID)
	id, err := api.AddTransition(data.AffiliateAccountID, data.AdvertiserAccountID, data.ClickID, data.StreamID,
		data.OfferID, data.ExpiresIn)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	// OK
	writeResult(http.StatusAccepted, "Transition added", id, w)
}

// GetTransitionsHandler uses FastAPI for search and list transitions.
// Query parameters: Query, Limit, Offset can be optional.
// Query - MongoDB query string.
// Limit - maximum 500 items.
// Offset - default 0.
func GetTransitionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var query *string
	var limit int
	var offset int
	var err error

	if q := r.URL.Query().Get("query"); q != "" {
		query = &q
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
	trs, err := api.SearchTransitions(searchReqStr)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", trs, w)
	return
}

// GetTransitionDetailsHandler uses FastAPI for get transition details by it id.
// Query parameters ID is required.
func GetTransitionDetailsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("id")
	if id == "" {
		writeResult(http.StatusBadRequest,
			"ID should not be empty", nil, w)
		return
	}
	api := client.NewAPI(endpoint, nil, "")
	trs, err := api.GetTransition(id)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", trs, w)
	return
}
