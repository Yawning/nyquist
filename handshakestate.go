package nyquist

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// DefaultMaxMessageSize is the default maximum message size.
	DefaultMaxMessageSize = 65535

	// PSKSize is the size of the pre-shared symmetric key.
	PSKSize = 32

	protocolPrefix = "Noise"
)

// Protocol is a the protocol to be used with a handshake.
type Protocol struct {
	Pattern HandshakePattern

	DH     DH
	Cipher Cipher
	Hash   Hash
}

// String returns the string representation of the protocol name.
func (pr *Protocol) String() string {
	parts := []string{
		protocolPrefix,
		pr.Pattern.String(),
		pr.DH.String(),
		pr.Cipher.String(),
		pr.Hash.String(),
	}
	return strings.Join(parts, "_")
}

// NewProtocol returns a Protocol from the provided (case-sensitive) protocol
// name.
//
// Note: Only protocols that can be built with the built-in crypto and patterns
// are supported.  Using custom crypto/patterns will require manually building
// a Protocol object.
func NewProtocol(s string) (*Protocol, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 5 || parts[0] != protocolPrefix {
		return nil, errors.New("nyquist: malformed protocol name")
	}

	var protocol Protocol
	protocol.Pattern = supportedPatterns[parts[1]]
	protocol.DH = supportedDHs[parts[2]]
	protocol.Cipher = supportedCiphers[parts[3]]
	protocol.Hash = supportedHashes[parts[4]]

	if protocol.Pattern == nil || protocol.DH == nil || protocol.Cipher == nil || protocol.Hash == nil {
		return nil, ErrProtocolNotSupported
	}

	return &protocol, nil
}

// HandshakeConfig is a handshake configuration.
//
// Warning: While the config may contain sensitive material like DH private
// keys or a pre-shared key, sanitizing such things are the responsibility of
// the caller, after the handshake completes (or aborts due to an error).
//
// Altering any of the members of this structure while a handshake is in
// progress will result in undefined behavior.
type HandshakeConfig struct {
	// Protocol is the noise protocol to use for this handshake.
	Protocol *Protocol

	// Prologue is the optional pre-handshake prologue input to be included
	// in the handshake hash.
	Prologue []byte

	// LocalStatic is the local static keypair, if any (`s`).
	LocalStatic Keypair

	// LocalEphemeral is the local ephemeral keypair, if any (`e`).
	LocalEphemeral Keypair

	// RemoteStatic is the remote static public key, if any (`rs`).
	RemoteStatic PublicKey

	// RemoteEphemeral is the remote ephemeral public key, if any (`re`).
	RemoteEphemeral PublicKey

	// PreSharedKey is the pre-shared symmetric key for PSK mode handshakes.
	PreSharedKey []byte

	// Observer is the optional handshake observer.
	Observer HandshakeObserver

	// Rng is the entropy source to be used when generating new DH key pairs.
	// If the value is `nil`, `crypto/rand.Reader` will be used.
	Rng io.Reader

	// MaxMessageSize specifies the maximum Noise message size the handshake
	// and session will process or generate.  If the value is `0`,
	// `DefaultMaxMessageSize` will be used.
	//
	// Warning: Values other than the default is a non-standard extension
	// to the protocol.
	MaxMessageSize int

	// IsInitiator should be set to true if this handshake is in the
	// initiator role.
	IsInitiator bool
}

// HandshakeStatus is the status of a handshake.
//
// Warning: It is the caller's responsibility to sanitize the CipherStates
// if desired.  Altering any of the members of this structure while a handshake
// is in progress will result in undefined behavior.
type HandshakeStatus struct {
	// Err is the error representing the status of the handshake.
	//
	// It will be `nil` if the handshake is in progess, `ErrDone` if the
	// handshake is complete, and any other error if the handshake has failed.
	Err error

	// LocalEphemeral is the local ephemeral public key, if any (`e`).
	LocalEphemeral PublicKey

	// RemoteStatic is the remote static public key, if any (`rs`).
	RemoteStatic PublicKey

	// RemoteEphemeral is the remote ephemeral public key, if any (`re`).
	RemoteEphemeral PublicKey

	// CipherStates is the resulting CipherState pair (`(cs1, cs2)`).
	//
	// Note: To prevent misuse, for one-way patterns `cs2` will be nil.
	CipherStates []*CipherState

	// HandshakeHash is the handshake hash (`h`).  This field is only set
	// once the handshake is completed.
	HandshakeHash []byte
}

// HandshakeObserver is a handshake observer for monitoring handshake status.
type HandshakeObserver interface {
	// OnPeerPublicKey will be called when a public key is received from
	// the peer, with the handshake pattern token (`Token_e`, `Token_s`)
	// and public key.
	//
	// Returning a non-nil error will abort the handshake immediately.
	OnPeerPublicKey(Token, PublicKey) error
}

func (cfg *HandshakeConfig) getRng() io.Reader {
	if cfg.Rng == nil {
		return rand.Reader
	}
	return cfg.Rng
}

func (cfg *HandshakeConfig) getMaxMessageSize() int {
	if cfg.MaxMessageSize > 0 {
		return cfg.MaxMessageSize
	}
	if cfg.MaxMessageSize == 0 {
		return DefaultMaxMessageSize
	}
	return 0
}

// HandshakeState is the per-handshake state.
type HandshakeState struct {
	cfg *HandshakeConfig

	dh       DH
	patterns []MessagePattern

	ss *SymmetricState

	s  Keypair
	e  Keypair
	rs PublicKey
	re PublicKey

	status *HandshakeStatus

	patternIndex   int
	maxMessageSize int
	dhLen          int
	isInitiator    bool
	processedE     bool
}

// SymmetricState returns the HandshakeState's encapsulated SymmetricState.
//
// Warning: There should be no reason to call this, ever.
func (hs *HandshakeState) SymmetricState() *SymmetricState {
	return hs.ss
}

// GetStatus returns the HandshakeState's status.
func (hs *HandshakeState) GetStatus() *HandshakeStatus {
	return hs.status
}

// Reset clears the HandshakeState of sensitive data.
//
// Warning: If either of the locak keypairs were provided by the
// HandshakeConfig, they will be left intact.
func (hs *HandshakeState) Reset() {
	if hs.ss != nil {
		hs.ss.Reset()
		hs.ss = nil
	}
	if hs.s != hs.cfg.LocalStatic {
		hs.s.Reset()
	}
	if hs.e != hs.cfg.LocalEphemeral {
		hs.e.Reset()
	}
	// TODO: Should this set hs.status.Err?
}

func (hs *HandshakeState) onWriteTokenE(dst []byte) []byte {
	if hs.processedE {
		hs.status.Err = errors.New("nyquist/HandshakeState/WriteMessage/e: e already set")
		return nil
	}
	hs.processedE = true

	// hs.cfg.LocalEphemeral can be used to pre-generate the ephemeral key,
	// so only generate when required.
	if hs.e == nil {
		if hs.e, hs.status.Err = hs.dh.GenerateKeypair(hs.cfg.getRng()); hs.status.Err != nil {
			return nil
		}
	}
	eBytes := hs.e.Public().Bytes()
	hs.ss.MixHash(eBytes)
	if hs.cfg.Protocol.Pattern.IsPSK() {
		hs.ss.MixKey(eBytes)
	}
	hs.status.LocalEphemeral = hs.e.Public()
	return append(dst, eBytes...)
}

func (hs *HandshakeState) onReadTokenE(payload []byte) []byte {
	if hs.re != nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/ReadMessage/e: re already set")
		return nil
	}
	if len(payload) < hs.dhLen {
		hs.status.Err = errors.New("nyquist/HandshakeState/ReadMessage/e: truncated message")
		return nil
	}
	eBytes, tail := payload[:hs.dhLen], payload[hs.dhLen:]
	if hs.re, hs.status.Err = hs.dh.ParsePublicKey(eBytes); hs.status.Err != nil {
		return nil
	}
	hs.status.RemoteEphemeral = hs.re
	if hs.cfg.Observer != nil {
		if hs.status.Err = hs.cfg.Observer.OnPeerPublicKey(Token_e, hs.re); hs.status.Err != nil {
			return nil
		}
	}
	hs.ss.MixHash(eBytes)
	if hs.cfg.Protocol.Pattern.IsPSK() {
		hs.ss.MixKey(eBytes)
	}
	return tail
}

func (hs *HandshakeState) onWriteTokenS(dst []byte) []byte {
	if hs.s == nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/WriteMessage/s: s not set")
		return nil
	}
	sBytes := hs.s.Public().Bytes()
	return hs.ss.EncryptAndHash(dst, sBytes)
}

func (hs *HandshakeState) onReadTokenS(payload []byte) []byte {
	if hs.rs != nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/ReadMessage/s: rs already set")
		return nil
	}
	tempLen := hs.dhLen
	if hs.ss.cs.HasKey() {
		// The spec says `DHLEN + 16`, but doing it this way allows this
		// implementation to support any AEAD implementation, regardless of
		// tag size.
		tempLen += hs.ss.cs.aead.Overhead()
	}
	if len(payload) < tempLen {
		hs.status.Err = errors.New("nyquist/HandshakeState/ReadMessage/s: truncated message")
		return nil
	}
	temp, tail := payload[:tempLen], payload[tempLen:]

	var sBytes []byte
	if sBytes, hs.status.Err = hs.ss.DecryptAndHash(nil, temp); hs.status.Err != nil {
		return nil
	}
	if hs.rs, hs.status.Err = hs.dh.ParsePublicKey(sBytes); hs.status.Err != nil {
		return nil
	}
	hs.status.RemoteStatic = hs.rs
	if hs.cfg.Observer != nil {
		if hs.status.Err = hs.cfg.Observer.OnPeerPublicKey(Token_s, hs.rs); hs.status.Err != nil {
			return nil
		}
	}
	return tail
}

func (hs *HandshakeState) onTokenEE() {
	if hs.e == nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/Message/ee: e not set")
		return
	}
	if hs.re == nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/Message/ee: re not set")
		return
	}
	var eeBytes []byte
	if eeBytes, hs.status.Err = hs.e.DH(hs.re); hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(eeBytes)
	explicitBzero(eeBytes)
}

func (hs *HandshakeState) onTokenES() {
	var esBytes []byte
	if hs.isInitiator {
		if hs.e == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/es: e not set")
			return
		}
		if hs.rs == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/es: rs not set")
			return
		}
		esBytes, hs.status.Err = hs.e.DH(hs.rs)
	} else {
		if hs.s == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/es: s not set")
			return
		}
		if hs.re == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/es: re not set")
			return
		}
		esBytes, hs.status.Err = hs.s.DH(hs.re)
	}
	if hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(esBytes)
	explicitBzero(esBytes)
}

func (hs *HandshakeState) onTokenSE() {
	var seBytes []byte
	if hs.isInitiator {
		if hs.s == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/se: s not set")
			return
		}
		if hs.re == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/se: re not set")
			return
		}
		seBytes, hs.status.Err = hs.s.DH(hs.re)
	} else {
		if hs.e == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/se: e not set")
			return
		}
		if hs.rs == nil {
			hs.status.Err = errors.New("nyquist/HandshakeState/Message/se: rs not set")
			return
		}
		seBytes, hs.status.Err = hs.e.DH(hs.rs)
	}
	if hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(seBytes)
	explicitBzero(seBytes)
}

func (hs *HandshakeState) onTokenSS() {
	if hs.s == nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/Message/ss: s not set")
		return
	}
	if hs.rs == nil {
		hs.status.Err = errors.New("nyquist/HandshakeState/Message/ss: rs not set")
		return
	}
	var ssBytes []byte
	if ssBytes, hs.status.Err = hs.s.DH(hs.rs); hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(ssBytes)
	explicitBzero(ssBytes)
}

func (hs *HandshakeState) onTokenPsk() {
	// PSK is validated at handshake creation.
	hs.ss.MixKeyAndHash(hs.cfg.PreSharedKey)
}

func (hs *HandshakeState) onDone(dst []byte) ([]byte, error) {
	hs.patternIndex++
	if hs.patternIndex < len(hs.patterns) {
		return dst, nil
	}

	hs.status.Err = ErrDone
	cs1, cs2 := hs.ss.Split()
	if hs.cfg.Protocol.Pattern.IsOneWay() {
		cs2.Reset()
		cs2 = nil
	}
	hs.status.CipherStates = []*CipherState{cs1, cs2}
	hs.status.HandshakeHash = hs.ss.GetHandshakeHash()

	// This will end up being called redundantly if the developer has any
	// sense at al, but it's cheap foot+gun avoidance.
	hs.Reset()

	return dst, hs.status.Err
}

// WriteMessage processes a write step of the handshake protocol, appending the
// handshake protocol message to dst, and returning the potentially new slice.
//
// Iff the handshake is complete, the error returned will be `ErrDone`.
func (hs *HandshakeState) WriteMessage(dst, payload []byte) ([]byte, error) {
	if hs.status.Err != nil {
		return nil, hs.status.Err
	}

	if hs.isInitiator != (hs.patternIndex&1 == 0) {
		hs.status.Err = ErrOutOfOrder
		return nil, hs.status.Err
	}

	baseLen := len(dst)
	for _, v := range hs.patterns[hs.patternIndex] {
		switch v {
		case Token_e:
			dst = hs.onWriteTokenE(dst)
		case Token_s:
			dst = hs.onWriteTokenS(dst)
		case Token_ee:
			hs.onTokenEE()
		case Token_es:
			hs.onTokenES()
		case Token_se:
			hs.onTokenSE()
		case Token_ss:
			hs.onTokenSS()
		case Token_psk:
			hs.onTokenPsk()
		default:
			hs.status.Err = errors.New("nyquist/HandshakeState/WriteMessage: invalid token: " + v.String())
		}

		if hs.status.Err != nil {
			return nil, hs.status.Err
		}
	}

	dst = hs.ss.EncryptAndHash(dst, payload)
	if hs.maxMessageSize > 0 && len(dst)-baseLen > hs.maxMessageSize {
		hs.status.Err = ErrMessageSize
		return nil, hs.status.Err
	}

	return hs.onDone(dst)
}

// ReadMessage processes a read step of the handshake protocol, appended the
// authentiated/decrypted message payload to dst, and returning the potentially
// new slice.
//
// Iff the handshake is complete, the error returned will be `ErrDone`.
func (hs *HandshakeState) ReadMessage(dst, payload []byte) ([]byte, error) {
	if hs.status.Err != nil {
		return nil, hs.status.Err
	}

	if hs.maxMessageSize > 0 && len(payload) > hs.maxMessageSize {
		hs.status.Err = ErrMessageSize
		return nil, hs.status.Err
	}

	if hs.isInitiator != (hs.patternIndex&1 != 0) {
		hs.status.Err = ErrOutOfOrder
		return nil, hs.status.Err
	}

	for _, v := range hs.patterns[hs.patternIndex] {
		switch v {
		case Token_e:
			payload = hs.onReadTokenE(payload)
		case Token_s:
			payload = hs.onReadTokenS(payload)
		case Token_ee:
			hs.onTokenEE()
		case Token_es:
			hs.onTokenES()
		case Token_se:
			hs.onTokenSE()
		case Token_ss:
			hs.onTokenSS()
		case Token_psk:
			hs.onTokenPsk()
		default:
			hs.status.Err = errors.New("nyquist/HandshakeState/ReadMessage: invalid token: " + v.String())
		}

		if hs.status.Err != nil {
			return nil, hs.status.Err
		}
	}

	dst, hs.status.Err = hs.ss.DecryptAndHash(dst, payload)
	if hs.status.Err != nil {
		return nil, hs.status.Err
	}

	return hs.onDone(dst)
}

func (hs *HandshakeState) handlePreMessages() error {
	preMessages := hs.cfg.Protocol.Pattern.PreMessages()
	if len(preMessages) == 0 {
		return nil
	}

	// Do everything from the point of view of the initiator to simplify
	// processing.
	var s, e, rs, re PublicKey
	rs, re = hs.rs, hs.re
	if hs.s != nil {
		s = hs.s.Public()
	}
	if hs.e != nil {
		e = hs.e.Public()
	}
	if !hs.isInitiator {
		s, e, rs, re = rs, re, s, e
	}

	if err := hs.handlePreMessage(preMessages[0], s, e, "initiator"); err != nil {
		return err
	}

	if len(preMessages) == 1 {
		return nil
	}

	if err := hs.handlePreMessage(preMessages[1], rs, re, "responder"); err != nil {
		return err
	}

	return nil
}

func (hs *HandshakeState) handlePreMessage(preMessage MessagePattern, s, e PublicKey, side string) error {
	for _, v := range preMessage {
		switch v {
		case Token_e:
			if e == nil {
				return fmt.Errorf("nyquist/New: %s e not set", side)
			}
			pkBytes := e.Bytes()
			hs.ss.MixHash(pkBytes)
			if hs.cfg.Protocol.Pattern.IsPSK() {
				hs.ss.MixKey(pkBytes)
			}
		case Token_s:
			if s == nil {
				return fmt.Errorf("nyquist/New: %s s not set", side)
			}
			hs.ss.MixHash(s.Bytes())
		default:
		}
	}
	return nil
}

// NewHandshake constructs a new HandshakeState with the provided configuration.
// This call is equivalent to the `Initialize` HandshakeState call in the
// Noise Protocol Framework specification.
func NewHandshake(cfg *HandshakeConfig) (*HandshakeState, error) {
	// TODO: Validate the config further?

	if cfg.Protocol.Pattern.IsPSK() {
		if len(cfg.PreSharedKey) != PSKSize {
			return nil, errors.New("nyquist/New: invalid or missing PreSharedKey")
		}
	}

	maxMessageSize := cfg.getMaxMessageSize()
	hs := &HandshakeState{
		cfg:      cfg,
		dh:       cfg.Protocol.DH,
		patterns: cfg.Protocol.Pattern.Messages(),
		ss:       newSymmetricState(cfg.Protocol.Cipher, cfg.Protocol.Hash, maxMessageSize),
		s:        cfg.LocalStatic,
		e:        cfg.LocalEphemeral,
		rs:       cfg.RemoteStatic,
		re:       cfg.RemoteEphemeral,
		status: &HandshakeStatus{
			RemoteStatic:    cfg.RemoteStatic,
			RemoteEphemeral: cfg.RemoteEphemeral,
		},
		maxMessageSize: maxMessageSize,
		dhLen:          cfg.Protocol.DH.Size(),
		isInitiator:    cfg.IsInitiator,
	}
	if cfg.LocalEphemeral != nil {
		hs.status.LocalEphemeral = cfg.LocalEphemeral.Public()
	}

	hs.ss.InitializeSymmetric([]byte(cfg.Protocol.String()))
	hs.ss.MixHash(cfg.Prologue)
	if err := hs.handlePreMessages(); err != nil {
		return nil, err
	}

	return hs, nil
}
