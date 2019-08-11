// Package nyquist implements the Noise Protocol Framework.
package nyquist

import (
	"errors"

	"gitlab.com/yawning/nyquist.git/internal"
)

// Version is the revision of the Noise specification implemented.
const Version = 34

var (
	// ErrNonceExhausted is the error returned when the CipherState's
	// nonce space is exhausted.
	ErrNonceExhausted = errors.New("nyquist: nonce exhausted")

	// ErrMessageSize is the error returned when an operation fails due
	// to the message size being exceeded.
	ErrMessageSize = errors.New("nyquist: oversized message")

	// ErrOpen is the error returned on a authenticated decryption failure.
	ErrOpen = errors.New("nyquist: decryption failure")

	// ErrInvalidConfig is the error returned when the configuration is invalid.
	ErrInvalidConfig = errors.New("nyquist: invalid configuration")

	// ErrOutOfOrder is the error returned when ReadMessage/WriteMessage
	// are called out of order, given the handshake's initiator status.
	ErrOutOfOrder = errors.New("nyquist: out of order handshake operation")

	// ErrDone is the error returned when the handshake is complete.
	ErrDone = errors.New("nyquist: handshake complete")

	// ErrProtocolNotSupported is the error returned when a requested protocol
	// is not supported.
	ErrProtocolNotSupported = errors.New("nyquist: protocol not supported")
)

func truncateTo32Bytes(b []byte) []byte {
	if len(b) <= 32 {
		return b
	}

	var tail []byte
	b, tail = b[:32], b[32:]
	internal.ExplicitBzero(tail)

	return b
}
