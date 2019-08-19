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

// Package pattern implements the Noise Protocol Framework handshake pattern
// abstract interface and standard patterns.
package pattern // import "gitlab.com/yawning/nyquist.git/pattern"

import "fmt"

var supportedPatterns = make(map[string]Pattern)

// Token is a Noise handshake pattern token.
type Token uint8

const (
	Token_invalid Token = iota
	Token_e
	Token_s
	Token_ee
	Token_es
	Token_se
	Token_ss
	Token_psk
)

// String returns the string representation of a Token.
func (t Token) String() string {
	switch t {
	case Token_e:
		return "e"
	case Token_s:
		return "s"
	case Token_ee:
		return "ee"
	case Token_es:
		return "es"
	case Token_se:
		return "se"
	case Token_ss:
		return "ss"
	case Token_psk:
		return "psk"
	default:
		return fmt.Sprintf("[invalid token: %d]", int(t))
	}
}

// Message is a sequence of pattern tokens.
type Message []Token

// Pattern is a handshake pattern.
type Pattern interface {
	fmt.Stringer

	// PreMessages returns the pre-message message patterns.
	PreMessages() []Message

	// Mesages returns the message patterns.
	Messages() []Message

	// NumPSKs returns the number of `psk` modifiers in the pattern.
	NumPSKs() int

	// IsOneWay returns true iff the pattern is one-way.
	IsOneWay() bool
}

// FromString returns a Pattern by pattern name, or nil.
func FromString(s string) Pattern {
	return supportedPatterns[s]
}

type builtIn struct {
	name        string
	preMessages []Message
	messages    []Message
	numPSKs     int
	isOneWay    bool
}

func (pa *builtIn) String() string {
	return pa.name
}

func (pa *builtIn) PreMessages() []Message {
	return pa.preMessages
}

func (pa *builtIn) Messages() []Message {
	return pa.messages
}

func (pa *builtIn) NumPSKs() int {
	return pa.numPSKs
}

func (pa *builtIn) IsOneWay() bool {
	return pa.isOneWay
}

// Register registers a new pattern for use with `FromString()`.
func Register(pa Pattern) error {
	if err := IsValid(pa); err != nil {
		return err
	}
	supportedPatterns[pa.String()] = pa

	return nil
}

func init() {
	for _, v := range []Pattern{
		// One-way patterns.
		N,
		K,
		X,
		Npsk0,
		Kpsk0,
		Xpsk1,

		// Interactive (fundemental) patterns.
		NN,
		NK,
		NX,
		XN,
		XK,
		XX,
		KN,
		KK,
		KX,
		IN,
		IK,
		IX,
		NNpsk0,
		NNpsk2,
		NKpsk0,
		NKpsk2,
		NXpsk2,
		XNpsk3,
		XKpsk3,
		XXpsk3,
		KNpsk0,
		KNpsk2,
		KKpsk0,
		KKpsk2,
		KXpsk2,
		INpsk1,
		INpsk2,
		IKpsk1,
		IKpsk2,
		IXpsk2,

		// Deferred patterns.
		NK1,
		NX1,
		X1N,
		X1K,
		XK1,
		X1K1,
		X1X,
		XX1,
		X1X1,
		K1N,
		K1K,
		KK1,
		K1K1,
		K1X,
		KX1,
		K1X1,
		I1N,
		I1K,
		IK1,
		I1K1,
		I1X,
		IX1,
		I1X1,
	} {
		if err := Register(v); err != nil {
			panic("nyquist/pattern: failed to register built-in pattern: " + err.Error())
		}
	}
}
