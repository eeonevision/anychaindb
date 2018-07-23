package client

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tendermint/tendermint/rpc/lib/types"

	"github.com/tendermint/tendermint/rpc/core/types"

	"github.com/eeonevision/anychaindb/crypto"
	"github.com/eeonevision/anychaindb/state"
	"github.com/eeonevision/anychaindb/transaction"
)

// Error variables defines empty response from ABCI and RPC.
var (
	ErrEmptyABCIResponse = errors.New("empty ABCI response")
	ErrEmptyRPCResponse  = errors.New("empty RPC response")
)

// Broadcast Modes.
var (
	sync   = "sync"
	async  = "async"
	commit = "commit"
)

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
		mode = sync
		break
	case "async":
		mode = async
		break
	case "commit":
		mode = commit
		break
	default:
		mode = sync
	}
	return &fastClient{key, endpoint, mode, accountID, &http.Client{Timeout: 30 * time.Second}}
}

func (c *fastClient) doPOSTRequest(method, data string) (*rpctypes.RPCResponse, error) {
	var rpcRes *rpctypes.RPCResponse
	txData := fmt.Sprintf(`{"jsonrpc":"2.0","id":"anything","method":"%s","params": %s}`, method, data)
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

	// Check RPC response
	err = json.Unmarshal(contents, &rpcRes)
	if err != nil {
		return nil, err
	}
	if rpcRes == nil {
		return nil, ErrEmptyRPCResponse
	}
	if rpcRes.Error != nil {
		return nil, rpcRes.Error
	}
	// Check ABCI result
	if rpcRes.Result == nil {
		return nil, ErrEmptyABCIResponse
	}

	return rpcRes, nil
}

func (c *fastClient) abciQuery(path string, data []byte) (*core_types.ResultABCIQuery, error) {
	var rpcRes *rpctypes.RPCResponse
	var abciRes *core_types.ResultABCIQuery

	rpcRes, err := c.doPOSTRequest("abci_query", fmt.Sprintf(`{"path": "%s", "data": "%s"}`, path, hex.EncodeToString(data)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rpcRes.Result, &abciRes)
	if err != nil {
		return nil, err
	}
	if abciRes.Response.IsErr() {
		return nil, errors.New(abciRes.Response.GetLog())
	}

	return abciRes, nil
}

func (c *fastClient) broadcastTx(tx []byte) (interface{}, error) {
	var rpcRes *rpctypes.RPCResponse

	rpcRes, err := c.doPOSTRequest("broadcast_tx_"+c.mode, fmt.Sprintf(`{"tx": "%s"}`, base64.StdEncoding.EncodeToString(tx)))
	if err != nil {
		return nil, err
	}

	// Check for commit mode
	if c.mode == commit {
		var data *core_types.ResultBroadcastTxCommit
		err = json.Unmarshal(rpcRes.Result, &data)
		if data.CheckTx.Code != 0 || data.DeliverTx.Code != 0 {
			return nil, errors.New("check tx error: " + data.CheckTx.Log + "; deliver tx error: " + data.DeliverTx.Log)
		}
		return data, nil
	}
	// Check for sync/async modes
	var data *core_types.ResultBroadcastTx
	err = json.Unmarshal(rpcRes.Result, &data)
	if err != nil {
		return nil, err
	}
	if data.Code != 0 {
		return nil, errors.New(data.Log)
	}

	return data, nil
}

func (c *fastClient) addAccount(acc *state.Account) error {
	var err error

	txBytes, err := acc.MarshalMsg(nil)
	if err != nil {
		return err
	}
	tx := transaction.New(transaction.AccountAdd, c.accountID, txBytes)
	bs, _ := tx.ToBytes()

	_, err = c.broadcastTx(bs)
	if err != nil {
		return err
	}
	return nil
}

func (c *fastClient) getAccount(id string) (*state.Account, error) {
	resp, err := c.abciQuery("accounts", []byte(id))
	if err != nil {
		return nil, err
	}
	acc := &state.Account{}
	if err := json.Unmarshal(resp.Response.GetValue(), &acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (c *fastClient) searchAccounts(searchQuery []byte) ([]state.Account, error) {
	resp, err := c.abciQuery("accounts/search", searchQuery)
	if err != nil {
		return nil, err
	}
	acc := []state.Account{}
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

func (c *fastClient) getPayload(id string) (*state.Payload, error) {
	resp, err := c.abciQuery("payloads", []byte(id))
	if err != nil {
		return nil, err
	}
	res := &state.Payload{}
	if err := json.Unmarshal(resp.Response.GetValue(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *fastClient) searchPayloads(searchQuery []byte) ([]state.Payload, error) {
	resp, err := c.abciQuery("payloads/search", searchQuery)
	if err != nil {
		return nil, err
	}
	res := []state.Payload{}
	if err := json.Unmarshal(resp.Response.GetValue(), &res); err != nil {
		return nil, err
	}
	return res, nil
}
