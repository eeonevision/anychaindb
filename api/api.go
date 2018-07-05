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

package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tendermint/tmlibs/log"

	"github.com/julienschmidt/httprouter"

	"github.com/leadschain/leadschain/api/handler"
)

type server struct {
	ListenHost string
	ListenPort string
	Logger     log.Logger
}

// NewHTTPServer method constructs new server object and handlers
func NewHTTPServer(GRPCEndpoint, listenIP, httpPort string) *server {
	handler.SetEndpoint(GRPCEndpoint)

	res := server{
		ListenHost: listenIP,
		ListenPort: httpPort,
		Logger:     log.NewNopLogger(),
	}

	m := httprouter.New()
	// Accounts
	m.GET("/v1/accounts", handler.GetAccountsHandler)
	m.GET("/v1/accounts/:id", handler.GetAccountDetailsHandler)
	m.POST("/v1/accounts", handler.PostAccountsHandler)
	// Payloads
	m.GET("/v1/payloads", handler.GetPayloadsHandler)
	m.GET("/v1/payloads/:id", handler.GetPayloadDetailsHandler)
	m.POST("/v1/payloads", handler.PostPayloadsHandler)

	http.Handle("/", m)

	return &res
}

func (s *server) SetLogger(l log.Logger) {
	s.Logger = l
}

func (s *server) Serve() {
	listenString := s.ListenHost + ":" + s.ListenPort

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	srv := &http.Server{
		Addr:         listenString,
		Handler:      nil,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		s.Logger.Info("Starting REST-API server...", "host", s.ListenHost, "port", s.ListenPort)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			s.Logger.Error("Starting error", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-stopChan // wait for SIGINT
	s.Logger.Info("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()

	s.Logger.Info("Server gracefully stopped")
}
