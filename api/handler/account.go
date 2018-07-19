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
	"github.com/julienschmidt/httprouter"
)

// Account struct describes account related fields
//
// ID - unique identifier of account in blockchain
// Priv - private key of account
// Pub - public key of account
type Account struct {
	ID   string `json:"_id"`
	Priv string `json:"private_key"`
	Pub  string `json:"public_key"`
}

// PostAccountsHandler uses FastAPI for sends new accounts requests in async mode to blockchain
func PostAccountsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	defer r.Body.Close()

	// Add account to blockchain
	api := client.NewAPI(endpoint, nil, "")
	id, pub, priv, err := api.CreateAccount()
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	writeResult(http.StatusAccepted, "Accepted",
		Account{
			ID:   id,
			Priv: priv,
			Pub:  pub,
		}, w)
	return
}

// GetAccountsHandler uses BaseAPI for search and list accounts.
// Query parameters: Query, Limit, Offset can be optional.
// Query - MongoDB query string.
// Limit - maximum 500 items.
// Offset - default 0.
func GetAccountsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	acc, err := api.SearchAccounts(searchReqStr)
	if err != nil {
		writeResult(http.StatusBadRequest, err.Error(), nil, w)
		return
	}

	writeResult(http.StatusOK, "OK", acc, w)
	return
}

// GetAccountDetailsHandler uses BaseAPI for get conversion details by it id.
// Query parameters ID is required.
func GetAccountDetailsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("id")
	if id == "" {
		writeResult(http.StatusBadRequest,
			"ID should not be empty", nil, w)
		return
	}
	api := client.NewAPI(endpoint, nil, "")
	acc, err := api.GetAccount(id)

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

	writeResult(http.StatusOK, "OK", acc, w)
	return
}
