package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fullURL = "https://mai.ru"
	token   = "mai"
)

func Test_StoreUrl(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		storage := New()
		defer func() {
			err := storage.Close()
			if err != nil {
				return
			}
		}()
		ctx := context.Background()

		err := storage.CreateShortURL(ctx, fullURL, token)
		require.NoError(t, err)

		url, _, err := storage.GetFullURL(ctx, token)
		require.NoError(t, err)
		assert.Equal(t, fullURL, url)
	})
	t.Run("not found", func(t *testing.T) {
		storage := New()
		defer func() {
			err := storage.Close()
			if err != nil {
				return
			}
		}()
		ctx := context.Background()

		_, ok, err := storage.GetFullURL(ctx, fullURL)
		require.NoError(t, err)
		assert.False(t, ok)
	})
	t.Run("check already exists", func(t *testing.T) {
		storage := New()
		defer func() {
			err := storage.Close()
			if err != nil {
				return
			}
		}()
		ctx := context.Background()

		_, found, err := storage.AlreadyExists(ctx, fullURL)
		require.NoError(t, err)
		assert.False(t, found)

		err = storage.CreateShortURL(ctx, fullURL, token)
		require.NoError(t, err)

		_, found, err = storage.AlreadyExists(ctx, fullURL)
		require.NoError(t, err)
		assert.True(t, found)
	})
}
