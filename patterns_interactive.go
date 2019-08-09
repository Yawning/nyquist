package nyquist

var (
	// Pattern_NN is the NN interactive (fundemental) pattern.
	Pattern_NN HandshakePattern = &builtInPattern{
		name: "NN",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee},
		},
	}

	// Pattern_NK is the NK interactive (fundemental) pattern.
	Pattern_NK HandshakePattern = &builtInPattern{
		name: "NK",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
		},
	}

	// Pattern_NX is the NX interactive (fundemental) pattern.
	Pattern_NX HandshakePattern = &builtInPattern{
		name: "NX",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
		},
	}

	// Pattern_XN is the XN interactive (fundemental) pattern.
	Pattern_XN HandshakePattern = &builtInPattern{
		name: "XN",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s, Token_se},
		},
	}

	// Pattern_XK is the XK interactive (fundemental) pattern.
	Pattern_XK HandshakePattern = &builtInPattern{
		name: "XK",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s, Token_se},
		},
	}

	// Pattern_XX is the XX interactive (fundemental) pattern.
	Pattern_XX HandshakePattern = &builtInPattern{
		name: "XX",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
			MessagePattern{Token_s, Token_se},
		},
	}

	// Pattern_KN is the KN interactive (fundemental) pattern.
	Pattern_KN HandshakePattern = &builtInPattern{
		name: "KN",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
	}

	// Pattern_KK is the KK interactive (fundemental) pattern.
	Pattern_KK HandshakePattern = &builtInPattern{
		name: "KK",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_ss},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
	}

	// Pattern_KX is the KX interactive (fundemental) pattern.
	Pattern_KX HandshakePattern = &builtInPattern{
		name: "KX",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s, Token_es},
		},
	}

	// Pattern_IN is the IN interactive (fundemental) pattern.
	Pattern_IN HandshakePattern = &builtInPattern{
		name: "IN",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
	}

	// Pattern_IK is the IK interactive (fundemental) pattern.
	Pattern_IK HandshakePattern = &builtInPattern{
		name: "IK",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s, Token_ss},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
	}

	// Pattern_IX is the IX interactive (fundemental) pattern.
	Pattern_IX HandshakePattern = &builtInPattern{
		name: "IX",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s, Token_es},
		},
	}

	// Pattern_NNpsk0 is the NNpsk0 interactive (fundemental) pattern.
	Pattern_NNpsk0 HandshakePattern = &builtInPattern{
		name: "NNpsk0",
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e},
			MessagePattern{Token_e, Token_ee},
		},
		isPSK: true,
	}

	// Pattern_NNpsk2 is the NNpsk2 interactive (fundemental) pattern.
	Pattern_NNpsk2 HandshakePattern = &builtInPattern{
		name: "NNpsk2",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_NKpsk0 is the NKpsk0 interactive (fundemental) pattern.
	Pattern_NKpsk0 HandshakePattern = &builtInPattern{
		name: "NKpsk0",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
		},
		isPSK: true,
	}

	// Pattern_NKpsk2 is the NKpsk2 interactive (fundemental) pattern.
	Pattern_NKpsk2 HandshakePattern = &builtInPattern{
		name: "NKpsk2",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_NXpsk2 is the NXpsk2 interactive (fundemental) pattern.
	Pattern_NXpsk2 HandshakePattern = &builtInPattern{
		name: "NXpsk2",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_XNpsk3 is the XNpsk3 interactive (fundemental) pattern.
	Pattern_XNpsk3 HandshakePattern = &builtInPattern{
		name: "XNpsk3",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_XKpsk3 is the XKpsk3 interactive (fundemental) pattern.
	Pattern_XKpsk3 HandshakePattern = &builtInPattern{
		name: "XKpsk3",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es},
			MessagePattern{Token_e, Token_ee},
			MessagePattern{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_XXpsk3 is the XXpsk3 interactive (fundemental) pattern.
	Pattern_XXpsk3 HandshakePattern = &builtInPattern{
		name: "XXpsk3",
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_s, Token_es},
			MessagePattern{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_KNpsk0 is the KNpsk0 interactive (fundemental) pattern.
	Pattern_KNpsk0 HandshakePattern = &builtInPattern{
		name: "KNpsk0",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// Pattern_KNpsk2 is the KNpsk2 interactive (fundemental) pattern.
	Pattern_KNpsk2 HandshakePattern = &builtInPattern{
		name: "KNpsk2",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_KKpsk0 is the KKpsk0 interactive (fundemental) pattern.
	Pattern_KKpsk0 HandshakePattern = &builtInPattern{
		name: "KKpsk0",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_psk, Token_e, Token_es, Token_ss},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// Pattern_KKpsk2 is the KKpsk2 interactive (fundemental) pattern.
	Pattern_KKpsk2 HandshakePattern = &builtInPattern{
		name: "KKpsk2",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_ss},
			MessagePattern{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_KXpsk2 is the KXpsk2 interactive (fundemental) pattern.
	Pattern_KXpsk2 HandshakePattern = &builtInPattern{
		name: "KXpsk2",
		preMessages: []MessagePattern{
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_INpsk1 is the INpsk1 interactive (fundemental) pattern.
	Pattern_INpsk1 HandshakePattern = &builtInPattern{
		name: "INpsk1",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s, Token_psk},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// Pattern_INpsk2 is the INpsk2 interactive (fundemental) pattern.
	Pattern_INpsk2 HandshakePattern = &builtInPattern{
		name: "INpsk2",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_IKpsk1 is the IKpsk1 interactive (fundemental) pattern.
	Pattern_IKpsk1 HandshakePattern = &builtInPattern{
		name: "IKpsk1",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s, Token_ss, Token_psk},
			MessagePattern{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// Pattern_IKpsk2 is the IKpsk2 interactive (fundemental) pattern.
	Pattern_IKpsk2 HandshakePattern = &builtInPattern{
		name: "IKpsk2",
		preMessages: []MessagePattern{
			nil,
			MessagePattern{Token_s},
		},
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_es, Token_s, Token_ss},
			MessagePattern{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// Pattern_IXpsk2 is the IXpsk2 interactive (fundemental) pattern.
	Pattern_IXpsk2 HandshakePattern = &builtInPattern{
		name: "IXpsk2",
		messages: []MessagePattern{
			MessagePattern{Token_e, Token_s},
			MessagePattern{Token_e, Token_ee, Token_se, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}
)
