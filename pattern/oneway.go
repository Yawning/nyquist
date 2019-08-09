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
	Npsk0 Pattern = &builtIn{
		name: "Npsk0",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_psk, Token_e, Token_es},
		},
		isPSK:    true,
		isOneWay: true,
	}

	// Kpsk0 is the Kpsk0 one-way handshake pattern.
	Kpsk0 Pattern = &builtIn{
		name: "Kpsk0",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_psk, Token_e, Token_es, Token_ss},
		},
		isPSK:    true,
		isOneWay: true,
	}

	// Xpsk1 is the Xpsk1 one-way handshake pattern.
	Xpsk1 Pattern = &builtIn{
		name: "Xpsk1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s, Token_ss, Token_psk},
		},
		isPSK:    true,
		isOneWay: true,
	}
)
