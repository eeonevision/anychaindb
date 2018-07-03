package tests

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/leadschain/leadschain/crypto"
)

var privateKey = flag.String("priv", "CHFHcLSe/eyDb1vLgX5LX11RQhe8D4l9zTFLnYXTTk0=", "ecdsa private key")
var publicKey = flag.String("pub", "BGh8r8qgPotDJb/TYwh/0GjMwW5Y7S4JvPqTTmVH5NkioI9WZt01gafkXXFjAJZRy5ix4vbFRABbHs68KO1LJyw=", "ecdsa public key")
var msg = flag.String("msg", "test message!", "test message")
var update = flag.Bool("update", false, "update .golden files")

func TestEncrypt(t *testing.T) {
	key, err := crypto.NewFromStrings(*publicKey, *privateKey)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	encrypted, err := key.Encrypt([]byte(*msg))
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	// Write encoded message to crypto.golden file
	if *update {
		// Create folder if not exists
		_ = os.Mkdir("artifacts", os.ModePerm)
		// Encode in base64
		res := base64.StdEncoding.EncodeToString([]byte(encrypted))
		// Write base64 encoded message
		err = ioutil.WriteFile("artifacts/crypto.golden", []byte(res), 0644)
		if err != nil {
			t.Errorf("%s", err)
			return
		}
	}
}

func TestDecrypt(t *testing.T) {
	// Build key from strings
	key, err := crypto.NewFromStrings(*publicKey, *privateKey)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	// Read crypto.golden file
	file, err := ioutil.ReadFile("artifacts/crypto.golden")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	// Decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(string(file))
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	// Decrypt message by ECDH algorithm
	decrypted, err := key.Decrypt(decoded)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	// Check if messages are not equals
	if string(decrypted) != *msg {
		t.Errorf("messages are not equals. Expected: %s, Output: %s", *msg, string(decrypted))
		return
	}
}
