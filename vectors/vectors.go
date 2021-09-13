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

// Package vectors provides types for the JSON formatted test vectors.
package vectors // import "gitlab.com/yawning/nyquist.git/vectors"

// Message is a test vector handshake message.
type Message struct {
	Payload    HexBuffer `json:"payload"`
	Ciphertext HexBuffer `json:"ciphertext"`
}

// Vector is a single test vector case.
type Vector struct {
	Name string `json:"name"`

	ProtocolName    string `json:"protocol_name"`
	Fail            bool   `json:"fail"`
	Fallback        bool   `json:"fallback"`
	FallbackPattern string `json:"fallback_pattern"`

	InitPrologue     HexBuffer   `json:"init_prologue"`
	InitPsks         []HexBuffer `json:"init_psks"`
	InitStatic       HexBuffer   `json:"init_static"`
	InitEphemeral    HexBuffer   `json:"init_ephemeral"`
	InitRemoteStatic HexBuffer   `json:"init_remote_static"`

	RespPrologue     HexBuffer   `json:"resp_prologue"`
	RespPsks         []HexBuffer `json:"resp_psks"`
	RespStatic       HexBuffer   `json:"resp_static"`
	RespEphemeral    HexBuffer   `json:"resp_ephemeral"`
	RespRemoteStatic HexBuffer   `json:"resp_remote_static"`

	HandshakeHash HexBuffer `json:"handshake_hash"`

	Messages []Message `json:"messages"`
}

// File is a collection of test vectors.
type File struct {
	Vectors []Vector `json:"vectors"`
}
