package inmemory

import (
	"context"
	"sync"

	"github.com/ilyakharev/url-short/internal/storage"
)

type Inmemory struct {
	mutex       sync.RWMutex
	shortToFull map[string]string
	fullToShort map[string]string
}

var _ storage.Storager = &Inmemory{}

func New() *Inmemory {
	return &Inmemory{
		mutex:       sync.RWMutex{},
		shortToFull: make(map[string]string),
		fullToShort: make(map[string]string),
	}
}

func (storage *Inmemory) GetFullURL(_ context.Context,
	token string,
) (fullURL string, found bool, err error) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	fullURL, found = storage.shortToFull[token]
	if !found {
		return "", found, err
	}
	return fullURL, found, err
}

func (storage *Inmemory) CreateShortURL(_ context.Context, fullURL string,
	token string,
) (err error) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	storage.fullToShort[fullURL] = token
	storage.shortToFull[token] = fullURL

	return nil
}

func (storage *Inmemory) AlreadyExists(_ context.Context,
	fullURL string,
) (token string, found bool, err error) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()
	token, found = storage.fullToShort[fullURL]
	if found {
		return token, found, nil
	}
	return "", found, nil
}

func (storage *Inmemory) Close() error {
	return nil
}
