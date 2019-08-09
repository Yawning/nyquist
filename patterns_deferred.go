package nyquist

var (
	// Pattern_NK1 is the NK1 deferred pattern.
	Pattern_NK1 HandshakePattern = &builtInPattern{
		name: "NK1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_es},
		},
	}

	// Pattern_NX1 is the NX1 deferred pattern.
	Pattern_NX1 HandshakePattern = &builtInPattern{
		name: "NX1",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s},
			MessagePattern{Token_es},
		},
	}

	// Pattern_X1N is the X1N deferred pattern.
	Pattern_X1N HandshakePattern = &builtInPattern{
		name: "X1N",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s},
			MessagePattern{Token_se},
		},
	}

	// Pattern_X1K is the X1K deferred pattern.
	Pattern_X1K HandshakePattern = &builtInPattern{
		name: "X1K",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s},
			MessagePattern{Token_se},
		},
	}

	// Pattern_XK1 is the XK1 deferred pattern.
	Pattern_XK1 HandshakePattern = &builtInPattern{
		name: "XK1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_es},
			MessagePattern{Token_s, Token_se},
		},
	}

	// Pattern_X1K1 is the X1K1 deferred pattern.
	Pattern_X1K1 HandshakePattern = &builtInPattern{
		name: "X1K1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_es},
			MessagePattern{Token_s},
			MessagePattern{Token_se},
		},
	}

	// Pattern_X1X is the X1X deferred pattern.
	Pattern_X1X HandshakePattern = &builtInPattern{
		name: "X1X",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
			MessagePattern{Token_s},
			MessagePattern{Token_se},
		},
	}

	// Pattern_XX1 is the XX1 deferred pattern.
	Pattern_XX1 HandshakePattern = &builtInPattern{
		name: "XX1",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s},
			MessagePattern{Token_es, Token_s, Token_se},
		},
	}

	// Pattern_X1X1 is the X1X1 deferred pattern.
	Pattern_X1X1 HandshakePattern = &builtInPattern{
		name: "X1X1",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s},
			MessagePattern{Token_es, Token_s},
			MessagePattern{Token_se},
		},
	}

	// Pattern_K1N is the K1N deferred pattern.
	Pattern_K1N HandshakePattern = &builtInPattern{
		name: "K1N",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_se},
		},
	}

	// Pattern_K1K is the K1K deferred pattern.
	Pattern_K1K HandshakePattern = &builtInPattern{
		name: "K1K",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_se},
		},
	}

	// Pattern_KK1 is the KK1 deferred pattern.
	Pattern_KK1 HandshakePattern = &builtInPattern{
		name: "KK1",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se, Token_es},
		},
	}

	// Pattern_K1K1 is the K1K1 deferred pattern.
	Pattern_K1K1 HandshakePattern = &builtInPattern{
		name: "K1K1",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_es},
			MessagePattern{Token_se},
		},
	}

	// Pattern_K1X is the K1X deferred pattern.
	Pattern_K1X HandshakePattern = &builtInPattern{
		name: "K1X",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
			MessagePattern{Token_se},
		},
	}

	// Pattern_KX1 is the KX1 deferred pattern.
	Pattern_KX1 HandshakePattern = &builtInPattern{
		name: "KX1",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s},
			MessagePattern{Token_es},
		},
	}

	// Pattern_K1X1 is the K1X1 deferred pattern.
	Pattern_K1X1 HandshakePattern = &builtInPattern{
		name: "K1X1",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s},
			MessagePattern{Token_se, Token_es},
		},
	}

	// Pattern_I1N is the I1N deferred pattern.
	Pattern_I1N HandshakePattern = &builtInPattern{
		name: "I1N",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_se},
		},
	}

	// Pattern_I1K is the I1K deferred pattern.
	Pattern_I1K HandshakePattern = &builtInPattern{
		name: "I1K",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_se},
		},
	}

	// Pattern_IK1 is the IK1 deferred pattern.
	Pattern_IK1 HandshakePattern = &builtInPattern{
		name: "IK1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se, Token_es},
		},
	}

	// Pattern_I1K1 is the I1K1 deferred pattern.
	Pattern_I1K1 HandshakePattern = &builtInPattern{
		name: "I1K1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_es},
			MessagePattern{Token_se},
		},
	}

	// Pattern_I1X is the I1X deferred pattern.
	Pattern_I1X HandshakePattern = &builtInPattern{
		name: "I1X",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
			MessagePattern{Token_se},
		},
	}

	// Pattern_IX1 is the IX1 deferred pattern.
	Pattern_IX1 HandshakePattern = &builtInPattern{
		name: "IX1",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s},
			MessagePattern{Token_es},
		},
	}

	// Pattern_I1X1 is the I1X1 deferred pattern.
	Pattern_I1X1 HandshakePattern = &builtInPattern{
		name: "I1X1",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_s},
			MessagePattern{Token_se, Token_es},
		},
	}
)
