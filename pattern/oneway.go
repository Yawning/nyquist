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

var (
	// N is the N one-way handshake pattern.
	N Pattern = &builtIn{
		name: "N",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
		},
		isOneWay: true,
	}

	// K is the K one-way handshake pattern.
	K Pattern = &builtIn{
		name: "K",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_ss},
		},
		isOneWay: true,
	}

	// X is the X one-way handshake pattern.
	X Pattern = &builtIn{
		name: "X",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s, Token_ss},
		},
		isOneWay: true,
	}

	// Npsk0 is the Npsk0 one-way handshake pattern.
	Npsk0 = mustMakePSK(N, "psk0")

	// Kpsk0 is the Kpsk0 one-way handshake pattern.
	Kpsk0 = mustMakePSK(K, "psk0")

	// Xpsk1 is the Xpsk1 one-way handshake pattern.
	Xpsk1 = mustMakePSK(X, "psk1")
)
