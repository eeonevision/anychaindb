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

package app

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/eeonevision/anychaindb/state"
	"github.com/eeonevision/anychaindb/transaction"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tendermint/tendermint/abci/types"
)

// Limit response elements count
const (
	resLimit  = 500
	resOffset = 0
)

// Application inherits BaseApplication and keeps state of anychaindb
type Application struct {
	types.BaseApplication
	state  *state.State
	logger log.Logger
}

// NewApplication method initializes new application with MongoDB state
func NewApplication(dbHost, dbName string) *Application {
	db, err := mgo.Dial(dbHost)
	if err != nil {
		panic("Error initialize Mongo DB: " + err.Error())
	}
	return &Application{state: state.NewStateFromDB(db.DB(dbName))}
}

// SetLogger method set logger for Application
func (app *Application) SetLogger(l log.Logger) {
	app.logger = l
}

// Info method returns information about current state.
// All sizes represented in kilobytes.
func (app *Application) Info() (resInfo types.ResponseInfo) {
	var stats map[string]interface{}
	if err := app.state.DB.Run(bson.M{"dbStats": 1}, &stats); err != nil {
		app.logger.Error("Getting state info error", "error", err.Error())
		return
	}
	res, err := json.Marshal(stats)
	if err != nil {
		app.logger.Error("Encoding state info error", "error", err.Error())
		return
	}
	return types.ResponseInfo{Data: string(res)}
}

// DeliverTx method responsible for deliver chosen transaction type.
func (app *Application) DeliverTx(txBytes []byte) types.ResponseDeliverTx {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}
	switch tx.Type {
	case transaction.AccountAdd:
		{
			if err := deliverAccountAddTransaction(tx, app.state); err != nil {
				return types.ResponseDeliverTx{
					Code: CodeTypeDeliverTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.PayloadAdd:
		{
			if err := deliverPayloadAddTransaction(tx, app.state); err != nil {
				return types.ResponseDeliverTx{
					Code: CodeTypeDeliverTxError,
					Log:  err.Error(),
				}
			}
		}
	default:
		{
			return types.ResponseDeliverTx{
				Code: CodeTypeUnknownRequest,
				Log:  "transaction type was not detected.",
			}
		}
	}
	return types.ResponseDeliverTx{
		Code: CodeTypeOK,
	}
}

// CheckTx method responsible for check transaction on validity.
func (app *Application) CheckTx(txBytes []byte) types.ResponseCheckTx {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ResponseCheckTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}
	switch tx.Type {
	case transaction.AccountAdd:
		{
			if err := checkAccountAddTransaction(tx, app.state); err != nil {
				return types.ResponseCheckTx{
					Code: CodeTypeCheckTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.PayloadAdd:
		{
			if err := checkPayloadAddTransaction(tx, app.state); err != nil {
				return types.ResponseCheckTx{
					Code: CodeTypeCheckTxError,
					Log:  err.Error(),
				}
			}
		}
	default:
		{
			return types.ResponseCheckTx{
				Code: CodeTypeUnknownRequest,
				Log:  "transaction type was not detected.",
			}
		}
	}
	return types.ResponseCheckTx{
		Code: CodeTypeOK,
	}
}

// Commit method generates hash of current state.
func (app *Application) Commit() types.ResponseCommit {
	var hash map[string]interface{}
	for {
		if err := app.state.DB.Run(bson.M{
			"dbhash":      1,
			"collections": []string{"accounts", "data"},
		}, &hash); err == nil {
			return types.ResponseCommit{Data: []byte(hash["md5"].(string))}
		}
	}
}

// mongoQuery is a struct for parse search query from a user.
type mongoQuery struct {
	Query  interface{} `json:"query,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
}

// Query method processes user's request.
// Search endpoint uses Mongo DB query syntax.
// For make search request uses mongoQuery struct.
// The limit and offset fields are optional. All mongo query places in query field of struct.
// Check mongo query syntax at:
// https://docs.mongodb.com/manual/tutorial/query-documents/
func (app *Application) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	var (
		result interface{}
		err    error
	)
	switch reqQuery.Path {
	case "accounts":
		{
			if reqQuery.Data == nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = "id is not presented in query"
				return
			}
			result, err = app.state.GetAccount(string(reqQuery.Data))
			if err != nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "accounts/search":
		{
			if reqQuery.Data == nil {
				resQuery.Code = CodeEmptySearchQuery
				resQuery.Log = "search query is empty"
				return
			}
			// Unmarshal search query
			var mgoQuery mongoQuery
			if err = json.Unmarshal(reqQuery.Data, &mgoQuery); err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			// Check limit and offset values
			if mgoQuery.Limit > resLimit || resOffset <= 0 {
				mgoQuery.Limit = resLimit
			}
			if mgoQuery.Offset < resOffset {
				mgoQuery.Offset = resOffset
			}
			// Search accounts in Database
			result, err = app.state.SearchAccounts(mgoQuery.Query, mgoQuery.Limit, mgoQuery.Offset)
			if err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "payloads":
		{
			if reqQuery.Data == nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = "id is not presented in query"
				return
			}
			result, err = app.state.GetPayload(string(reqQuery.Data))
			if err != nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "payloads/search":
		{
			if reqQuery.Data == nil {
				resQuery.Code = CodeEmptySearchQuery
				resQuery.Log = "search query is empty"
				return
			}
			// Unmarshal search query
			var mgoQuery mongoQuery
			if err = json.Unmarshal(reqQuery.Data, &mgoQuery); err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			// Check limit and offset values
			if mgoQuery.Limit > resLimit || resOffset <= 0 {
				mgoQuery.Limit = resLimit
			}
			if mgoQuery.Offset < resOffset {
				mgoQuery.Offset = resOffset
			}
			// Search transaction data in Database
			result, err = app.state.SearchPayloads(mgoQuery.Query, mgoQuery.Limit, mgoQuery.Offset)
			if err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	default:
		{
			resQuery.Code = CodeTypeUnknownRequest
			resQuery.Log = "the path was not detected"
			return
		}
	}
	resQuery.Code = CodeTypeOK
	return
}
