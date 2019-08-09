package internal

func ExplicitBzero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
