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
	"github.com/leadschain/leadschain/client"
	"github.com/leadschain/leadschain/crypto"
	"github.com/mitchellh/mapstructure"
)

// Conversion struct describes conversion related fields
// All fields, except AdvertiserAccountID and AffiliateAccountID are optional.
//
// AdvertiserAccountID - unique identifier of advertiser account
// AffiliateAccountID - unique identifier of affiliate account
// ClickID - unique identifier of user's click (usually transmitted from affiliate network)
// StreamID (a.k.a. PlatformID) - uniques identifier of webmaster's platform in afiliate network
// OfferID - unique identifier of offer in affiliate network
// ClientID - unique identifer of client in CRM (database) of advertiser
// GoalID - type of conversion in affiliate network (maybe any, like: "per lead, filled form or so on...")
// Comment - commentary from advertiser for current conversion
// Status - status of conversion (PENDING, DECLINED, CONFIRMED)
type Conversion struct {
	AdvertiserAccountID string `json:"advertiser_account_id" mapstructure:"advertiser_account_id" bson:"advertiser_account_id"`
	AffiliateAccountID  string `json:"affiliate_account_id" mapstructure:"affiliate_account_id" bson:"affiliate_account_id"`
	ClickID             string `json:"click_id" mapstructure:"click_id" bson:"click_id"`
	StreamID            string `json:"stream_id" mapstructure:"stream_id" bson:"stream_id"`
	OfferID             string `json:"offer_id" mapstructure:"offer_id" bson:"offer_id"`
	ClientID            string `json:"client_id" mapstructure:"client_id" bson:"client_id"`
	GoalID              string `json:"goal_id" mapstructure:"goal_id" bson:"goal_id"`
	Comment             string `json:"comment" mapstructure:"comment,omitempty" bson:"comment"`
	Status              string `json:"status" mapstructure:"status" bson:"status"`
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
	id, err := api.AddConversion(data.AffiliateAccountID, data.AdvertiserAccountID, data.ClickID, data.StreamID,
		data.ClientID, data.GoalID, data.OfferID, data.Status, data.Comment)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	// OK
	writeResult(http.StatusAccepted, "Conversion added", id, w)
}

// GetConversionsHandler uses FastAPI for search and list conversions.
// Query parameters: Query, Limit, Offset can be optional.
// Query - MongoDB query string.
// Limit - maximum 500 items.
// Offset - default 0.
func GetConversionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	trs, err := api.SearchConversions(searchReqStr)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", trs, w)
	return
}

// GetConversionDetailsHandler uses FastAPI for get conversion details by it id.
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
	trs, err := api.GetConversion(id)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	writeResult(http.StatusOK, "OK", trs, w)
	return
}
