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
package tests

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"gitlab.com/leadschain/leadschain/client"

	"gitlab.com/leadschain/leadschain/api/handler"
)

var host = flag.String("host", "localhost", "machine host")
var apiPort = flag.String("apiport", "8889", "api port")
var rpcPort = flag.String("rpcport", "46657", "rpc port")
var count = flag.Int64("count", 0, "requests count")

func TestGenerateYandexTankAccounts(t *testing.T) {
	apiHost := *host
	apiPort := *apiPort
	accountsCount := *count
	// Create folder if not exists
	_ = os.Mkdir("artifacts", os.ModePerm)
	// Generate transactions batches
	var requests []string
	for i := int64(0); i < accountsCount; i++ {
		body := []string{}
		data, _ := json.Marshal(body)
		query := fmt.Sprintf("/v1/accounts")
		tr := fmt.Sprintf("POST %s HTTP/1.1\nHost: %s\nContent-Type: application/json\nUser-Agent: Tank\nContent-Length: %v\n\n%s",
			query,
			apiHost+":"+apiPort,
			len(data), data)
		tr = fmt.Sprintf("%v\n%s\n", len(tr), tr)
		requests = append(requests, tr)
	}

	err := writeLines(requests, "artifacts/accounts_tank.txt")
	if err != nil {
		t.Errorf("Error write to file")
	}
}

func TestGenerateYandexTankTransitions(t *testing.T) {
	endpointURL := fmt.Sprintf("http://%s:%s", *host, *rpcPort)
	apiHost := *host
	apiPort := *apiPort
	transitionsCount := *count
	// Create folder if not exists
	_ = os.Mkdir("artifacts", os.ModePerm)
	// Create new affiliate account
	api := client.NewAPI(endpointURL, nil, "")
	id, pub, priv, err := api.CreateAccount()
	if err != nil {
		t.Errorf("Error to create account: " + err.Error())
	}
	t.Logf("Created account: %s, %s, %s", id, pub, priv)
	//Generate transitions requests
	var requests []string
	for i := int64(0); i < transitionsCount; i++ {
		body := handler.Transition{
			AffiliateAccountID:  id,
			AdvertiserAccountID: id,
			ClickID:             "test click",
			OfferID:             "test offer",
			StreamID:            "test stream",
			ExpiresIn:           15343242343,
		}
		req := handler.Request{
			AccountID: id,
			PrivKey:   priv,
			PubKey:    pub,
			Data:      body,
		}
		data, _ := json.Marshal(req)
		query := fmt.Sprintf("/v1/transitions")
		tr := fmt.Sprintf("POST %s HTTP/1.1\nHost: %s\nContent-Type: application/json\nUser-Agent: Tank\nContent-Length: %v\n\n%s",
			query,
			apiHost+":"+apiPort,
			len(data), data)
		tr = fmt.Sprintf("%v\n%s\n", len(tr), tr)
		requests = append(requests, tr)
	}
	// Write to file
	err = writeLines(requests, "artifacts/transitions_tank.txt")
	if err != nil {
		t.Errorf("Error write to file")
	}
}

func TestGenerateYandexTankConversions(t *testing.T) {
	endpointURL := fmt.Sprintf("http://%s:%s", *host, *rpcPort)
	apiHost := *host
	apiPort := *apiPort
	conversionsCount := *count
	// Create folder if not exists
	_ = os.Mkdir("artifacts", os.ModePerm)
	// Create folder if not exists
	_ = os.Mkdir("artifacts/tank_"+apiHost, os.ModePerm)
	// Create new affiliate account
	api := client.NewAPI(endpointURL, nil, "")
	id, pub, priv, err := api.CreateAccount()
	if err != nil {
		t.Errorf("Error to create account: " + err.Error())
	}
	t.Logf("Created account: %s, %s, %s", id, pub, priv)
	//Generate conversions requests
	var requests []string
	for i := int64(0); i < conversionsCount; i++ {
		body := handler.Conversion{
			AffiliateAccountID:  id,
			AdvertiserAccountID: id,
			ClickID:             "test click",
			OfferID:             "test offer",
			ClientID:            "test client",
			GoalID:              "test goal",
			StreamID:            "test stream",
			Comment:             "Test comment",
			Status:              "CONFIRMED",
		}
		req := handler.Request{
			AccountID: id,
			PrivKey:   priv,
			PubKey:    pub,
			Data:      body,
		}
		data, _ := json.Marshal(req)
		tr := fmt.Sprintf("POST %s HTTP/1.1\nHost: %s\nContent-Type: application/json\nUser-Agent: Tank\nContent-Length: %v\n\n%s",
			"/v1/conversions",
			apiHost+":"+apiPort,
			len(data), data)
		tr = fmt.Sprintf("%v\n%s\n", len(tr), tr)
		requests = append(requests, tr)
	}
	// Write to file
	err = writeLines(requests, "artifacts/tank_"+apiHost+"/"+"conversions_tank_"+apiHost+".txt")
	if err != nil {
		t.Errorf("Error write to file")
	}
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprint(w, line)
	}
	return w.Flush()
}
