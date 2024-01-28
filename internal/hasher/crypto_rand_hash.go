package hasher

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type CryptoRandHash struct{}

var _ Hasher = CryptoRandHash{}

func (h CryptoRandHash) GenerateToken() (token string, err error) {
	var b strings.Builder

	for i := 0; i < tokenLength; i++ {
		n, _ := rand.Int(rand.Reader, alphabetSize)

		_, _ = fmt.Fprint(&b, string(alphabet[n.Int64()]))
	}

	return b.String(), nil
}

func New() CryptoRandHash {
	return CryptoRandHash{}
}
