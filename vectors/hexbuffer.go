package vectors

import "encoding/hex"

// HexBuffer is a byte slice that will marshal to/unmarshal from a hex encoded
// string.
type HexBuffer []byte

// MarshalText implements the TextMarshaler interface.
func (x *HexBuffer) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(*x)), nil
}

// UnmarshalText implements the TextUnmarshaler interface.
func (x *HexBuffer) UnmarshalText(data []byte) error {
	b, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	if len(b) == 0 {
		*x = nil
	} else {
		*x = b
	}

	return nil
}
