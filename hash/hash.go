// Package hash implments the Noise Protocol Framework hash function abstract
// interface and standard hash functions.
package hash // import "gitlab.com/yawning/nyquist.git/hash"

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

var (
	// SHA256 is the SHA256 hash function.
	SHA256 Hash = &hashSha256{}

	// SHA512 is the SHA512 hash function.
	SHA512 Hash = &hashSha512{}

	// BLAKE2s is the BLAKE2s hash function.
	BLAKE2s Hash = &hashBlake2s{}

	// BLAKE2b is the BLAKE2b hash function.
	BLAKE2b Hash = &hashBlake2b{}

	supportedHashes = map[string]Hash{
		"SHA256":  SHA256,
		"SHA512":  SHA512,
		"BLAKE2s": BLAKE2s,
		"BLAKE2b": BLAKE2b,
	}
)

// Hash is a collision-resistant cryptographic hash function factory.
type Hash interface {
	fmt.Stringer

	// New constructs a new `hash.Hash` instance.
	New() hash.Hash

	// Size returns the hash function's digest size in bytes (`HASHLEN`).
	Size() int
}

// FromString returns a Hash by algorithm name, or nil.
func FromString(s string) Hash {
	return supportedHashes[s]
}

type hashSha256 struct{}

func (h *hashSha256) String() string {
	return "SHA256"
}

func (h *hashSha256) New() hash.Hash {
	return sha256.New()
}

func (h *hashSha256) Size() int {
	return sha256.Size
}

type hashSha512 struct{}

func (h *hashSha512) String() string {
	return "SHA512"
}

func (h *hashSha512) New() hash.Hash {
	return sha512.New()
}

func (h *hashSha512) Size() int {
	return sha512.Size
}

type hashBlake2s struct{}

func (h *hashBlake2s) String() string {
	return "BLAKE2s"
}

func (h *hashBlake2s) New() hash.Hash {
	ret, _ := blake2s.New256(nil)
	return ret
}

func (h *hashBlake2s) Size() int {
	return blake2s.Size
}

type hashBlake2b struct{}

func (h *hashBlake2b) String() string {
	return "BLAKE2b"
}

func (h *hashBlake2b) New() hash.Hash {
	ret, _ := blake2b.New512(nil)
	return ret
}

func (h *hashBlake2b) Size() int {
	return blake2b.Size
}
