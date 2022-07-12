package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func SignatureFromString(s string) *Signature {
	x, y := StringToBigIntTuple(s)

	return &Signature{&x, &y}
}

func StringToBigIntTuple(s string) (bix big.Int, biy big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	bix.SetBytes(bx)
	biy.SetBytes(by)

	return
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := StringToBigIntTuple(s)

	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])
	var bi big.Int

	bi.SetBytes(b)

	return &ecdsa.PrivateKey{*publicKey, &bi}
}
