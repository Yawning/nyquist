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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/yawning/nyquist.git/dh"
	"gitlab.com/yawning/nyquist.git/pattern"
	"gitlab.com/yawning/nyquist.git/vectors"
)

func TestVectors(t *testing.T) {
	// Register the multi-PSK suites covered by the snow test vectors.
	t.Run("RegisterMultiPSK", doRegisterMultiPSK)

	srcImpls := []struct {
		name   string
		skipOk bool
	}{
		{"cacophony", false},
		{"snow", false},
		{"noise-c-basic", true}, // PSK patterns use a non-current name.
	}

	for _, v := range srcImpls {
		t.Run(v.name, func(t *testing.T) {
			doTestVectorsFile(t, v.name, v.skipOk)
		})
	}
}

func doRegisterMultiPSK(t *testing.T) {
	require := require.New(t)

	for _, v := range []struct {
		base     pattern.Pattern
		modifier string
	}{
		{pattern.NN, "psk0+psk2"},
		{pattern.NX, "psk0+psk1+psk2"},
		{pattern.XN, "psk1+psk3"},
		{pattern.XK, "psk0+psk3"},
		{pattern.KN, "psk1+psk2"},
		{pattern.KK, "psk0+psk2"},
		{pattern.IN, "psk1+psk2"},
		{pattern.IK, "psk0+psk2"},
		{pattern.IX, "psk0+psk2"},
		{pattern.XX, "psk0+psk1"},
		{pattern.XX, "psk0+psk2"},
		{pattern.XX, "psk0+psk3"},
		{pattern.XX, "psk0+psk1+psk2+psk3"},
	} {
		pa, err := pattern.MakePSK(v.base, v.modifier)
		require.NoError(err, "MakePSK(%s, %s)", v.base.String(), v.modifier)

		err = pattern.Register(pa)
		require.NoError(err, "Register(%s)", pa)
	}
}

func doTestVectorsFile(t *testing.T, impl string, skipOk bool) {
	require := require.New(t)
	fn := filepath.Join("./testdata/", impl+".txt")
	b, err := ioutil.ReadFile(fn)
	require.NoError(err, "ReadFile(%v)", fn)

	var vectorsFile vectors.File
	err = json.Unmarshal(b, &vectorsFile)
	require.NoError(err, "json.Unmarshal")

	for _, v := range vectorsFile.Vectors {
		if v.Name == "" {
			v.Name = v.ProtocolName
		}
		if v.ProtocolName == "" {
			// The noise-c test vectors have `name` but not `protocol_name`.
			v.ProtocolName = v.Name
		}
		if v.ProtocolName == "" {
			continue
		}
		t.Run(v.Name, func(t *testing.T) {
			doTestVector(t, &v, skipOk)
		})
	}
}

func doTestVector(t *testing.T, v *vectors.Vector, skipOk bool) {
	if v.Fail {
		t.Skip("fail tests not supported")
	}
	if v.Fallback || v.FallbackPattern != "" {
		t.Skip("fallback patterns not supported")
	}

	require := require.New(t)
	initCfg, respCfg := configsFromVector(t, v, skipOk)

	initHs, err := NewHandshake(initCfg)
	require.NoError(err, "NewHandshake(initCfg)")
	defer initHs.Reset()

	respHs, err := NewHandshake(respCfg)
	require.NoError(err, "NewHandshake(respCfg)")
	defer respHs.Reset()

	t.Run("Initiator", func(t *testing.T) {
		doTestVectorMessages(t, initHs, v)
	})
	t.Run("Responder", func(t *testing.T) {
		doTestVectorMessages(t, respHs, v)
	})
}

func doTestVectorMessages(t *testing.T, hs *HandshakeState, v *vectors.Vector) {
	require := require.New(t)

	writeOnEven := hs.isInitiator

	var (
		status     *HandshakeStatus
		txCs, rxCs *CipherState
	)
	for idx, msg := range v.Messages {
		var (
			dst, expectedDst []byte
			err              error
		)

		if status == nil {
			// Handshake message(s).
			if (idx&1 == 0) == writeOnEven {
				dst, err = hs.WriteMessage(nil, msg.Payload)
				expectedDst = msg.Ciphertext
			} else {
				dst, err = hs.ReadMessage(nil, msg.Ciphertext)
				expectedDst = msg.Payload
			}

			switch err {
			case ErrDone:
				status = hs.GetStatus()
				require.Equal(status.Err, ErrDone, "Status.Err indicates normal completion")
				if len(v.HandshakeHash) > 0 {
					// The handshake hash is an optional field in the test vectors,
					// and the ones generated by snow, don't include it.
					require.EqualValues(v.HandshakeHash, status.HandshakeHash, "HandshakeHash matches")
				}
				require.Len(status.CipherStates, 2, "Status has 2 CipherState objects")
				if hs.cfg.Protocol.Pattern.IsOneWay() {
					require.Nil(status.CipherStates[1], "Status CipherStates[1] is nil")
				}

				txCs, rxCs = status.CipherStates[0], status.CipherStates[1]
				if !hs.isInitiator {
					txCs, rxCs = rxCs, txCs
				}
			case nil:
			default:
				require.NoError(err, "Handshake Message - %d", idx)
			}
		} else {
			// The messages that use the derived cipherstates just follow the
			// handshake message(s), and the flow continues.
			if hs.cfg.Protocol.Pattern.IsOneWay() {
				// Except one-way patterns which go from initiator to responder.
				if hs.isInitiator {
					dst, err = txCs.EncryptWithAd(nil, nil, msg.Payload)
					expectedDst = msg.Ciphertext
				} else {
					dst, err = rxCs.DecryptWithAd(nil, nil, msg.Ciphertext)
					expectedDst = msg.Payload
				}
			} else {
				if (idx&1 == 0) == writeOnEven {
					dst, err = txCs.EncryptWithAd(nil, nil, msg.Payload)
					expectedDst = msg.Ciphertext
				} else {
					dst, err = rxCs.DecryptWithAd(nil, nil, msg.Ciphertext)
					expectedDst = msg.Payload
				}
			}
			require.NoError(err, "Transport Message - %d", idx)
		}
		require.EqualValues(expectedDst, dst, "Message - #%d, output matches", idx)
	}

	// Sanity-check the test vectors for stupidity.
	require.NotNil(status, "Status != nil (test vector sanity check)")

	// These usually would be done by defer, but invoke them manually to make
	// sure nothing panics.
	if txCs != nil {
		txCs.Reset()
	}
	if rxCs != nil {
		rxCs.Reset()
	}
}

func configsFromVector(t *testing.T, v *vectors.Vector, skipOk bool) (*HandshakeConfig, *HandshakeConfig) {
	require := require.New(t)

	protoName := v.ProtocolName
	protocol, err := NewProtocol(protoName)
	if err == ErrProtocolNotSupported && skipOk {
		t.Skipf("protocol not supported")
	}
	require.NoError(err, "NewProtocol(%v)", protoName)
	require.Equal(protoName, protocol.String(), "derived protocol name matches test case")
	err = pattern.IsValid(protocol.Pattern)
	require.NoError(err, "IsValid(protocol.Pattern)")

	// Initiator side.
	var initStatic dh.Keypair
	if len(v.InitStatic) != 0 {
		initStatic, err = protocol.DH.ParsePrivateKey(v.InitStatic)
		require.NoError(err, "parse InitStatic")
	}

	var initEphemeral dh.Keypair
	if len(v.InitEphemeral) != 0 {
		initEphemeral, err = protocol.DH.ParsePrivateKey(v.InitEphemeral)
		require.NoError(err, "parse InitEphemeral")
	}

	var initRemoteStatic dh.PublicKey
	if len(v.InitRemoteStatic) != 0 {
		initRemoteStatic, err = protocol.DH.ParsePublicKey(v.InitRemoteStatic)
		require.NoError(err, "parse InitRemoteStatic")
	}

	initCfg := &HandshakeConfig{
		Protocol:       protocol,
		Prologue:       v.InitPrologue,
		LocalStatic:    initStatic,
		LocalEphemeral: initEphemeral,
		RemoteStatic:   initRemoteStatic,
		Rng:            &failReader{},
		IsInitiator:    true,
	}

	require.Len(v.InitPsks, protocol.Pattern.NumPSKs(), "test vector has the expected number of InitPsks")
	for _, psk := range v.InitPsks {
		initCfg.PreSharedKeys = append(initCfg.PreSharedKeys, []byte(psk))
	}

	// Responder side.
	var respStatic dh.Keypair
	if len(v.RespStatic) != 0 {
		respStatic, err = protocol.DH.ParsePrivateKey(v.RespStatic)
		require.NoError(err, "parse RespStatic")
	}

	var respEphemeral dh.Keypair
	if len(v.RespEphemeral) != 0 {
		respEphemeral, err = protocol.DH.ParsePrivateKey(v.RespEphemeral)
		require.NoError(err, "parse RespEphemeral")
	}

	var respRemoteStatic dh.PublicKey
	if len(v.RespRemoteStatic) != 0 {
		respRemoteStatic, err = protocol.DH.ParsePublicKey(v.RespRemoteStatic)
		require.NoError(err, "parse RespRemoteStatic")
	}

	respCfg := &HandshakeConfig{
		Protocol:       protocol,
		Prologue:       v.RespPrologue,
		LocalStatic:    respStatic,
		LocalEphemeral: respEphemeral,
		RemoteStatic:   respRemoteStatic,
		Rng:            &failReader{},
		IsInitiator:    false,
	}

	require.Len(v.RespPsks, protocol.Pattern.NumPSKs(), "test vector has the expected number of RespPsks")
	for _, psk := range v.RespPsks {
		respCfg.PreSharedKeys = append(respCfg.PreSharedKeys, []byte(psk))
	}

	return initCfg, respCfg
}
