package nyquist

var (
	// Patern_N is the N one-way handshake pattern.
	Pattern_N HandshakePattern = &builtInPattern{
		name: "N",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
		},
		isOneWay: true,
	}

	// Patern_K is the K one-way handshake pattern.
	Pattern_K HandshakePattern = &builtInPattern{
		name: "K",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_ss},
		},
		isOneWay: true,
	}

	// Patern_X is the X one-way handshake pattern.
	Pattern_X HandshakePattern = &builtInPattern{
		name: "X",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s, Token_ss},
		},
		isOneWay: true,
	}

	// Patern_Npsk0 is the Npsk0 one-way handshake pattern.
	Pattern_Npsk0 HandshakePattern = &builtInPattern{
		name: "Npsk0",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e, Token_es},
		},
		isPSK:    true,
		isOneWay: true,
	}

	// Patern_Kpsk0 is the Kpsk0 one-way handshake pattern.
	Pattern_Kpsk0 HandshakePattern = &builtInPattern{
		name: "Kpsk0",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e, Token_es, Token_ss},
		},
		isPSK:    true,
		isOneWay: true,
	}

	// Patern_Xpsk1 is the Xpsk1 one-way handshake pattern.
	Pattern_Xpsk1 HandshakePattern = &builtInPattern{
		name: "Xpsk1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s, Token_ss, Token_psk},
		},
		isPSK:    true,
		isOneWay: true,
	}
)
