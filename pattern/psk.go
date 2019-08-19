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

package pattern

import (
	"errors"
	"strconv"
	"strings"
)

const prefixPSK = "psk"

// MakePSK applies `psk` modifiers to an existing pattern, returning the new
// pattern.
func MakePSK(template Pattern, modifier string) (Pattern, error) {
	if template.NumPSKs() > 0 {
		return nil, errors.New("nyquist/pattern: PSK template pattern already has PSKs")
	}

	pa := &builtIn{
		name:        template.String() + modifier,
		preMessages: template.PreMessages(),
		isOneWay:    template.IsOneWay(),
	}

	// Deep-copy the messages.
	templateMessages := template.Messages()
	pa.messages = make([]Message, 0, len(templateMessages))
	for _, v := range templateMessages {
		pa.messages = append(pa.messages, append(Message{}, v...))
	}

	// Apply the psk modifiers to all of the patterns.
	indexes := make(map[int]bool)
	splitModifier := strings.Split(modifier, "+")
	for _, v := range splitModifier {
		if !strings.HasPrefix(v, prefixPSK) {
			return nil, errors.New("nyquist/pattern: non-PSK modifier: " + v)
		}
		v = strings.TrimPrefix(v, prefixPSK)
		pskIndex, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("nyquist/pattern: failed to parse PSK index: " + err.Error())
		}

		if indexes[pskIndex] {
			return nil, errors.New("nyquist/pattern: redundant PSK modifier: " + prefixPSK + v)
		}
		if pskIndex < 0 || pskIndex > len(templateMessages) {
			return nil, errors.New("nyquist/pattern: invalid PSK modifier: " + prefixPSK + v)
		}
		switch pskIndex {
		case 0:
			pa.messages[0] = append(Message{Token_psk}, pa.messages[0]...)
		default:
			idx := pskIndex - 1
			pa.messages[idx] = append(pa.messages[idx], Token_psk)
		}
		indexes[pskIndex] = true
	}
	pa.numPSKs = len(indexes)

	return pa, nil
}

func mustMakePSK(template Pattern, modifier string) Pattern {
	pa, err := MakePSK(template, modifier)
	if err != nil {
		panic(err)
	}
	return pa
}
