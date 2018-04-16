package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

type Key struct {
	pub  *ecdsa.PublicKey
	priv *ecdsa.PrivateKey
}

func CreateKeyPair() (*Key, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
