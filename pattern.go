package nyquist

import "fmt"

var supportedPatterns = map[string]HandshakePattern{
	// One-way patterns.
	"N":     Pattern_N,
	"K":     Pattern_K,
	"X":     Pattern_X,
	"Npsk0": Pattern_Npsk0,
	"Kpsk0": Pattern_Kpsk0,
	"Xpsk1": Pattern_Xpsk1,

	// Interactive (fundemental) patterns.
	"NN":     Pattern_NN,
	"NK":     Pattern_NK,
	"NX":     Pattern_NX,
	"XN":     Pattern_XN,
	"XK":     Pattern_XK,
	"XX":     Pattern_XX,
	"KN":     Pattern_KN,
	"KK":     Pattern_KK,
	"KX":     Pattern_KX,
	"IN":     Pattern_IN,
	"IK":     Pattern_IK,
	"IX":     Pattern_IX,
	"NNpsk0": Pattern_NNpsk0,
	"NNpsk2": Pattern_NNpsk2,
	"NKpsk0": Pattern_NKpsk0,
	"NKpsk2": Pattern_NKpsk2,
	"NXpsk2": Pattern_NXpsk2,
	"XNpsk3": Pattern_XNpsk3,
	"XKpsk3": Pattern_XKpsk3,
	"XXpsk3": Pattern_XXpsk3,
	"KNpsk0": Pattern_KNpsk0,
	"KNpsk2": Pattern_KNpsk2,
	"KKpsk0": Pattern_KKpsk0,
	"KKpsk2": Pattern_KKpsk2,
	"KXpsk2": Pattern_KXpsk2,
	"INpsk1": Pattern_INpsk1,
	"INpsk2": Pattern_INpsk2,
	"IKpsk1": Pattern_IKpsk1,
	"IKpsk2": Pattern_IKpsk2,
	"IXpsk2": Pattern_IXpsk2,

	// Deferred patterns.
	"NK1":  Pattern_NK1,
	"NX1":  Pattern_NX1,
	"X1N":  Pattern_X1N,
	"X1K":  Pattern_X1K,
	"XK1":  Pattern_XK1,
	"X1K1": Pattern_X1K1,
	"X1X":  Pattern_X1X,
	"XX1":  Pattern_XX1,
	"X1X1": Pattern_X1X1,
	"K1N":  Pattern_K1N,
	"K1K":  Pattern_K1K,
	"KK1":  Pattern_KK1,
	"K1K1": Pattern_K1K1,
	"K1X":  Pattern_K1X,
	"KX1":  Pattern_KX1,
	"K1X1": Pattern_K1X1,
	"I1N":  Pattern_I1N,
	"I1K":  Pattern_I1K,
	"IK1":  Pattern_IK1,
	"I1K1": Pattern_I1K1,
	"I1X":  Pattern_I1X,
	"IX1":  Pattern_IX1,
	"I1X1": Pattern_I1X1,
}

// Token is a Noise handshake pattern token.
type Token uint8

const (
	Token_e Token = iota
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

// MessagePattern is a sequence of pattern tokens.
type MessagePattern []Token

// HandshakePattern is a handshake pattern.
type HandshakePattern interface {
	fmt.Stringer

	// PreMessages returns the pre-message message patterns.
	PreMessages() []MessagePattern

	// Mesages returns the message patterns.
	Messages() []MessagePattern

	// IsPSK returns true iff the pattern has a `psk` modifier.
	IsPSK() bool

	// IsOneWay returns true iff the pattern is one-way.
	IsOneWay() bool
}

type builtInPattern struct {
	name        string
	preMessages []MessagePattern
	messages    []MessagePattern
	isPSK       bool
	isOneWay    bool
}

func (pa *builtInPattern) String() string {
	return pa.name
}

func (pa *builtInPattern) PreMessages() []MessagePattern {
	return pa.preMessages
}

func (pa *builtInPattern) Messages() []MessagePattern {
	return pa.messages
}

func (pa *builtInPattern) IsPSK() bool {
	return pa.isPSK
}

func (pa *builtInPattern) IsOneWay() bool {
	return pa.isOneWay
}
