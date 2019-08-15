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
