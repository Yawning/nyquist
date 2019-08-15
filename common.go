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

// Package nyquist implements the Noise Protocol Framework.
package nyquist // import "gitlab.com/yawning/nyquist.git"

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
