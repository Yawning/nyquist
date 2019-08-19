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
			return nil, errors.New("nyquist/pattern: reduntant PSK modifiler: " + prefixPSK + v)
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
