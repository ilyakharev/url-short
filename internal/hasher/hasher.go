package hasher

import (
	"math/big"
)

var (
	alphabet     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")
	alphabetSize = big.NewInt(int64(len(alphabet)))
	tokenLength  = 10
)

//go:generate mockgen -source=hasher.go -destination=./mock/hasher.go
type Hasher interface {
	GenerateToken() (token string, err error)
}
