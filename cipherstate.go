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

package nyquist

import (
	goCipher "crypto/cipher"
	"errors"
	"math"

	"gitlab.com/yawning/nyquist.git/cipher"
	"gitlab.com/yawning/nyquist.git/internal"
)

const (
	// SymmetricKeySize is the size a symmetric key in bytes.
	SymmetricKeySize = 32

	maxnonce = math.MaxUint64
)

var (
	errInvalidKeySize = errors.New("nyquist/CipherState: invalid key size")
	errNoExistingKey  = errors.New("nyquist/CipherState: failed to rekey, no existing key")

	zeroes [32]byte
)

// CipherState is a keyed AEAD algorithm instance.
type CipherState struct {
	cipher cipher.Cipher

	aead goCipher.AEAD
	k    []byte
	n    uint64

	maxMessageSize int
	aeadOverhead   int
}

// InitializeKey initializes sets the cipher key to `key`, and nonce to 0.
func (cs *CipherState) InitializeKey(key []byte) {
	if err := cs.setKey(key); err != nil {
		panic("nyquist/CipherState: failed to initialize key: " + err.Error())
	}
	cs.n = 0
}

func (cs *CipherState) setKey(key []byte) error {
	cs.Reset()

	switch len(key) {
	case 0:
	case SymmetricKeySize:
		var err error
		if cs.aead, err = cs.cipher.New(key); err != nil {
			return err
		}

		cs.aeadOverhead = cs.aead.Overhead()

		cs.k = make([]byte, SymmetricKeySize)
		copy(cs.k, key)
	default:
		return errInvalidKeySize
	}

	return nil
}

// HasKey returns true iff the CipherState is keyed.
func (cs *CipherState) HasKey() bool {
	return cs.aead != nil
}

// SetNonce sets the CipherState's nonce to `nonce`.
func (cs *CipherState) SetNonce(nonce uint64) {
	cs.n = nonce
}

// EncryptWithAd encrypts and authenticates the additional data and plaintext
// and increments the nonce iff the CipherState is keyed, and otherwise returns
// the plaintext.
//
// Note: The ciphertext is appended to `dst`, and the new slice is returned.
func (cs *CipherState) EncryptWithAd(dst, ad, plaintext []byte) ([]byte, error) {
	aead := cs.aead
	if aead == nil {
		return append(dst, plaintext...), nil
	}

	if cs.n == maxnonce {
		return nil, ErrNonceExhausted
	}

	if cs.maxMessageSize > 0 && len(plaintext)+cs.aeadOverhead > cs.maxMessageSize {
		return nil, ErrMessageSize
	}

	nonce := cs.cipher.EncodeNonce(cs.n)
	ciphertext := aead.Seal(dst, nonce, plaintext, ad)
	cs.n++

	return ciphertext, nil
}

// DecryptWihtAd authenticates and decrypts the additional data and ciphertext
// and increments the nonce iff the CipherState is keyed, and otherwise returns
// the plaintext.  If an authentication failure occurs, the nonce is not
// incremented.
//
// Note: The plaintext is appended to `dst`, and the new slice is returned.
func (cs *CipherState) DecryptWithAd(dst, ad, ciphertext []byte) ([]byte, error) {
	aead := cs.aead
	if aead == nil {
		return append(dst, ciphertext...), nil
	}

	if cs.n == maxnonce {
		return nil, ErrNonceExhausted
	}

	if cs.maxMessageSize > 0 && len(ciphertext) > cs.maxMessageSize {
		return nil, ErrMessageSize
	}

	nonce := cs.cipher.EncodeNonce(cs.n)
	plaintext, err := aead.Open(dst, nonce, ciphertext, ad)
	if err != nil {
		return nil, ErrOpen
	}
	cs.n++

	return plaintext, nil
}

// Rekey sets the CipherState's key to `REKEY(k)`.
func (cs *CipherState) Rekey() error {
	if !cs.HasKey() {
		return errNoExistingKey
	}

	var newKey []byte
	if rekeyer, ok := (cs.cipher).(cipher.Rekeyable); ok {
		// The cipher function set has a specific `REKEY` function defined.
		newKey = rekeyer.Rekey(cs.k)
	} else {
		// The cipher function set has no `REKEY` function defined, use the
		// default generic implementation.
		nonce := cs.cipher.EncodeNonce(maxnonce)
		newKey = cs.aead.Seal(nil, nonce, zeroes[:], nil)

		// "defaults to returning the first 32 bytes"
		newKey = truncateTo32Bytes(newKey)
	}

	err := cs.setKey(newKey)
	internal.ExplicitBzero(newKey)

	return err
}

// Reset clears the CipherState of sensitive data.
func (cs *CipherState) Reset() {
	if cs.k != nil {
		internal.ExplicitBzero(cs.k)
		cs.k = nil
	}
	if cs.aead != nil {
		if resetter, ok := (cs.aead).(cipher.Resetable); ok {
			resetter.Reset()
		}
		cs.aead = nil
		cs.aeadOverhead = 0
	}
}

func newCipherState(cipher cipher.Cipher, maxMessageSize int) *CipherState {
	return &CipherState{
		cipher:         cipher,
		maxMessageSize: maxMessageSize,
	}
}
