// Copyright (C) 2019, 2021 Yawning Angel. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package nyquist

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/yawning/nyquist.git/dh"
	"gitlab.com/yawning/nyquist.git/pattern"
)

const xFixedSize = 32 + 32 + 16 + 16

var errFailReader = errors.New("nyquist/test: failReader.Read")

type failReader struct{}

func (r *failReader) Read(p []byte) (int, error) {
	return 0, errFailReader
}

func mustMakeX(t *testing.T, mms int) (*HandshakeState, *HandshakeState) {
	require := require.New(t)

	protocol, err := NewProtocol("Noise_X_25519_ChaChaPoly_BLAKE2s")
	require.NotNil(protocol, "NewProtocol")
	require.NoError(err, "NewProtocol")

	aliceStatic, err := protocol.DH.GenerateKeypair(rand.Reader)
	require.NoError(err, "Generate Alice's static keypair")

	bobStatic, err := protocol.DH.GenerateKeypair(rand.Reader)
	require.NoError(err, "Generate Bob's static keypair")

	aliceCfg := &HandshakeConfig{
		Protocol:       protocol,
		LocalStatic:    aliceStatic,
		RemoteStatic:   bobStatic.Public(),
		MaxMessageSize: mms,
		IsInitiator:    true,
	}
	aliceHs, err := NewHandshake(aliceCfg)
	require.NoError(err, "NewHandshake(aliceCfg)")

	bobCfg := &HandshakeConfig{
		Protocol:       protocol,
		LocalStatic:    bobStatic,
		MaxMessageSize: mms,
	}
	bobHs, err := NewHandshake(bobCfg)
	require.NoError(err, "NewHandshake(bobCfg)")

	return aliceHs, bobHs
}

func TestHandshakeState(t *testing.T) {
	for _, v := range []struct {
		n  string
		fn func(*testing.T)
	}{
		{"BadProtocol", testHandshakeStateBadProtocol},
		{"KeygenFailure", testHandshakeStateKeygenFailure},
		{"TruncatedE", testHandshakeStateTruncatedE},
		{"TruncatedS", testHandshakeStateTruncatedS},
		{"OutOfOrder", testHandshakeStateOutOfOrder},
		{"MaxMessageSize", testHandshakeStateMaxMessageSize},
		{"Observer", testHandshakeStateObserver},
		{"BadPSK", testHandshakeStateBadPSK},
		{"MissingS", testHandshakeStateMissingS},
	} {
		t.Run(v.n, v.fn)
	}
}

func testHandshakeStateBadProtocol(t *testing.T) {
	require := require.New(t)

	var nilProtocol Protocol
	s := nilProtocol.String()
	require.Equal(invalidProtocol, s, "Bad protocol object String()s")

	protocol, err := NewProtocol("Signal_XX_25519_ChaChaPoly_BLAKE2s")
	require.Nil(protocol, "NewProtocol(invalid)")
	require.Equal(ErrProtocolNotSupported, err, "NewProtocol(invalid)")
}

func testHandshakeStateKeygenFailure(t *testing.T) {
	require := require.New(t)

	protocol, err := NewProtocol("Noise_N_25519_ChaChaPoly_BLAKE2s")
	require.NotNil(protocol, "NewProtocol")
	require.NoError(err, "NewProtocol")

	bobStatic, err := protocol.DH.GenerateKeypair(rand.Reader)
	require.NoError(err, "Generate Bob's static keypair")

	aliceCfg := &HandshakeConfig{
		Protocol:     protocol,
		RemoteStatic: bobStatic.Public(),
		Rng:          &failReader{},
		IsInitiator:  true,
	}
	aliceHs, err := NewHandshake(aliceCfg)
	require.NoError(err, "NewHandshake(aliceCfg)")

	dst, err := aliceHs.WriteMessage(nil, nil)
	require.Nil(dst, "aliceHs.WriteMessage - e generation will fail")
	require.Equal(errFailReader, err, "aliceHs.WriteMessage - e generation will fail")
}

func testHandshakeStateTruncatedE(t *testing.T) {
	require := require.New(t)

	_, bobHs := mustMakeX(t, 0)
	dst, err := bobHs.ReadMessage(nil, make([]byte, 31))
	require.Nil(dst, "bobHs.ReadMessage - truncated E")
	require.Equal(errTruncatedE, err)
}

func testHandshakeStateTruncatedS(t *testing.T) {
	require := require.New(t)

	aliceHs, bobHs := mustMakeX(t, 0)
	dst, err := aliceHs.WriteMessage(nil, nil)
	require.Equal(ErrDone, err, "aliceHs.WriteMessage")
	require.Len(dst, xFixedSize, "aliceHs.WriteMessage") // e, es, s, ss

	dst, err = bobHs.ReadMessage(nil, dst[:32+32]) // Clip off both tags.
	require.Nil(dst, "bobHs.ReadMessage - truncated s")
	require.Equal(errTruncatedS, err)
}

func testHandshakeStateOutOfOrder(t *testing.T) {
	require := require.New(t)

	aliceHs, bobHs := mustMakeX(t, 0)

	dst, err := aliceHs.ReadMessage(nil, []byte("never read, whatever"))
	require.Nil(dst, "aliceHs.ReadMessage - out of order")
	require.Equal(ErrOutOfOrder, err, "aliceHs.ReadMessage - out of order")

	dst, err = bobHs.WriteMessage(nil, []byte("placeholder plaintext pls ignore"))
	require.Nil(dst, "bobHs.WriteMessage - after critical failure")
	require.Equal(ErrOutOfOrder, err, "bobHs.WriteMessage - after critical failure")

	// While we are here and have two busted HandshakeState objects, make
	// sure that the errors are sticky.
	dst, err = aliceHs.WriteMessage(nil, []byte("placeholder plaintext pls ignore"))
	require.Nil(dst, "aliceHs.WriteMessage - after critical failure")
	require.Equal(ErrOutOfOrder, err, "aliceHs.WriteMessage - after critical failure")
	require.Equal(err, aliceHs.GetStatus().Err)

	dst, err = bobHs.ReadMessage(nil, []byte("never read, whatever"))
	require.Nil(dst, "bobHs.ReadMessage - after critical failure")
	require.Equal(ErrOutOfOrder, err, "bobHs.WriteMessage - after critical failure")
	require.Equal(err, bobHs.GetStatus().Err)
}

func testHandshakeStateMaxMessageSize(t *testing.T) {
	const testMMS = 127

	require := require.New(t)

	// Ensure that exactly the maximum message size passes.
	aliceHs, bobHs := mustMakeX(t, testMMS)
	maxSizedPayload := make([]byte, testMMS-xFixedSize)
	_, _ = rand.Read(maxSizedPayload)
	dst, err := aliceHs.WriteMessage(nil, maxSizedPayload)
	require.Equal(ErrDone, err, "aliceHs.WriteMessage(maxSize)")
	require.Len(dst, testMMS, "aliceHs.WriteMessage(maxSize)")

	dst, err = bobHs.ReadMessage(nil, dst)
	require.Equal(ErrDone, err, "bobHs.ReadMessage(maxSize)")
	require.Equal(maxSizedPayload, dst, "bobHs.ReadMessage(maxSize)")

	// Ensure that the payloads at 1 past the limit fail.
	aliceHs, bobHs = mustMakeX(t, testMMS)
	oversizedPayload := append(maxSizedPayload, 23)
	dst, err = aliceHs.WriteMessage(nil, oversizedPayload)
	require.Equal(ErrMessageSize, err, "aliceHs.WriteMessage(overSize)")
	require.Nil(dst, "aliceHs.WriteMessage(overSize)")

	dst, err = bobHs.ReadMessage(nil, make([]byte, testMMS+1))
	require.Equal(ErrMessageSize, err, "bobHs.ReadMessage(overSize)")
	require.Nil(dst, "bobHs.ReadMessage(overSize)")

	// Ensure that a negative mms disables limit enforcement.
	aliceHs, bobHs = mustMakeX(t, -1)
	giantPayload := make([]byte, DefaultMaxMessageSize*10)
	_, _ = rand.Read(giantPayload)
	dst, err = aliceHs.WriteMessage(nil, giantPayload)
	require.Equal(ErrDone, err, "aliceHs.WriteMessage(giantSize)")
	require.Len(dst, len(giantPayload)+xFixedSize, "aliceHs.WriteMessage(giantSize)")

	dst, err = bobHs.ReadMessage(nil, dst)
	require.Equal(ErrDone, err, "bobHs.ReadMessage(giantSize)")
	require.Equal(giantPayload, dst, "bobHs.ReadMessage(giantSize)")
}

type proxyObserver struct {
	callbackFn func(pattern.Token, dh.PublicKey) error
}

func (proxy *proxyObserver) OnPeerPublicKey(token pattern.Token, pk dh.PublicKey) error {
	return proxy.callbackFn(token, pk)
}

func testHandshakeStateObserver(t *testing.T) {
	require := require.New(t)

	aliceHs, bobHs := mustMakeX(t, 0)

	var seenE, seenS bool
	proxy := &proxyObserver{
		callbackFn: func(token pattern.Token, pk dh.PublicKey) error {
			switch token {
			case pattern.Token_e:
				require.False(seenE)
				require.Equal(pk.Bytes(), aliceHs.GetStatus().LocalEphemeral.Bytes())
				seenE = true
			case pattern.Token_s:
				require.False(seenS)
				require.Equal(pk.Bytes(), aliceHs.cfg.LocalStatic.Public().Bytes())
				seenS = true
			default:
				panic("unknown token: " + token.String())
			}
			return nil
		},
	}
	bobHs.cfg.Observer = proxy // Yeah this is ugly, but it works.

	dst, err := aliceHs.WriteMessage(nil, nil)
	require.Equal(ErrDone, err, "aliceHs.WriteMessage()")
	require.Len(dst, xFixedSize, "aliceHs.WriteMessage()")

	dst, err = bobHs.ReadMessage(nil, dst)
	require.Equal(ErrDone, err, "bobHs.ReadMessage()")
	require.Len(dst, 0, "bobHs.ReadMessage()")
	require.True(seenE, "bobHs observer saw alice e")
	require.True(seenS, "bobHs observer saw alice s")
}

func testHandshakeStateBadPSK(t *testing.T) {
	require := require.New(t)

	protocol, err := NewProtocol("Noise_Xpsk1_25519_ChaChaPoly_BLAKE2s")
	require.NotNil(protocol, "NewProtocol")
	require.NoError(err, "NewProtocol")

	aliceCfg := &HandshakeConfig{
		Protocol:    protocol,
		IsInitiator: true,
	}
	_, err = NewHandshake(aliceCfg)
	require.Equal(errMissingPSK, err, "NewHandshake() - missing PSK")

	aliceCfg.PreSharedKeys = [][]byte{
		make([]byte, PreSharedKeySize+1),
	}
	_, err = NewHandshake(aliceCfg)
	require.Equal(errBadPSK, err, "NewHandshake() - malformed PSK")
}

func testHandshakeStateMissingS(t *testing.T) {
	require := require.New(t)

	protocol, err := NewProtocol("Noise_X_25519_ChaChaPoly_BLAKE2s")
	require.NotNil(protocol, "NewProtocol")
	require.NoError(err, "NewProtocol")

	aliceCfg := &HandshakeConfig{
		Protocol:    protocol,
		IsInitiator: true,
	}
	_, err = NewHandshake(aliceCfg)
	require.EqualError(err, "nyquist/New: responder s not set", "NewHandshake() - missing s")

	aliceHs, _ := mustMakeX(t, 0)
	aliceHs.s = nil // Not the best way to do this, but this also works.
	dst, err := aliceHs.WriteMessage(nil, nil)
	require.Equal(errMissingS, err, "aliceHs.WriteMessage()")
	require.Nil(dst, "aliceHs.WriteMessage()")
}
