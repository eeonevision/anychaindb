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

package main

import (
	"flag"
	"os"

	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"

	labci "github.com/eeonevision/anychaindb/abci-app"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/common"
)

func main() {
	// Parse CLI arguments
	addrPtr := flag.String("addr", "tcp://localhost:26658", "Listen address")
	abciPtr := flag.String("abci", "socket", "socket | grpc")
	dbHost := flag.String("dbhost", "localhost", "database host path")
	dbName := flag.String("dbname", "anychaindb", "database name")
	logLevel := flag.String("loglevel", "*:info", "log level for abci modules: abci-app:info,abci-server:info,*:error")
	flag.Parse()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err := tmflags.ParseLogLevel(*logLevel, logger, "info")
	if err != nil {
		panic(err)
	}

	app := labci.NewPersistentApplication(*dbHost, *dbName)
	app.SetLogger(logger.With("module", "abci-app"))

	// Start the listener
	srv, err := server.NewServer(*addrPtr, *abciPtr, app)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Wait forever
	common.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}
