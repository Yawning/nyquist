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
	"fmt"
)

// IsValid checks a pattern for validity according to the handshake pattern
// validity rules, and implementation limitations.
//
// Warning: This is not particularly fast, and should only be called when
// validating custom patterns, or testing.
func IsValid(pa Pattern) error {
	initTokens := make(map[Token]bool)
	respTokens := make(map[Token]bool)

	getSide := func(idx int) (map[Token]bool, bool, string) {
		isInitiator := idx&1 == 0
		if isInitiator {
			return initTokens, true, "initiator"
		}
		return respTokens, false, "responder"
	}

	inEither := func(t Token) bool {
		return initTokens[t] || respTokens[t]
	}

	inBoth := func(t Token) bool {
		return initTokens[t] && respTokens[t]
	}

	// Sanity-check the pre-messages.
	preMessages := pa.PreMessages()
	if len(preMessages) > 2 {
		return errors.New("nyquist/pattern: excessive pre-messages")
	}
	for i, msg := range preMessages {
		m, _, side := getSide(i)
		for _, v := range msg {
			switch v {
			case Token_e, Token_s:
				// 2. Parties must not send their static public key or ephemeral
				// public key more than once per handshake.
				if m[v] {
					return fmt.Errorf("nyquist/pattern: redundant pre-message token (%s): %s", side, v)
				}
				m[v] = true
			default:
				return fmt.Errorf("nyquist/pattern: invalid pre-message token: %s", v)
			}
		}
	}

	// Validate the messages.
	messages := pa.Messages()
	if len(messages) == 0 {
		return errors.New("nyquist/pattern: no handshake messages")
	}
	if pa.IsOneWay() && len(messages) != 1 {
		return errors.New("nyquist/pattern: excessive messages for one-way pattern")
	}
	var numDHs int
	for i, msg := range messages {
		m, isInitiator, side := getSide(i)
		for _, v := range msg {
			switch v {
			case Token_e, Token_s:
				// 2. Parties must not send their static public key or ephemeral
				// public key more than once per handshake.
				if m[v] {
					return fmt.Errorf("nyquist/pattern: redundant public key (%s): %s", side, v)
				}
			case Token_ee, Token_es, Token_se, Token_ss:
				// 3. Parties must not perform a DH calculation more than once
				// per handshake.
				if inEither(v) {
					return fmt.Errorf("nyquist/pattern: redundant DH calcuation: %s", v)
				}
				numDHs++
			case Token_psk:
				// Technically the spec supports multiple PSKs, though no
				// standard patterns are defined at this time.  Some
				// implementations support this, this one does not.
				if inEither(v) {
					return fmt.Errorf("nyquist/pattern: redundant pre-shared key")
				}
			default:
				return fmt.Errorf("nyquist/pattern: invalid message token: %s", v)
			}

			// 1. Parties can only perform DH between private keys and public
			// keys they posess.
			var impossibleDH Token
			switch v {
			case Token_ee:
				if !inBoth(Token_e) {
					impossibleDH = v
				}
			case Token_ss:
				if !inBoth(Token_s) {
					impossibleDH = v
				}
			case Token_es:
				if !initTokens[Token_e] || !respTokens[Token_s] {
					impossibleDH = v
				}
			case Token_se:
				if !initTokens[Token_s] || !respTokens[Token_e] {
					impossibleDH = v
				}
			default:
			}
			if impossibleDH != Token_invalid {
				return fmt.Errorf("nyquist/pattern: impossible DH: %s", v)
			}

			m[v] = true
		}

		// 4. After performing a DH between a remote public key (either static
		// or ephemeral) and the local static key, the local party must not
		// call ENCRYPT() unless it has also performed a DH between its local
		// ephemeral key and the remote public key.
		var missingDH Token
		if isInitiator {
			if inEither(Token_se) && !inEither(Token_ee) {
				missingDH = Token_ee
			}
			if inEither(Token_ss) && !inEither(Token_es) {
				missingDH = Token_es
			}
		} else {
			if inEither(Token_es) && !inEither(Token_ee) {
				missingDH = Token_ee
			}
			if inEither(Token_ss) && !inEither(Token_se) {
				missingDH = Token_se
			}
		}
		if missingDH != Token_invalid {
			return fmt.Errorf("nyquist/pattern: missing DH calculation (%s): %s", side, missingDH)
		}

		if inEither(Token_psk) {
			// A party may not send any encrypted data after it processes a
			// "psk" token unless it has previously sent an epmeheral public
			// key (an "e" token), either before or after the "psk" token.
			if !m[Token_e] {
				return fmt.Errorf("nyquist/pattern: payload after pre-shared key without ephemeral (%s)", side)
			}
		}
	}

	// Patterns without any DH calculations may be "valid", but are
	// nonsensical.
	if numDHs == 0 {
		return errors.New("nyquist/pattern: no DH calculations at all")
	}

	// Make sure the PSK hint interface function is implemented correctly.
	if inEither(Token_psk) != pa.IsPSK() {
		return errors.New("nyquist/pattern: IsPSK() mismatch with (pre-)messages")
	}

	return nil
}
