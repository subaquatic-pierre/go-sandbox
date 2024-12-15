package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"go-chain/types"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (privKey *PrivateKey) Sign(msg []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey.key, msg)

	if err != nil {
		return nil, fmt.Errorf("unable to sign message, %s", err)
	}

	return &Signature{
		r: r,
		s: s,
	}, nil
}

func NewPrivateKey() PrivateKey {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("unable to generate key, %s", err))
	}
	return PrivateKey{key: priv}
}

func (k *PrivateKey) PublicKey() PublicKey {
	return PublicKey{key: &k.key.PublicKey}
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

func (pubKey *PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(elliptic.P256(), pubKey.key.X, pubKey.key.Y)
}

func (pubKey *PublicKey) Address() types.Address {
	hex := sha256.Sum256(pubKey.ToSlice())

	return types.AddressFromBytes(hex[len(hex)-20:])
}

type Signature struct {
	s, r *big.Int
}

func (sig *Signature) Verify(pubKey PublicKey, msg []byte) bool {

	return ecdsa.Verify(pubKey.key, msg, sig.r, sig.s)
}
