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
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"

	"gitlab.com/yawning/nyquist.git/cipher"
	"gitlab.com/yawning/nyquist.git/dh"
	"gitlab.com/yawning/nyquist.git/hash"
	"gitlab.com/yawning/nyquist.git/internal"
	"gitlab.com/yawning/nyquist.git/pattern"
)

const (
	// DefaultMaxMessageSize is the default maximum message size.
	DefaultMaxMessageSize = 65535

	// PreSharedKeySize is the size of the pre-shared symmetric key.
	PreSharedKeySize = 32

	protocolPrefix  = "Noise"
	invalidProtocol = "[invalid protocol]"
)

var (
	errTruncatedE = errors.New("nyquist/HandshakeState/ReadMessage/e: truncated message")
	errTruncatedS = errors.New("nyquist/HandshakeState/ReadMessage/s: truncated message")
	errMissingS   = errors.New("nyquist/HandshakeState/WriteMessage/s: s not set")

	errMissingPSK = errors.New("nyquist/New: missing or excessive PreSharedKey(s)")
	errBadPSK     = errors.New("nyquist/New: malformed PreSharedKey(s)")
)

// Protocol is a the protocol to be used with a handshake.
type Protocol struct {
	Pattern pattern.Pattern

	DH     dh.DH
	Cipher cipher.Cipher
	Hash   hash.Hash
}

// String returns the string representation of the protocol name.
func (pr *Protocol) String() string {
	if pr.Pattern == nil || pr.DH == nil || pr.Cipher == nil || pr.Hash == nil {
		return invalidProtocol
	}

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
// name.  Returned protocol objects may be reused across multiple
// HandshakeConfigs.
//
// Note: Only protocols that can be built with the built-in crypto and patterns
// are supported.  Using custom crypto/patterns will require manually building
// a Protocol object.
func NewProtocol(s string) (*Protocol, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 5 || parts[0] != protocolPrefix {
		return nil, ErrProtocolNotSupported
	}

	var pr Protocol
	pr.Pattern = pattern.FromString(parts[1])
	pr.DH = dh.FromString(parts[2])
	pr.Cipher = cipher.FromString(parts[3])
	pr.Hash = hash.FromString(parts[4])

	if pr.Pattern == nil || pr.DH == nil || pr.Cipher == nil || pr.Hash == nil {
		return nil, ErrProtocolNotSupported
	}

	return &pr, nil
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
	LocalStatic dh.Keypair

	// LocalEphemeral is the local ephemeral keypair, if any (`e`).
	LocalEphemeral dh.Keypair

	// RemoteStatic is the remote static public key, if any (`rs`).
	RemoteStatic dh.PublicKey

	// RemoteEphemeral is the remote ephemeral public key, if any (`re`).
	RemoteEphemeral dh.PublicKey

	// PreSharedKeys is the vector of pre-shared symmetric key for PSK mode
	// handshakes.
	PreSharedKeys [][]byte

	// Observer is the optional handshake observer.
	Observer HandshakeObserver

	// Rng is the entropy source to be used when generating new DH key pairs.
	// If the value is `nil`, `crypto/rand.Reader` will be used.
	Rng io.Reader

	// MaxMessageSize specifies the maximum Noise message size the handshake
	// and session will process or generate.  If the value is `0`,
	// `DefaultMaxMessageSize` will be used.  A negative value will disable
	// the maximum message size enforcement entirely.
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
	LocalEphemeral dh.PublicKey

	// RemoteStatic is the remote static public key, if any (`rs`).
	RemoteStatic dh.PublicKey

	// RemoteEphemeral is the remote ephemeral public key, if any (`re`).
	RemoteEphemeral dh.PublicKey

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
	// the peer, with the handshake pattern token (`pattern.Token_e`,
	// `pattern.Token_s`) and public key.
	//
	// Returning a non-nil error will abort the handshake immediately.
	OnPeerPublicKey(pattern.Token, dh.PublicKey) error
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

	dh       dh.DH
	patterns []pattern.Message

	ss *SymmetricState

	s  dh.Keypair
	e  dh.Keypair
	rs dh.PublicKey
	re dh.PublicKey

	status *HandshakeStatus

	patternIndex   int
	pskIndex       int
	maxMessageSize int
	dhLen          int
	isInitiator    bool
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
// Warning: If either of the local keypairs were provided by the
// HandshakeConfig, they will be left intact.
func (hs *HandshakeState) Reset() {
	if hs.ss != nil {
		hs.ss.Reset()
		hs.ss = nil
	}
	if hs.s != nil && hs.s != hs.cfg.LocalStatic {
		// Having a local static key, that isn't from the config currently can't
		// happen, but having the sanitization is harmless.
		hs.s.Reset()
	}
	if hs.e != nil && hs.e != hs.cfg.LocalEphemeral {
		hs.e.Reset()
	}
	// TODO: Should this set hs.status.Err?
}

func (hs *HandshakeState) onWriteTokenE(dst []byte) []byte {
	// hs.cfg.LocalEphemeral can be used to pre-generate the ephemeral key,
	// so only generate when required.
	if hs.e == nil {
		if hs.e, hs.status.Err = hs.dh.GenerateKeypair(hs.cfg.getRng()); hs.status.Err != nil {
			return nil
		}
	}
	eBytes := hs.e.Public().Bytes()
	hs.ss.MixHash(eBytes)
	if hs.cfg.Protocol.Pattern.NumPSKs() > 0 {
		hs.ss.MixKey(eBytes)
	}
	hs.status.LocalEphemeral = hs.e.Public()
	return append(dst, eBytes...)
}

func (hs *HandshakeState) onReadTokenE(payload []byte) []byte {
	if len(payload) < hs.dhLen {
		hs.status.Err = errTruncatedE
		return nil
	}
	eBytes, tail := payload[:hs.dhLen], payload[hs.dhLen:]
	if hs.re, hs.status.Err = hs.dh.ParsePublicKey(eBytes); hs.status.Err != nil {
		return nil
	}
	hs.status.RemoteEphemeral = hs.re
	if hs.cfg.Observer != nil {
		if hs.status.Err = hs.cfg.Observer.OnPeerPublicKey(pattern.Token_e, hs.re); hs.status.Err != nil {
			return nil
		}
	}
	hs.ss.MixHash(eBytes)
	if hs.cfg.Protocol.Pattern.NumPSKs() > 0 {
		hs.ss.MixKey(eBytes)
	}
	return tail
}

func (hs *HandshakeState) onWriteTokenS(dst []byte) []byte {
	if hs.s == nil {
		hs.status.Err = errMissingS
		return nil
	}
	sBytes := hs.s.Public().Bytes()
	return hs.ss.EncryptAndHash(dst, sBytes)
}

func (hs *HandshakeState) onReadTokenS(payload []byte) []byte {
	tempLen := hs.dhLen
	if hs.ss.cs.HasKey() {
		// The spec says `DHLEN + 16`, but doing it this way allows this
		// implementation to support any AEAD implementation, regardless of
		// tag size.
		tempLen += hs.ss.cs.aead.Overhead()
	}
	if len(payload) < tempLen {
		hs.status.Err = errTruncatedS
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
		if hs.status.Err = hs.cfg.Observer.OnPeerPublicKey(pattern.Token_s, hs.rs); hs.status.Err != nil {
			return nil
		}
	}
	return tail
}

func (hs *HandshakeState) onTokenEE() {
	var eeBytes []byte
	if eeBytes, hs.status.Err = hs.e.DH(hs.re); hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(eeBytes)
	internal.ExplicitBzero(eeBytes)
}

func (hs *HandshakeState) onTokenES() {
	var esBytes []byte
	if hs.isInitiator {
		esBytes, hs.status.Err = hs.e.DH(hs.rs)
	} else {
		esBytes, hs.status.Err = hs.s.DH(hs.re)
	}
	if hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(esBytes)
	internal.ExplicitBzero(esBytes)
}

func (hs *HandshakeState) onTokenSE() {
	var seBytes []byte
	if hs.isInitiator {
		seBytes, hs.status.Err = hs.s.DH(hs.re)
	} else {
		seBytes, hs.status.Err = hs.e.DH(hs.rs)
	}
	if hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(seBytes)
	internal.ExplicitBzero(seBytes)
}

func (hs *HandshakeState) onTokenSS() {
	var ssBytes []byte
	if ssBytes, hs.status.Err = hs.s.DH(hs.rs); hs.status.Err != nil {
		return
	}
	hs.ss.MixKey(ssBytes)
	internal.ExplicitBzero(ssBytes)
}

func (hs *HandshakeState) onTokenPsk() {
	// PSK is validated at handshake creation.
	hs.ss.MixKeyAndHash(hs.cfg.PreSharedKeys[hs.pskIndex])
	hs.pskIndex++
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
		case pattern.Token_e:
			dst = hs.onWriteTokenE(dst)
		case pattern.Token_s:
			dst = hs.onWriteTokenS(dst)
		case pattern.Token_ee:
			hs.onTokenEE()
		case pattern.Token_es:
			hs.onTokenES()
		case pattern.Token_se:
			hs.onTokenSE()
		case pattern.Token_ss:
			hs.onTokenSS()
		case pattern.Token_psk:
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

// ReadMessage processes a read step of the handshake protocol, appending the
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
		case pattern.Token_e:
			payload = hs.onReadTokenE(payload)
		case pattern.Token_s:
			payload = hs.onReadTokenS(payload)
		case pattern.Token_ee:
			hs.onTokenEE()
		case pattern.Token_es:
			hs.onTokenES()
		case pattern.Token_se:
			hs.onTokenSE()
		case pattern.Token_ss:
			hs.onTokenSS()
		case pattern.Token_psk:
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

	// Gather all the public keys from the config, from the initiator's
	// point of view.
	var s, e, rs, re dh.PublicKey
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

	for i, keys := range []struct {
		s, e dh.PublicKey
		side string
	}{
		{s, e, "initiator"},
		{rs, re, "responder"},
	} {
		if i+1 > len(preMessages) {
			break
		}

		for _, v := range preMessages[i] {
			switch v {
			case pattern.Token_e:
				// While the specification allows for `e` tokens in the
				// pre-messages, there are currently no patterns that use
				// such a construct.
				//
				// While it is possible to generate `e` if it is the local
				// one that is missing, that would be stretching a use-case
				// that is already somewhat nonsensical.
				if keys.e == nil {
					return fmt.Errorf("nyquist/New: %s e not set", keys.side)
				}
				pkBytes := keys.e.Bytes()
				hs.ss.MixHash(pkBytes)
				if hs.cfg.Protocol.Pattern.NumPSKs() > 0 {
					hs.ss.MixKey(pkBytes)
				}
			case pattern.Token_s:
				if keys.s == nil {
					return fmt.Errorf("nyquist/New: %s s not set", keys.side)
				}
				hs.ss.MixHash(keys.s.Bytes())
			default:
				return errors.New("nyquist/New: invalid pre-message token: " + v.String())
			}
		}
	}

	return nil
}

// NewHandshake constructs a new HandshakeState with the provided configuration.
// This call is equivalent to the `Initialize` HandshakeState call in the
// Noise Protocol Framework specification.
func NewHandshake(cfg *HandshakeConfig) (*HandshakeState, error) {
	// TODO: Validate the config further?

	if cfg.Protocol.Pattern.NumPSKs() != len(cfg.PreSharedKeys) {
		return nil, errMissingPSK
	}
	for _, v := range cfg.PreSharedKeys {
		if len(v) != PreSharedKeySize {
			return nil, errBadPSK
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
