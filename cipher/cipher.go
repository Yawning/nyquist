// Copryright (C) 2019 Yawning Angel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package cipher implments the Noise Protocol Framework cipher function
// abstract interface and standard cipher functions.
package cipher // import "gitlab.com/yawning/nyquist.git/cipher"

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"github.com/oasislabs/deoxysii"
	"gitlab.com/yawning/bsaes.git"
	"golang.org/x/crypto/chacha20poly1305"
)

var supportedCiphers = map[string]Cipher{
	"ChaChaPoly": ChaChaPoly,
	"AESGCM":     AESGCM,
	"DeoxysII":   DeoxysII,
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

// Resetable is the interface implemented by `crypto/cipher.AEAD` instances
// that are capable of sanitizing themselves.
type Resetable interface {
	// Reset clears the object of sensitive data.
	Reset()
}

// FromString returns a Cipher by algorithm name, or nil.
func FromString(s string) Cipher {
	return supportedCiphers[s]
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

// DeoxysII is the DeoxysII cipher functions.
//
// Warning: This cipher is non-standard.
var DeoxysII Cipher = &cipherDeoxysII{}

type cipherDeoxysII struct{}

func (ci *cipherDeoxysII) String() string {
	return "DeoxysII"
}

func (ci *cipherDeoxysII) New(key []byte) (cipher.AEAD, error) {
	return deoxysii.New(key)
}

func (ci *cipherDeoxysII) EncodeNonce(nonce uint64) []byte {
	// Using the full nonce-space is fine, and big endian follows how
	// Deoxys-II encodes things internally.
	var encodedNonce [deoxysii.NonceSize]byte // 120 bits
	binary.BigEndian.PutUint64(encodedNonce[7:], nonce)
	return encodedNonce[:]
}

// Register registers a new cipher for use with `FromString()`.
func Register(cipher Cipher) {
	supportedCiphers[cipher.String()] = cipher
}
