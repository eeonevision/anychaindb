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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/leadschain/leadschain/api/handler"
)

var host = flag.String("host", "localhost", "machine host")
var apiPort = flag.String("apiport", "8889", "api port")
var rpcPort = flag.String("rpcport", "46657", "rpc port")
var update = flag.Bool("update", false, "update .golden files")

func doPOSTRequest(endpoint, url string, data []byte) ([]byte, error) {
	// Send request to Leadschain
	respRaw, err := http.Post("http://"+url+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer respRaw.Body.Close()
	contents, err := ioutil.ReadAll(respRaw.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func doGETRequest(endpoint, url string) ([]byte, error) {
	respRaw, err := http.Get("http://" + url + endpoint)
	if err != nil {
		return nil, err
	}
	defer respRaw.Body.Close()
	contents, err := ioutil.ReadAll(respRaw.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func TestCreateAccount(t *testing.T) {
	// Generate transaction request
	endpoint := fmt.Sprintf("/v1/accounts")
	url := *host + ":" + *apiPort
	data, _ := json.Marshal(handler.Request{})
	contents, err := doPOSTRequest(endpoint, url, data)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Check data in results
	resp := handler.Result{}
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	account := resp.Data.(map[string]interface{})
	if account["ID"] == "" {
		t.Errorf("account has no ID: %s", account)
		return
	}
	// Write account data to account.golden file
	if *update {
		// Create folder if not exists
		_ = os.Mkdir("artifacts", os.ModePerm)
		// Marshal account data to JSON string
		res, err := json.Marshal(resp.Data)
		if err != nil {
			t.Errorf("%s", err)
			return
		}
		// Write account id
		err = ioutil.WriteFile("artifacts/account.golden", res, 0644)
		if err != nil {
			t.Errorf("%s", err)
			return
		}
		// Wait for transaction approve
		time.Sleep(time.Second * 5)
	}
}

func TestGetAccount(t *testing.T) {
	// Generate transaction request
	endpoint := fmt.Sprintf("/v1/accounts")
	url := *host + ":" + *apiPort
	// Read account.golden file
	file, err := ioutil.ReadFile("artifacts/account.golden")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Parse account.golden file
	acc1 := handler.Account{}
	err = json.Unmarshal(file, &acc1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if len(acc1.ID) == 0 {
		t.Errorf("account.golden is empty: %s", acc1)
		return
	}
	endpoint = endpoint + "/" + acc1.ID
	// Find account in Leadschain server
	contents, err := doGETRequest(endpoint, url)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	resp := handler.Result{}
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Compare accounts
	acc2 := resp.Data.(map[string]interface{})
	if acc2["_id"] != acc1.ID {
		t.Errorf("accounts are not equal. Expected: %v, Output: %v", acc1.ID, acc2["_id"])
		return
	}
}

func TestCreateConversion(t *testing.T) {
	// Generate transaction request
	endpoint := fmt.Sprintf("/v1/conversions")
	url := *host + ":" + *apiPort
	// Read account.golden file
	file, err := ioutil.ReadFile("artifacts/account.golden")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Parse account.golden file
	acc1 := handler.Account{}
	err = json.Unmarshal(file, &acc1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if len(acc1.ID) == 0 {
		t.Errorf("account.golden is empty: %s", acc1)
		return
	}
	// Send request
	data, _ := json.Marshal(handler.Request{
		AccountID: acc1.ID,
		PrivKey:   acc1.Priv,
		PubKey:    acc1.Pub,
		Data: handler.Conversion{
			AffiliateAccountID: acc1.ID,
			PrivateData:        "test_data",
			PublicData:         "test_public_data",
		}})
	contents, err := doPOSTRequest(endpoint, url, data)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Check data in results
	resp := handler.Result{}
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	cnv := resp.Data.(map[string]interface{})
	if cnv["_id"] == "" {
		t.Errorf("conversion has no id: %s", cnv)
		return
	}
	// Write conversion data to conversion.golden file
	if *update {
		// Create folder if not exists
		_ = os.Mkdir("artifacts", os.ModePerm)
		// Marshal conversion data to JSON string
		res, err := json.Marshal(resp.Data)
		if err != nil {
			t.Errorf("%s", err)
			return
		}
		// Write account id
		err = ioutil.WriteFile("artifacts/conversion.golden", res, 0644)
		if err != nil {
			t.Errorf("%s", err)
			return
		}
		// Wait for transaction approve
		time.Sleep(time.Second * 5)
	}
}

func TestGetConversion(t *testing.T) {
	// Generate transaction request
	endpoint := fmt.Sprintf("/v1/conversions")
	url := *host + ":" + *apiPort
	// Read conversion.golden file
	file, err := ioutil.ReadFile("artifacts/conversion.golden")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Parse conversion.golden file
	cnv1 := handler.Conversion{}
	err = json.Unmarshal(file, &cnv1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if len(cnv1.ID) == 0 {
		t.Errorf("conversion.golden is empty: %v", cnv1)
		return
	}
	endpoint = endpoint + "/" + cnv1.ID
	// Find conversion in Leadschain server
	contents, err := doGETRequest(endpoint, url)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	resp := handler.Result{}
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Compare conversions
	cnv2 := resp.Data.(map[string]interface{})
	if cnv2["_id"] != cnv1.ID {
		t.Errorf("conversions are not equal. Expected: %v, Output: %v", cnv1.ID, cnv2["_id"])
		return
	}
}
