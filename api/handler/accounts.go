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
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/leadschain/leadschain/client"
)

// Account struct describes account related fields
//
// ID - unique identifier of account in blockchain
// Priv - private key of account
// Pub - public key of account
type Account struct {
	ID   string `json:"id" mapstructure:"id"`
	Priv string `json:"private_key" mapstructure:"private_key"`
	Pub  string `json:"public_key" mapstructure:"public_key"`
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

	// OK
	writeResult(http.StatusAccepted, "Accepted",
		Account{
			ID:   id,
			Priv: priv,
			Pub:  pub,
		}, w)
}
