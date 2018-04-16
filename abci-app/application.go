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

package app

import (
	"encoding/json"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tendermint/abci/types"
	"gitlab.com/leadschain/leadschain/state"
	"gitlab.com/leadschain/leadschain/transaction"
)

type Application struct {
	types.BaseApplication
	state *state.State
}

func NewApplication(dbHost, dbName string) *Application {
	db, err := mgo.Dial(dbHost)
	if err != nil {
		panic("Error initialize Mongo DB: " + err.Error())
	}
	return &Application{state: state.NewStateFromDB(db.DB(dbName))}
}

// Info method returns infomation about current state.
// All sizes represented in kilobytes.
func (app *Application) Info() (resInfo types.ResponseInfo) {
	var stats map[string]interface{}
	if err := app.state.DB.Run(bson.M{"dbStats": 1, "scale": 1024}, &stats); err != nil {
		return types.ResponseInfo{Data: err.Error()}
	}
	res, err := json.Marshal(stats)
	if err != nil {
		return types.ResponseInfo{Data: err.Error()}
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

	case transaction.AccountDel:
		{
			if err := deliverAccountDelTransaction(tx, app.state); err != nil {
				return types.ResponseDeliverTx{
					Code: CodeTypeDeliverTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.TransitionAdd:
		{
			if err := deliverTransitionAddTransaction(tx, app.state); err != nil {
				return types.ResponseDeliverTx{
					Code: CodeTypeDeliverTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.ConversionAdd:
		{
			if err := deliverConversionAddTransaction(tx, app.state); err != nil {
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
				Log:  "Transaction type was not detected.",
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
	case transaction.AccountDel:
		{
			if err := checkAccountDelTransaction(tx, app.state); err != nil {
				return types.ResponseCheckTx{
					Code: CodeTypeCheckTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.TransitionAdd:
		{
			if err := checkTransitionAddTransaction(tx, app.state); err != nil {
				return types.ResponseCheckTx{
					Code: CodeTypeCheckTxError,
					Log:  err.Error(),
				}
			}
		}
	case transaction.ConversionAdd:
		{
			if err := checkConversionAddTransaction(tx, app.state); err != nil {
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
				Log:  "Transaction type was not detected.",
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
			"collections": []string{"accounts", "transitions", "conversions"},
		}, &hash); err == nil {
			return types.ResponseCommit{Data: []byte(hash["md5"].(string))}
		}
	}
}

// MongoQuery is a struct for parse search query from a user.
type MongoQuery struct {
	Query  interface{} `json:"query,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
}

// Query method processes user's request.
// Search endpoint uses Mongo DB query syntax.
// For make search request uses MongoQuery struct.
// The limit and offset fields are optional. All mongo query places in query field of struct.
// Check mongo query syntax at:
// https://docs.mongodb.com/manual/tutorial/query-documents/
func (app *Application) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	switch reqQuery.Path {
	case "accounts":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.state.ListAccounts()
				log.Printf("Got account list: %+s", result)
			} else {
				id := string(reqQuery.Data)
				result, err = app.state.GetAccount(id)
				log.Printf("Got account: %+v", result)
			}
			if err != nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "transitions":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.state.ListTransitions()
				log.Printf("Got transitions list: %+s", result)
			} else {
				result, err = app.state.GetTransition(string(reqQuery.Data))
				log.Printf("Got transition: %+v", result)
			}
			if err != nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "transitions/search":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				resQuery.Code = CodeEmptySearchQuery
				resQuery.Log = "Search query is empty"
				return
			}
			// Unmarshal search query
			var mgoQuery MongoQuery
			if err = json.Unmarshal(reqQuery.Data, &mgoQuery); err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			// Search transitions in Database
			result, err = app.state.SearchTransitions(mgoQuery.Query, mgoQuery.Limit, mgoQuery.Offset)
			if err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			log.Printf("Got transitions: %+s", result)
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "conversions":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.state.ListConversions()
				log.Printf("Got conversions list: %+s", result)
			} else {
				result, err = app.state.GetConversion(string(reqQuery.Data))
				log.Printf("Got conversion: %+v", result)
			}
			if err != nil {
				resQuery.Code = CodeTypeQueryError
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "conversions/search":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				resQuery.Code = CodeEmptySearchQuery
				resQuery.Log = "Search query is empty"
				return
			}
			// Unmarshal search query
			var mgoQuery MongoQuery
			if err = json.Unmarshal(reqQuery.Data, &mgoQuery); err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			// Search conversions in Database
			result, err = app.state.SearchConversions(mgoQuery.Query, mgoQuery.Limit, mgoQuery.Offset)
			if err != nil {
				resQuery.Code = CodeParseSearchQueryError
				resQuery.Log = err.Error()
				return
			}
			log.Printf("Got conversions: %+s", result)
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	default:
		{
			resQuery.Code = CodeTypeUnknownRequest
			resQuery.Log = "The path was not detected"
			return
		}
	}
	resQuery.Code = CodeTypeOK
	return
}
