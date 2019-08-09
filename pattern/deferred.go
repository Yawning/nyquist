package pattern

var (
	// NK1 is the NK1 deferred pattern.
	NK1 Pattern = &builtIn{
		name: "NK1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_es},
		},
	}

	// NX1 is the NX1 deferred pattern.
	NX1 Pattern = &builtIn{
		name: "NX1",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s},
			Message{Token_es},
		},
	}

	// X1N is the X1N deferred pattern.
	X1N Pattern = &builtIn{
		name: "X1N",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee},
			Message{Token_s},
			Message{Token_se},
		},
	}

	// X1K is the X1K deferred pattern.
	X1K Pattern = &builtIn{
		name: "X1K",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee},
			Message{Token_s},
			Message{Token_se},
		},
	}

	// XK1 is the XK1 deferred pattern.
	XK1 Pattern = &builtIn{
		name: "XK1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_es},
			Message{Token_s, Token_se},
		},
	}

	// X1K1 is the X1K1 deferred pattern.
	X1K1 Pattern = &builtIn{
		name: "X1K1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_es},
			Message{Token_s},
			Message{Token_se},
		},
	}

	// X1X is the X1X deferred pattern.
	X1X Pattern = &builtIn{
		name: "X1X",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es},
			Message{Token_s},
			Message{Token_se},
		},
	}

	// XX1 is the XX1 deferred pattern.
	XX1 Pattern = &builtIn{
		name: "XX1",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s},
			Message{Token_es, Token_s, Token_se},
		},
	}

	// X1X1 is the X1X1 deferred pattern.
	X1X1 Pattern = &builtIn{
		name: "X1X1",
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s},
			Message{Token_es, Token_s},
			Message{Token_se},
		},
	}

	// K1N is the K1N deferred pattern.
	K1N Pattern = &builtIn{
		name: "K1N",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee},
			Message{Token_se},
		},
	}

	// K1K is the K1K deferred pattern.
	K1K Pattern = &builtIn{
		name: "K1K",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es},
			Message{Token_e, Token_ee},
			Message{Token_se},
		},
	}

	// KK1 is the KK1 deferred pattern.
	KK1 Pattern = &builtIn{
		name: "KK1",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se, Token_es},
		},
	}

	// K1K1 is the K1K1 deferred pattern.
	K1K1 Pattern = &builtIn{
		name: "K1K1",
		preMessages: []Message{
			Message{Token_s},
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_es},
			Message{Token_se},
		},
	}

	// K1X is the K1X deferred pattern.
	K1X Pattern = &builtIn{
		name: "K1X",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s, Token_es},
			Message{Token_se},
		},
	}

	// KX1 is the KX1 deferred pattern.
	KX1 Pattern = &builtIn{
		name: "KX1",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_se, Token_s},
			Message{Token_es},
		},
	}

	// K1X1 is the K1X1 deferred pattern.
	K1X1 Pattern = &builtIn{
		name: "K1X1",
		preMessages: []Message{
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e},
			Message{Token_e, Token_ee, Token_s},
			Message{Token_se, Token_es},
		},
	}

	// I1N is the I1N deferred pattern.
	I1N Pattern = &builtIn{
		name: "I1N",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee},
			Message{Token_se},
		},
	}

	// I1K is the I1K deferred pattern.
	I1K Pattern = &builtIn{
		name: "I1K",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_es, Token_s},
			Message{Token_e, Token_ee},
			Message{Token_se},
		},
	}

	// IK1 is the IK1 deferred pattern.
	IK1 Pattern = &builtIn{
		name: "IK1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se, Token_es},
		},
	}

	// I1K1 is the I1K1 deferred pattern.
	I1K1 Pattern = &builtIn{
		name: "I1K1",
		preMessages: []Message{
			nil,
			Message{Token_s},
		},
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_es},
			Message{Token_se},
		},
	}

	// I1X is the I1X deferred pattern.
	I1X Pattern = &builtIn{
		name: "I1X",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_s, Token_es},
			Message{Token_se},
		},
	}

	// IX1 is the IX1 deferred pattern.
	IX1 Pattern = &builtIn{
		name: "IX1",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_se, Token_s},
			Message{Token_es},
		},
	}

	// I1X1 is the I1X1 deferred pattern.
	I1X1 Pattern = &builtIn{
		name: "I1X1",
		messages: []Message{
			Message{Token_e, Token_s},
			Message{Token_e, Token_ee, Token_s},
			Message{Token_se, Token_es},
		},
	}
)
