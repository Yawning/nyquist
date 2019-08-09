// Package vectors provides types for the JSON formatted test vectors.
package vectors

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
