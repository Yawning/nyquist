// Package pattern implements the Noise Protocol Framework handshake pattern
// abstract interface and standard patterns.
package pattern

import "fmt"

var supportedPatterns = map[string]Pattern{
	// One-way patterns.
	"N":     N,
	"K":     K,
	"X":     X,
	"Npsk0": Npsk0,
	"Kpsk0": Kpsk0,
	"Xpsk1": Xpsk1,

	// Interactive (fundemental) patterns.
	"NN":     NN,
	"NK":     NK,
	"NX":     NX,
	"XN":     XN,
	"XK":     XK,
	"XX":     XX,
	"KN":     KN,
	"KK":     KK,
	"KX":     KX,
	"IN":     IN,
	"IK":     IK,
	"IX":     IX,
	"NNpsk0": NNpsk0,
	"NNpsk2": NNpsk2,
	"NKpsk0": NKpsk0,
	"NKpsk2": NKpsk2,
	"NXpsk2": NXpsk2,
	"XNpsk3": XNpsk3,
	"XKpsk3": XKpsk3,
	"XXpsk3": XXpsk3,
	"KNpsk0": KNpsk0,
	"KNpsk2": KNpsk2,
	"KKpsk0": KKpsk0,
	"KKpsk2": KKpsk2,
	"KXpsk2": KXpsk2,
	"INpsk1": INpsk1,
	"INpsk2": INpsk2,
	"IKpsk1": IKpsk1,
	"IKpsk2": IKpsk2,
	"IXpsk2": IXpsk2,

	// Deferred patterns.
	"NK1":  NK1,
	"NX1":  NX1,
	"X1N":  X1N,
	"X1K":  X1K,
	"XK1":  XK1,
	"X1K1": X1K1,
	"X1X":  X1X,
	"XX1":  XX1,
	"X1X1": X1X1,
	"K1N":  K1N,
	"K1K":  K1K,
	"KK1":  KK1,
	"K1K1": K1K1,
	"K1X":  K1X,
	"KX1":  KX1,
	"K1X1": K1X1,
	"I1N":  I1N,
	"I1K":  I1K,
	"IK1":  IK1,
	"I1K1": I1K1,
	"I1X":  I1X,
	"IX1":  IX1,
	"I1X1": I1X1,
}

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

	// IsPSK returns true iff the pattern has a `psk` modifier.
	IsPSK() bool

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
	isPSK       bool
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

func (pa *builtIn) IsPSK() bool {
	return pa.isPSK
}

func (pa *builtIn) IsOneWay() bool {
	return pa.isOneWay
}
