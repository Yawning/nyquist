package pattern

var (
	// NN is the NN interactive (fundemental) pattern.
	NN Pattern = &builtIn{
		name: "NN",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee},
		},
	}

	// NK is the NK interactive (fundemental) pattern.
	NK Pattern = &builtIn{
		name: "NK",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee},
		},
	}

	// NX is the NX interactive (fundemental) pattern.
	NX Pattern = &builtIn{
		name: "NX",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es},
		},
	}

	// XN is the XN interactive (fundemental) pattern.
	XN Pattern = &builtIn{
		name: "XN",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee},
			Message{Token_s, Token_se},
		},
	}

	// XK is the XK interactive (fundemental) pattern.
	XK Pattern = &builtIn{
		name: "XK",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee},
			Message{Token_s, Token_se},
		},
	}

	// XX is the XX interactive (fundemental) pattern.
	XX Pattern = &builtIn{
		name: "XX",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es},
			Message{Token_s, Token_se},
		},
	}

	// KN is the KN interactive (fundemental) pattern.
	KN Pattern = &builtIn{
		name: "KN",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se},
		},
	}

	// KK is the KK interactive (fundemental) pattern.
	KK Pattern = &builtIn{
		name: "KK",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_ss},
			Message{Token_e, Token_ee, Token_se},
		},
	}

	// KX is the KX interactive (fundemental) pattern.
	KX Pattern = &builtIn{
		name: "KX",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se, Token_s, Token_es},
		},
	}

	// IN is the IN interactive (fundemental) pattern.
	IN Pattern = &builtIn{
		name: "IN",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se},
		},
	}

	// IK is the IK interactive (fundemental) pattern.
	IK Pattern = &builtIn{
		name: "IK",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s, Token_ss},
			Message{Token_e, Token_ee, Token_se},
		},
	}

	// IX is the IX interactive (fundemental) pattern.
	IX Pattern = &builtIn{
		name: "IX",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se, Token_s, Token_es},
		},
	}

	// NNpsk0 is the NNpsk0 interactive (fundemental) pattern.
	NNpsk0 Pattern = &builtIn{
		name: "NNpsk0",
		messages: []Message{
			Message{Token_psk, Token_e},
			Message{Token_e, Token_ee},
		},
		isPSK: true,
	}

	// NNpsk2 is the NNpsk2 interactive (fundemental) pattern.
	NNpsk2 Pattern = &builtIn{
		name: "NNpsk2",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_psk},
		},
		isPSK: true,
	}

	// NKpsk0 is the NKpsk0 interactive (fundemental) pattern.
	NKpsk0 Pattern = &builtIn{
		name: "NKpsk0",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_psk, Token_e, Token_es},
			Message{Token_e, Token_ee},
		},
		isPSK: true,
	}

	// NKpsk2 is the NKpsk2 interactive (fundemental) pattern.
	NKpsk2 Pattern = &builtIn{
		name: "NKpsk2",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee, Token_psk},
		},
		isPSK: true,
	}

	// NXpsk2 is the NXpsk2 interactive (fundemental) pattern.
	NXpsk2 Pattern = &builtIn{
		name: "NXpsk2",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}

	// XNpsk3 is the XNpsk3 interactive (fundemental) pattern.
	XNpsk3 Pattern = &builtIn{
		name: "XNpsk3",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee},
			Message{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// XKpsk3 is the XKpsk3 interactive (fundemental) pattern.
	XKpsk3 Pattern = &builtIn{
		name: "XKpsk3",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee},
			Message{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// XXpsk3 is the XXpsk3 interactive (fundemental) pattern.
	XXpsk3 Pattern = &builtIn{
		name: "XXpsk3",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es},
			Message{Token_s, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// KNpsk0 is the KNpsk0 interactive (fundemental) pattern.
	KNpsk0 Pattern = &builtIn{
		name: "KNpsk0",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_psk, Token_e},
			Message{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// KNpsk2 is the KNpsk2 interactive (fundemental) pattern.
	KNpsk2 Pattern = &builtIn{
		name: "KNpsk2",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// KKpsk0 is the KKpsk0 interactive (fundemental) pattern.
	KKpsk0 Pattern = &builtIn{
		name: "KKpsk0",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_psk, Token_e, Token_es, Token_ss},
			Message{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// KKpsk2 is the KKpsk2 interactive (fundemental) pattern.
	KKpsk2 Pattern = &builtIn{
		name: "KKpsk2",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_ss},
			Message{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// KXpsk2 is the KXpsk2 interactive (fundemental) pattern.
	KXpsk2 Pattern = &builtIn{
		name: "KXpsk2",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}

	// INpsk1 is the INpsk1 interactive (fundemental) pattern.
	INpsk1 Pattern = &builtIn{
		name: "INpsk1",
		messages: []Message{
			Message{Token_e, Token_s, Token_psk},
			Message{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// INpsk2 is the INpsk2 interactive (fundemental) pattern.
	INpsk2 Pattern = &builtIn{
		name: "INpsk2",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// IKpsk1 is the IKpsk1 interactive (fundemental) pattern.
	IKpsk1 Pattern = &builtIn{
		name: "IKpsk1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s, Token_ss, Token_psk},
			Message{Token_e, Token_ee, Token_se},
		},
		isPSK: true,
	}

	// IKpsk2 is the IKpsk2 interactive (fundemental) pattern.
	IKpsk2 Pattern = &builtIn{
		name: "IKpsk2",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s, Token_ss},
			Message{Token_e, Token_ee, Token_se, Token_psk},
		},
		isPSK: true,
	}

	// IXpsk2 is the IXpsk2 interactive (fundemental) pattern.
	IXpsk2 Pattern = &builtIn{
		name: "IXpsk2",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se, Token_s, Token_es, Token_psk},
		},
		isPSK: true,
	}
)
