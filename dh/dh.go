// Package dh implments the Noise Protocol Framework Diffie-Hellman function
// abstract interface and standard hash functions.
package dh

import (
	"encoding"
	"errors"
	"fmt"
	"io"

	"github.com/oasislabs/ed25519/extra/x25519"

	"gitlab.com/yawning/nyquist.git/internal"
)

var (
	// ErrMalformedPrivateKey is the error returned when a serialized
	// private key is malformed.
	ErrMalformedPrivateKey = errors.New("nyquist/dh: malformed private key")

	// ErrMalformedPublicKey is the error returned when a serialized public
	// key is malformed.
	ErrMalformedPublicKey = errors.New("nyquist/dh: malformed public key")

	// ErrMismatchedPublicKey is the error returned when a public key for an
	// unexpected algorithm is provided to a DH calculation.
	ErrMismatchedPublicKey = errors.New("nyquist/dh: mismatched public key")

	supportedDHs = map[string]DH{
		"25519": X25519,
	}
)

// DH is a Diffie-Hellman key exchange algorithm.
type DH interface {
	fmt.Stringer

	// GenerateKeypair generates a new Diffie-Hellman keypair using the
	// provided entropy source.
	GenerateKeypair(rng io.Reader) (Keypair, error)

	// ParsePublicKey parses a binary encoded public key.
	ParsePublicKey(data []byte) (PublicKey, error)

	// Size returns the size of public keys and DH outputs in bytes (`DHLEN`).
	Size() int
}

// FromString returns a DH by algorithm name, or nil.
func FromString(s string) DH {
	return supportedDHs[s]
}

// Keypair is a Diffie-Hellman keypair.
type Keypair interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Reset clears the object of sensitive data.
	Reset()

	// Public returns the public key of the keypair.
	Public() PublicKey

	// DH performs a Diffie-Hellman calculation between the private key
	// in the keypair and the provided public key.
	DH(publicKey PublicKey) ([]byte, error)
}

// PublicKey is a Diffie-Hellman public key.
type PublicKey interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Bytes returns the binary serialized public key.
	//
	// Warning: Altering the returned slice is unsupported and will lead
	// to unexpected behavior.
	Bytes() []byte
}

// X25519 is the 25519 DH function.
var X25519 DH = &dh25519{}

type dh25519 struct{}

func (dh *dh25519) String() string {
	return "25519"
}

func (dh *dh25519) GenerateKeypair(rng io.Reader) (Keypair, error) {
	var kp Keypair25519
	if _, err := io.ReadFull(rng, kp.rawPrivateKey[:]); err != nil {
		return nil, err
	}

	x25519.ScalarBaseMult(&kp.publicKey.rawPublicKey, &kp.rawPrivateKey)

	return &kp, nil
}

func (dh *dh25519) ParsePublicKey(data []byte) (PublicKey, error) {
	var pk PublicKey25519
	if err := pk.UnmarshalBinary(data); err != nil {
		return nil, err
	}

	return &pk, nil
}

func (dh *dh25519) Size() int {
	return 32
}

// Keypair25519 is a X25519 keypair.
type Keypair25519 struct {
	rawPrivateKey [32]byte
	publicKey     PublicKey25519
}

// MarshalBinary marshals the keypair's private key to binary form.
func (kp *Keypair25519) MarshalBinary() ([]byte, error) {
	out := make([]byte, 0, len(kp.rawPrivateKey))
	return append(out, kp.rawPrivateKey[:]...), nil
}

// UnmarshalBinary unmarshals the keypair's private key from binary form,
// and re-derives the corresponding public key.
func (kp *Keypair25519) UnmarshalBinary(data []byte) error {
	if len(data) != 32 {
		return ErrMalformedPrivateKey
	}

	copy(kp.rawPrivateKey[:], data)
	x25519.ScalarBaseMult(&kp.publicKey.rawPublicKey, &kp.rawPrivateKey)

	return nil
}

// Public returns the public key of the keypair.
func (kp *Keypair25519) Public() PublicKey {
	return &kp.publicKey
}

// DH performs a Diffie-Hellman calculation between the private key in the
// keypair and the provided public key.
func (kp *Keypair25519) DH(publicKey PublicKey) ([]byte, error) {
	pubKey, ok := publicKey.(*PublicKey25519)
	if !ok {
		return nil, ErrMismatchedPublicKey
	}

	var sharedSecret [32]byte
	x25519.ScalarMult(&sharedSecret, &kp.rawPrivateKey, &pubKey.rawPublicKey)

	return sharedSecret[:], nil
}

// Reset clears the keypair of sensitive data.
func (kp *Keypair25519) Reset() {
	internal.ExplicitBzero(kp.rawPrivateKey[:])
}

// PublicKey25519 is a X25519 public key.
type PublicKey25519 struct {
	rawPublicKey [32]byte
}

// MarshalBinary marshals the public key to binary form.
func (pk *PublicKey25519) MarshalBinary() ([]byte, error) {
	out := make([]byte, 0, len(pk.rawPublicKey))
	return append(out, pk.rawPublicKey[:]...), nil
}

// UnmarshalBinary unmarshals the public key from binary form.
func (pk *PublicKey25519) UnmarshalBinary(data []byte) error {
	if len(data) != 32 {
		return ErrMalformedPublicKey
	}

	copy(pk.rawPublicKey[:], data)

	return nil
}

// Bytes returns the binary serialized public key.
//
// Warning: Altering the returned slice is unsupported and will lead to
// unexpected behavior.
func (pk *PublicKey25519) Bytes() []byte {
	return pk.rawPublicKey[:]
}
