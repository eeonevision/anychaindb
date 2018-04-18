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
	"fmt"
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leadschain/leadschain/state"
	"github.com/tendermint/abci/types"
)

type PersistentApplication struct {
	app         *Application
	blockHeader types.Header
	changes     []types.Validator
	logger      *log.Logger
}

func NewPersistentApplication(dbHost, dbName string) *PersistentApplication {
	stateDB, err := mgo.Dial(dbHost)
	if err != nil {
		panic("Error initialize Mongo DB: " + err.Error())
	}
	return &PersistentApplication{
		app:    &Application{state: state.NewStateFromDB(stateDB.DB(dbName))},
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (app *PersistentApplication) SetLogger(l *log.Logger) {
	app.logger = l
}

func (app *PersistentApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	resInfo = app.app.Info()
	lastBlock := app.LoadLastBlock()
	resInfo.LastBlockHeight = lastBlock.Height
	resInfo.LastBlockAppHash = lastBlock.AppHash
	return resInfo
}

func (app *PersistentApplication) SetOption(reqSetOpt types.RequestSetOption) types.ResponseSetOption {
	return app.app.SetOption(reqSetOpt)
}

func (app *PersistentApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	return app.app.DeliverTx(tx)
}

func (app *PersistentApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	return app.app.CheckTx(tx)
}

func (app *PersistentApplication) Commit() types.ResponseCommit {
	appCommit := app.app.Commit()

	lastBlock := LastBlockInfo{
		Height:  app.blockHeader.Height,
		AppHash: appCommit.Data, // this hash will be in the next block header
	}

	app.SaveLastBlock(lastBlock)
	return appCommit
}

func (app *PersistentApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	return app.app.Query(reqQuery)
}

func (app *PersistentApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	return types.ResponseInitChain{}
}

// BeginBlock method tracks the block hash and header information
func (app *PersistentApplication) BeginBlock(req types.RequestBeginBlock) (resBeginBlock types.ResponseBeginBlock) {
	// update latest block info
	app.blockHeader = req.Header

	// reset valset changes
	app.changes = make([]types.Validator, 0)
	return types.ResponseBeginBlock{}
}

// EndBlock method should in future update the validator set
func (app *PersistentApplication) EndBlock(reqEndBlock types.RequestEndBlock) (resEndBlock types.ResponseEndBlock) {
	return types.ResponseEndBlock{ValidatorUpdates: app.changes}
}

type LastBlockInfo struct {
	Height  int64  `bson:"height"`
	AppHash []byte `bson:"app_hash"`
}

// LoadLastBlock method load last confirmed block from DB
func (app *PersistentApplication) LoadLastBlock() (lastBlock LastBlockInfo) {
	if err := app.app.state.DB.C("blocks").Find(nil).One(&lastBlock); err != mgo.ErrNotFound && err != nil {
		panic(err)
	}
	fmt.Println("Last block is: ", string(lastBlock.AppHash))
	return lastBlock
}

// SaveLastBlock method saves appHash of loast confirmed block in DB
func (app *PersistentApplication) SaveLastBlock(lastBlock LastBlockInfo) {
	lb := app.LoadLastBlock()
	selector := bson.M{"height": lb.Height}
	updator := bson.M{"$set": bson.M{"height": lastBlock.Height, "app_hash": lastBlock.AppHash}}
	_, err := app.app.state.DB.C("blocks").Upsert(selector, updator)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Saved block %v with height %s", lastBlock.Height, string(lastBlock.AppHash))
}
