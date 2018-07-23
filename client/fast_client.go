package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tendermint/tendermint/rpc/core/types"

	"github.com/eeonevision/anychaindb/crypto"
	"github.com/eeonevision/anychaindb/state"
	"github.com/eeonevision/anychaindb/transaction"
)

type ResultWrapper struct {
	Result *core_types.ResultABCIQuery `json:"result"`
}

// fastClient struct contains config
// parameters for performing requests.
type fastClient struct {
	key       *crypto.Key
	endpoint  string
	mode      string
	accountID string
	client    *http.Client
}

// newFastClient initializes new fast client instance.
func newFastClient(endpoint, mode string, key *crypto.Key, accountID string) *fastClient {
	// Set default mode
	switch mode {
	case "sync":
		mode = "sync"
		break
	case "async":
		mode = "async"
		break
	case "commit":
		mode = "commit"
		break
	default:
		mode = "sync"
	}

	return &fastClient{key, endpoint, mode, accountID, &http.Client{Timeout: 30 * time.Second}}
}

func (c *fastClient) abciQuery(path, data string) (*core_types.ResultABCIQuery, error) {
	var res *ResultWrapper

	req, err := http.NewRequest("GET", c.endpoint+"/abci_query?path=\""+path+"\"&data=\""+data+"\"", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("empty ABCI response")
	}
	if res.Result.Response.IsErr() {
		return nil, errors.New(res.Result.Response.GetLog())
	}

	return res.Result, nil
}

func (c *fastClient) broadcastTx(tx []byte) (interface{}, error) {
	var res interface{}

	ba := base64.StdEncoding.EncodeToString(tx)
	txData := fmt.Sprintf(`{"jsonrpc":"2.0","id":"anything","method":"broadcast_tx_%s","params": {"tx": "%s"}}`, c.mode, ba)
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer([]byte(txData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &res)
	// Check transport errors
	if err != nil {
		return nil, err
	}
	// Check special empty case
	if res == nil {
		return nil, errors.New("empty ABCI response")
	}
	// Check for async/sync response
	if r, ok := res.(*core_types.ResultBroadcastTx); ok && r.Code != 0 {
		return nil, errors.New(r.Log)
	}
	// Check for commit response
	if r, ok := res.(*core_types.ResultBroadcastTxCommit); ok && (r.CheckTx.Code != 0 || r.DeliverTx.Code != 0) {
		return nil, errors.New("check tx error: " + r.CheckTx.Log + "; deliver tx error: " + r.DeliverTx.Log)
	}

	return res, nil
}

func (c *fastClient) getAccount(id string) (*state.Account, error) {
	resp, err := c.abciQuery("accounts", id)
	if err != nil {
		return nil, err
	}
	acc := &state.Account{}
	if err := json.Unmarshal(resp.Response.GetValue(), &acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (c *fastClient) addPayload(cv *state.Payload) error {
	txBytes, err := cv.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.PayloadAdd, c.accountID, txBytes)
	if err := tx.Sign(c.key); err != nil {
		return err
	}
	bs, _ := tx.ToBytes()

	_, err = c.broadcastTx(bs)
	if err != nil {
		return err
	}

	return nil
}
