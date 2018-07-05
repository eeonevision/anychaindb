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

	lapi "github.com/eeonevision/anychaindb/api"
	tmflags "github.com/tendermint/tmlibs/cli/flags"
	"github.com/tendermint/tmlibs/log"
)

func main() {
	// Parse CLI arguments
	endpointPtr := flag.String("endpoint", "http://0.0.0.0:46657", "Validator grpc endpoint address")
	ipPtr := flag.String("ip", "localhost", "Listen host ip")
	portPtr := flag.String("port", "8888", "Listen host port")
	logLevel := flag.String("loglevel", "*:info", "log level for anychaindb api module: rest-api:info")
	flag.Parse()

	// Create server
	api := lapi.NewHTTPServer(*endpointPtr, *ipPtr, *portPtr)

	// Define logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err := tmflags.ParseLogLevel(*logLevel, logger, "info")
	if err != nil {
		panic(err)
	}
	api.SetLogger(logger.With("module", "rest-api"))

	// Start listener
	api.Serve()
}
