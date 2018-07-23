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
package tests

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/eeonevision/anychaindb/api/handler"
	"github.com/eeonevision/anychaindb/client"
)

var count = flag.Int64("count", 10000, "test data count")

func TestGenerateAccountsForYandexTank(t *testing.T) {
	// Create folder if not exists
	_ = os.Mkdir("artifacts/tank", os.ModePerm)
	// Generate transactions batches
	var requests []string
	for i := int64(0); i < *count; i++ {
		body := []string{}
		data, _ := json.Marshal(body)
		query := fmt.Sprintf("/v1/accounts")
		tr := fmt.Sprintf("POST %s HTTP/1.1\nHost: %s\nContent-Type: application/json\nUser-Agent: Tank\nContent-Length: %v\n\n%s",
			query,
			*host+":"+*apiPort,
			len(data), data)
		tr = fmt.Sprintf("%v\n%s\n", len(tr), tr)
		requests = append(requests, tr)
	}

	err := writeLines(requests, "artifacts/tank/accounts_tank.txt")
	if err != nil {
		t.Errorf("error write to file")
	}
}

func TestGeneratePayloadsForYandexTank(t *testing.T) {
	endpointURL := fmt.Sprintf("http://%s:%s", *host, *rpcPort)
	// Create folder if not exists
	_ = os.Mkdir("artifacts/tank", os.ModePerm)
	// Create folder if not exists
	_ = os.Mkdir("artifacts/tank/tank_"+*host, os.ModePerm)
	// Create new affiliate account
	api := client.NewAPI(endpointURL, *apiMode, nil, "")
	id, pub, priv, err := api.CreateAccount()
	if err != nil {
		t.Errorf("error to create account: " + err.Error())
	}
	t.Logf("created account: %s, %s, %s", id, pub, priv)
	//Generate conversions requests
	var requests []string
	for i := int64(0); i < *count; i++ {
		body := handler.Payload{
			PublicData: "test_public",
			PrivateData: []*handler.PrivateData{&handler.PrivateData{
				ReceiverAccountID: id,
				Data:              "test_private",
			},
			},
		}
		req := handler.Request{
			AccountID: id,
			PrivKey:   priv,
			PubKey:    pub,
			Data:      body,
		}
		data, _ := json.Marshal(req)
		tr := fmt.Sprintf("POST %s HTTP/1.1\nHost: %s\nContent-Type: application/json\nUser-Agent: Tank\nContent-Length: %v\n\n%s",
			"/v1/payloads",
			*host+":"+*apiPort,
			len(data), data)
		tr = fmt.Sprintf("%v\n%s\n", len(tr), tr)
		requests = append(requests, tr)
	}
	// Write to file
	err = writeLines(requests, "artifacts/tank/tank_"+*host+"/"+"payloads_tank_"+*host+".txt")
	if err != nil {
		t.Errorf("error write to file")
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
