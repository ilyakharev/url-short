package storage

import "context"

//go:generate mockgen -source=storager.go -destination=./mock/storager.go
type Storager interface {
	GetFullURL(ctx context.Context, token string) (fullURL string, found bool, err error)
	CreateShortURL(ctx context.Context, fullURL string, token string) (err error)
	AlreadyExists(ctx context.Context, fullURL string) (token string, found bool, err error)
	Close() error
}
