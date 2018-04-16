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

package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/leadschain/leadschain/api/handler"
)

type server struct {
	ListenHost string
	ListenPort string
	Logger     *log.Logger
}

// NewHTTPServer method constructs new server object and handlers
func NewHTTPServer(GRPCEndpoint, listenIP, httpPort string) *server {
	handler.SetEndpoint(GRPCEndpoint)

	res := server{
		ListenHost: listenIP,
		ListenPort: httpPort,
		Logger:     log.New(os.Stderr, "", log.LstdFlags),
	}

	m := httprouter.New()
	// Accounts
	m.POST("/v1/accounts", handler.PostAccountsHandler)
	// Transitions
	m.GET("/v1/transitions", handler.GetTransitionsHandler)
	m.GET("/v1/transitions/:id", handler.GetTransitionDetailsHandler)
	m.POST("/v1/transitions", handler.PostTransitionsHandler)
	//Conversions
	m.GET("/v1/conversions", handler.GetConversionsHandler)
	m.GET("/v1/conversions/:id", handler.GetConversionDetailsHandler)
	m.POST("/v1/conversions", handler.PostConversionsHandler)

	http.Handle("/", m)

	return &res
}

func (s *server) SetLogger(l *log.Logger) {
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
		s.Logger.Printf("Listening on http://%s:%s\n", s.ListenHost, s.ListenPort)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			s.Logger.Fatal(err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
