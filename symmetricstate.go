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
	"io"

	"golang.org/x/crypto/hkdf"

	"gitlab.com/yawning/nyquist.git/cipher"
	"gitlab.com/yawning/nyquist.git/hash"
	"gitlab.com/yawning/nyquist.git/internal"
)

// SymmetricState encapsulates all symmetric cryptography used by the Noise
// protocol during a handshake.
//
// Warning: There should be no reason to interact directly with this ever.
type SymmetricState struct {
	cipher cipher.Cipher
	hash   hash.Hash

	cs *CipherState

	ck []byte
	h  []byte

	hashLen int
}

// InitializeSymmetric initializes the SymmetricSatate with the initial
// chaining key and handshake hash, based on the protocol name.
func (ss *SymmetricState) InitializeSymmetric(protocolName []byte) {
	if len(protocolName) <= ss.hashLen {
		ss.h = make([]byte, ss.hashLen)
		copy(ss.h, protocolName)
	} else {
		h := ss.hash.New()
		_, _ = h.Write(protocolName)
		ss.h = h.Sum(nil)
	}

	ss.ck = make([]byte, 0, ss.hashLen)
	ss.ck = append(ss.ck, ss.h...)

	ss.cs.InitializeKey(nil)
}

// MixKey mixes the provided material with the chaining key, and initializes
// the encapsulated CipherState's key with the output.
func (ss *SymmetricState) MixKey(inputKeyMaterial []byte) {
	tempK := make([]byte, ss.hashLen)

	ss.hkdfHash(inputKeyMaterial, ss.ck, tempK)
	tempK = truncateTo32Bytes(tempK)
	ss.cs.InitializeKey(tempK)

	internal.ExplicitBzero(tempK)
}

// MixHash mixes the provided data with the handshake hash.
func (ss *SymmetricState) MixHash(data []byte) {
	h := ss.hash.New()
	_, _ = h.Write(ss.h)
	_, _ = h.Write(data)
	ss.h = h.Sum(ss.h[:0])

	h.Reset()
}

// MixKeyAndHash mises the provided material with the chaining key, and mixes
// the handshake and initializes the encapsulated CipherState with the output.
func (ss *SymmetricState) MixKeyAndHash(inputKeyMaterial []byte) {
	tempH, tempK := make([]byte, ss.hashLen), make([]byte, ss.hashLen)

	ss.hkdfHash(inputKeyMaterial, ss.ck, tempH, tempK)
	ss.MixHash(tempH)
	tempK = truncateTo32Bytes(tempK)
	ss.cs.InitializeKey(tempK)

	internal.ExplicitBzero(tempK)
}

// GetHandshakeHash returns the handshake hash `h`.
func (ss *SymmetricState) GetHandshakeHash() []byte {
	return ss.h
}

// EncryptAndHash encrypts and authenticates the plaintext, mixes the
// ciphertext with the handshake hash, appends the ciphertext to dst,
// and returns the potentially new slice.
func (ss *SymmetricState) EncryptAndHash(dst, plaintext []byte) []byte {
	var err error
	ciphertextOff := len(dst)
	if dst, err = ss.cs.EncryptWithAd(dst, ss.h, plaintext); err != nil {
		panic("nyquist/SymmetricState: encryptAndHash() failed: " + err.Error())
	}
	ss.MixHash(dst[ciphertextOff:])
	return dst
}

// DecryptAndHash authenticates and decrypts the ciphertext, mixes the
// ciphertext with the handshake hash, appends the plaintext to dst,
// and returns the potentially new slice.
func (ss *SymmetricState) DecryptAndHash(dst, ciphertext []byte) ([]byte, error) {
	// `dst` and `ciphertext` could alias, so save a copy of `h` so that the
	// `MixHash()` call can be called prior to `DecryptWithAd`.
	hPrev := make([]byte, 0, len(ss.h))
	hPrev = append(hPrev, ss.h...)

	ss.MixHash(ciphertext)

	return ss.cs.DecryptWithAd(dst, hPrev, ciphertext)
}

// Split returns a pair of CipherState objects for encrypted transport messages.
func (ss *SymmetricState) Split() (*CipherState, *CipherState) {
	tempK1, tempK2 := make([]byte, ss.hashLen), make([]byte, ss.hashLen)

	ss.hkdfHash(nil, tempK1, tempK2)
	tempK1 = truncateTo32Bytes(tempK1)
	tempK2 = truncateTo32Bytes(tempK2)

	c1, c2 := newCipherState(ss.cipher, ss.cs.maxMessageSize), newCipherState(ss.cipher, ss.cs.maxMessageSize)
	c1.InitializeKey(tempK1)
	c2.InitializeKey(tempK2)

	internal.ExplicitBzero(tempK1)
	internal.ExplicitBzero(tempK2)

	return c1, c2
}

// CipherState returns the SymmetricState's encapsualted CipherState.
//
// Warning: There should be no reason to call this, ever.
func (ss *SymmetricState) CipherState() *CipherState {
	return ss.cs
}

func (ss *SymmetricState) hkdfHash(inputKeyMaterial []byte, outputs ...[]byte) {
	// There is no way to sanitize the HKDF reader state.  While it is tempting
	// to just write a HKDF implementation that supports sanitization, neither
	// `crypto/hmac` nor the actual hash function implementations support
	// sanitization correctly either due to:
	//
	//  * `Reset()`ing a HMAC instance resets it to the keyed (initialized)
	//     state.
	//  * All of the concrete hash function implementations do not `Reset()`
	//    the cloned instance when `Sum([]byte)` is called.

	r := hkdf.New(ss.hash.New, inputKeyMaterial, ss.ck, nil)
	for _, output := range outputs {
		if len(output) != ss.hashLen {
			panic("nyquist/SymmetricState: non-HASHLEN sized output to HKDF-HASH")
		}
		_, _ = io.ReadFull(r, output)
	}
}

// Reset clears the SymmetricState of sensitive data.
func (ss *SymmetricState) Reset() {
	// `ss.h` is not sensitive, and not explicitly clearing it allows
	// `Reset()` to be called immediately after `Split()` while preserving
	// the ability to call `GetHandshakeHash()`.
	if ss.ck != nil {
		internal.ExplicitBzero(ss.ck)
		ss.ck = nil
	}
	if ss.cs != nil {
		ss.cs.Reset()
		ss.cs = nil
	}
}

func newSymmetricState(cipher cipher.Cipher, hash hash.Hash, maxMessageSize int) *SymmetricState {
	return &SymmetricState{
		cipher:  cipher,
		hash:    hash,
		cs:      newCipherState(cipher, maxMessageSize),
		hashLen: hash.Size(),
	}
}
