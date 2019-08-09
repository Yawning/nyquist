// Package nyquist implements the Noise Protocol Framework.
package nyquist

import "errors"

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

	// ErrMalformedPrivateKey is the error returned when a serialized
	// private key is malformed.
	ErrMalformedPrivateKey = errors.New("nyquist: malformed private key")

	// ErrMalformedPublicKey is the error returned when a serialized public
	// key is malformed.
	ErrMalformedPublicKey = errors.New("nyquist: malformed public key")

	// ErrMismatchedPublicKey is the error returned when a public key for an
	// unexpected algorithm is provided to a DH calculation.
	ErrMismatchedPublicKey = errors.New("nyquist: mismatched public key")

	// ErrInvalidConfig is the error returned when the configuration is invalid.
	ErrInvalidConfig = errors.New("nyquist: invalid configuration")

	// ErrOutOfOrder is the error returned when ReadMessage/WriteMessage
	// are called out of order, given the handshake's initiator status.
	ErrOutOfOrder = errors.New("nyquist: out of order handshake operation")

	// ErrDone is the error returned when further WrtieMessage/ReadMessage
	// calls are attempted on an already completed HandshakeState.
	ErrDone = errors.New("nyquist: handshake already complete")

	// ErrProtocolNotSupported is the error returned when a requested protocol
	// is not supported.
	ErrProtocolNotSupported = errors.New("nyquist: protocol not supported")
)

// Resetable is the interface implemented by objects capable of sanitizing
// themselves.
//
// Warning: In some cases this is strictly a best-effort process.
type Resetable interface {
	// Reset clears the object of sensitive data.
	Reset()
}

func explicitBzero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func truncateTo32Bytes(b []byte) []byte {
	if len(b) <= 32 {
		return b
	}

	var tail []byte
	b, tail = b[:32], b[32:]
	explicitBzero(tail)

	return b
}
