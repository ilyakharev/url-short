package hasher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	lenShortURLPattern = 10
	shortURLPattern    = "^[a-zA-Z0-9_]{10}$"
)

func TestEncode(t *testing.T) {
	h := New()
	r := regexp.MustCompile(shortURLPattern)
	t.Run("shortUrl consists of provided alphabet", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			token, err := h.GenerateToken()
			require.NoError(t, err)
			assert.Regexp(t, r, token)
		}
	})

	t.Run("shortUrl has a specified length", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			token, err := h.GenerateToken()
			require.NoError(t, err)
			assert.Len(t, token, lenShortURLPattern)
		}
	})
}
