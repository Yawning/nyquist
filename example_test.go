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
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/yawning/nyquist.git/cipher"
	"gitlab.com/yawning/nyquist.git/dh"
	"gitlab.com/yawning/nyquist.git/hash"
	"gitlab.com/yawning/nyquist.git/pattern"
)

func TestExample(t *testing.T) {
	require := require.New(t)

	// Protocols can be constructed by parsing a protocol name.
	protocol, err := NewProtocol("Noise_XX_25519_ChaChaPoly_BLAKE2s")
	require.NoError(err, "NewProtocol")

	// Protocols can also be constructed manually.
	protocol2 := &Protocol{
		Pattern: pattern.XX,
		DH:      dh.X25519,
		Cipher:  cipher.ChaChaPoly,
		Hash:    hash.BLAKE2s,
	}
	require.Equal(protocol, protocol2)

	// Each side needs a HandshakeConfig, properly filled out.
	aliceStatic, err := protocol.DH.GenerateKeypair(rand.Reader)
	require.NoError(err, "Generate Alice's static keypair")
	aliceCfg := &HandshakeConfig{
		Protocol:    protocol,
		LocalStatic: aliceStatic,
		IsInitiator: true,
	}

	bobStatic, err := protocol.DH.GenerateKeypair(rand.Reader)
	require.NoError(err, "Generate Bob's static keypair")
	bobCfg := &HandshakeConfig{
		Protocol:    protocol,
		LocalStatic: bobStatic,
		IsInitiator: false,
	}

	// Each side then constructs a HandshakeState.
	aliceHs, err := NewHandshake(aliceCfg)
	require.NoError(err, "NewHandshake(aliceCfg)")

	bobHs, err := NewHandshake(bobCfg)
	require.NoError(err, "NewHandshake(bobCfg")

	// Ensuring that HandshakeState.Reset() is called, will make sure that
	// the HandshakeState isn't inadvertently reused.
	defer aliceHs.Reset()
	defer bobHs.Reset()

	// The SymmetricState and CipherState objects embedded in the
	// HandshakeState can be accessed while the handshake is in progress,
	// though most users likely will not need to do this.
	aliceSs := aliceHs.SymmetricState()
	require.NotNil(aliceSs, "aliceHs.SymmetricState()")
	aliceCs := aliceSs.CipherState()
	require.NotNil(aliceCs, "aliceSS.CipherState()")

	// Then, each side calls hs.ReadMessage/hs.WriteMessage as appropriate.
	alicePlaintextE := []byte("alice e plaintext") // Handshake message payloads are optional.
	aliceMsg1, err := aliceHs.WriteMessage(nil, alicePlaintextE)
	require.NoError(err, "aliceHs.WriteMessage(1)") // (alice) -> e

	bobRecv, err := bobHs.ReadMessage(nil, aliceMsg1)
	require.NoError(err, "bobHs.ReadMessage(alice1)")
	require.Equal(bobRecv, alicePlaintextE)

	bobMsg1, err := bobHs.WriteMessage(nil, nil) // (bob) -> e, ee, s, es
	require.NoError(err, "bobHS.WriteMessage(bob1)")

	// When the handshake is complete, ErrDone is returned as the error.
	_, err = aliceHs.ReadMessage(nil, bobMsg1)
	require.NoError(err, "aliceHS.ReadMessage(bob1)")

	aliceMsg2, err := aliceHs.WriteMessage(nil, nil) // (alice) -> s, se
	require.Equal(ErrDone, err, "aliceHs.WriteMessage(alice2)")

	_, err = bobHs.ReadMessage(nil, aliceMsg2)
	require.Equal(ErrDone, err, "bobHs.ReadMessage(alice2)")

	// Once a handshake is completed, the CipherState objects, handshake hash
	// and various public keys can be pulled out of the HandshakeStatus object.
	aliceStatus := aliceHs.GetStatus()
	bobStatus := bobHs.GetStatus()

	require.Equal(aliceStatus.HandshakeHash, bobStatus.HandshakeHash, "Handshake hashes match")
	require.Equal(aliceStatus.LocalEphemeral.Bytes(), bobStatus.RemoteEphemeral.Bytes())
	require.Equal(bobStatus.LocalEphemeral.Bytes(), aliceStatus.RemoteEphemeral.Bytes())
	require.Equal(aliceStatus.RemoteStatic.Bytes(), bobStatic.Public().Bytes())
	require.Equal(bobStatus.RemoteStatic.Bytes(), aliceStatic.Public().Bytes())

	// Then the CipherState objects can be used to exchange messages.
	aliceTx, aliceRx := aliceStatus.CipherStates[0], aliceStatus.CipherStates[1]
	bobRx, bobTx := bobStatus.CipherStates[0], bobStatus.CipherStates[1] // Reversed from alice!

	// Naturally CipherState.Reset() also exists.
	defer func() {
		aliceTx.Reset()
		aliceRx.Reset()
	}()
	defer func() {
		bobTx.Reset()
		bobRx.Reset()
	}()

	// Alice -> Bob, post-handshake.
	alicePlaintext := []byte("alice transport plaintext")
	aliceMsg3, err := aliceTx.EncryptWithAd(nil, nil, alicePlaintext)
	require.NoError(err, "aliceTx.EncryptWithAd()")

	bobRecv, err = bobRx.DecryptWithAd(nil, nil, aliceMsg3)
	require.NoError(err, "bobRx.DecryptWithAd()")
	require.Equal(alicePlaintext, bobRecv)

	// Bob -> Alice, post-handshake.
	bobPlaintext := []byte("bob transport plaintext")
	bobMsg2, err := bobTx.EncryptWithAd(nil, nil, bobPlaintext)
	require.NoError(err, "bobTx.EncryptWithAd()")

	aliceRecv, err := aliceRx.DecryptWithAd(nil, nil, bobMsg2)
	require.NoError(err, "aliceRx.DecryptWithAd")
	require.Equal(bobPlaintext, aliceRecv)
}
