package nyquist

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"gitlab.com/yawning/bsaes.git"
	"golang.org/x/crypto/chacha20poly1305"
)

var supportedCiphers = map[string]Cipher{
	"ChaChaPoly": ChaChaPoly,
	"AESGCM":     AESGCM,
}

// Cipher is an AEAD algorithm factory.
type Cipher interface {
	fmt.Stringer

	// New constructs a new keyed `cipher.AEAD` instance, with the provided
	// key.
	New(key []byte) (cipher.AEAD, error)

	// EncodeNonce encodes a Noise nonce to a nonce suitable for use with
	// the `cipher.AEAD` instances created by `Cipher.New`.
	EncodeNonce(nonce uint64) []byte
}

// Rekeyable is the interface implemented by Cipher instances that have a
// `REKEY(k)` function specifically defined.
type Rekeyable interface {
	// Rekey returns a new 32-byte cipher key as a pseudorandom function of `k`.
	Rekey(k []byte) []byte
}

// ChaChaPoly is the ChaChaPoly cipher functions.
//
// Note: Due to upstream limitiations, key sanitization is currently not
// supported.
var ChaChaPoly Cipher = &cipherChaChaPoly{}

type cipherChaChaPoly struct{}

func (ci *cipherChaChaPoly) String() string {
	return "ChaChaPoly"
}

func (ci *cipherChaChaPoly) New(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.New(key)
}

func (ci *cipherChaChaPoly) EncodeNonce(nonce uint64) []byte {
	var encodedNonce [12]byte // 96 bits
	binary.LittleEndian.PutUint64(encodedNonce[4:], nonce)
	return encodedNonce[:]
}

// AESGCM is the AESGCM cipher functions.
//
// Note: This Cipher implementation is always constant time, even on systems
// where the Go runtime library's is not.  Due to runrime library limitiations,
// key sanitization is currently not universally supported.
var AESGCM Cipher = &cipherAesGcm{}

type cipherAesGcm struct{}

func (ci *cipherAesGcm) String() string {
	return "AESGCM"
}

func (ci *cipherAesGcm) New(key []byte) (cipher.AEAD, error) {
	block, err := bsaes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func (ci *cipherAesGcm) EncodeNonce(nonce uint64) []byte {
	var encodedNonce [12]byte // 96 bits
	binary.BigEndian.PutUint64(encodedNonce[4:], nonce)
	return encodedNonce[:]
}
