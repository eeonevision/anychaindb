package crypto

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"

	"github.com/cloudflare/redoctober/padding"
	"github.com/cloudflare/redoctober/symcrypt"
)

var Curve = elliptic.P256

type Key struct {
	pub  *ecdsa.PublicKey
	priv *ecdsa.PrivateKey
}

func CreateKeyPair() (*Key, error) {
	priv, err := ecdsa.GenerateKey(Curve(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{priv: priv, pub: &priv.PublicKey}, nil
}

func NewFromStrings(pub, priv string) (*Key, error) {
	k := &Key{}
	if pub == "" && priv == "" {
		return nil, errors.New("No key material supplied")
	}
	if pub != "" {
		if err := k.SetPubString(pub); err != nil {
			return nil, err
		}
	}
	if priv != "" && pub == "" {
		return nil, errors.New("No pubkey to privkey supplied")
	}
	if priv != "" {
		if err := k.SetPrivString(priv); err != nil {
			return nil, err
		}
	}
	return k, nil
}

func (k *Key) GetPubString() string {
	pub := elliptic.Marshal(elliptic.P256(), k.pub.X, k.pub.Y)
	return base64.StdEncoding.EncodeToString(pub)
}

func (k *Key) GetPrivString() string {
	priv := k.priv.D.Bytes()
	return base64.StdEncoding.EncodeToString(priv)
}

func (k *Key) SetPubString(pub string) error {
	bs, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), bs)
	k.pub = &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	return nil
}

func (k *Key) SetPrivString(priv string) error {
	bs, err := base64.StdEncoding.DecodeString(priv)
	if err != nil {
		return err
	}
	d := &big.Int{}
	d.SetBytes(bs)
	k.priv = &ecdsa.PrivateKey{D: d, PublicKey: *k.pub}
	k.pub = &k.priv.PublicKey
	return nil
}

func (k *Key) Sign(hash []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.priv, hash)
	if err != nil {
		return "", err
	}
	rStr := base64.StdEncoding.EncodeToString(r.Bytes())
	sStr := base64.StdEncoding.EncodeToString(s.Bytes())
	return rStr + ":" + sStr, nil
}

func (k *Key) Verify(hash []byte, signature string) error {
	parts := strings.Split(signature, ":")
	if len(parts) != 2 {
		return errors.New("Malformed signature")
	}
	rBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	sBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	r := (&big.Int{}).SetBytes(rBytes)
	s := (&big.Int{}).SetBytes(sBytes)
	if !ecdsa.Verify(k.pub, hash, r, s) {
		return errors.New("Bad signature")
	}
	return nil
}

// Encrypt secures and authenticates its input using the public key
// using ECDHE with AES-128-CBC-HMAC-SHA1.
func (k *Key) Encrypt(in []byte) (out []byte, err error) {
	ephemeral, err := ecdsa.GenerateKey(Curve(), rand.Reader)
	if err != nil {
		return
	}
	x, _ := k.pub.Curve.ScalarMult(k.pub.X, k.pub.Y, ephemeral.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}
	shared := sha256.Sum256(x.Bytes())
	iv, err := symcrypt.MakeRandom(16)
	if err != nil {
		return
	}

	paddedIn := padding.AddPadding(in)
	ct, err := symcrypt.EncryptCBC(paddedIn, iv, shared[:16])
	if err != nil {
		return
	}

	ephPub := elliptic.Marshal(k.pub.Curve, ephemeral.PublicKey.X, ephemeral.PublicKey.Y)
	out = make([]byte, 1+len(ephPub)+16)
	out[0] = byte(len(ephPub))
	copy(out[1:], ephPub)
	copy(out[1+len(ephPub):], iv)
	out = append(out, ct...)

	h := hmac.New(sha1.New, shared[16:])
	h.Write(iv)
	h.Write(ct)
	out = h.Sum(out)
	return
}

// Decrypt authentications and recovers the original message from
// its input using the private key and the ephemeral key included in
// the message.
func (k *Key) Decrypt(in []byte) (out []byte, err error) {
	ephLen := int(in[0])
	ephPub := in[1 : 1+ephLen]
	ct := in[1+ephLen:]
	if len(ct) < (sha1.Size + aes.BlockSize) {
		return nil, errors.New("Invalid ciphertext")
	}

	x, y := elliptic.Unmarshal(Curve(), ephPub)
	if x == nil {
		return nil, errors.New("Invalid public key")
	}

	x, _ = k.priv.Curve.ScalarMult(x, y, k.priv.D.Bytes())
	if x == nil {
		return nil, errors.New("Failed to generate encryption key")
	}
	shared := sha256.Sum256(x.Bytes())

	tagStart := len(ct) - sha1.Size
	h := hmac.New(sha1.New, shared[16:])
	h.Write(ct[:tagStart])
	mac := h.Sum(nil)
	if !hmac.Equal(mac, ct[tagStart:]) {
		return nil, errors.New("Invalid MAC")
	}

	paddedOut, err := symcrypt.DecryptCBC(ct[aes.BlockSize:tagStart], ct[:aes.BlockSize], shared[:16])
	if err != nil {
		return
	}
	out, err = padding.RemovePadding(paddedOut)
	return
}
