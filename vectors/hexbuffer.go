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
