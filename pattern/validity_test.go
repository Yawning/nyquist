package pattern

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuiltInPatternValidity(t *testing.T) {
	for _, v := range supportedPatterns {
		t.Run(v.String(), func(t *testing.T) {
			require := require.New(t)

			err := IsValid(v)
			require.NoError(err, "IsValid(pattern)")
		})
	}
}
