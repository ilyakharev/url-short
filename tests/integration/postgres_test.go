package integration

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ilyakharev/url-short/internal/storage/postgres"
)

func TestPostgres(t *testing.T) {
	t.Run("Test postgres", func(t *testing.T) {
		tt := struct {
			token   string
			fullURL string
		}{
			token:   "qwertyuiop",
			fullURL: "http://ozon.ru",
		}

		ctx := context.Background()
		st, err := postgres.New(os.Getenv("POSTGRES_URL"))
		require.NoError(t, err)
		defer func() {
			err = st.Close()
			if err != nil {
				return
			}
		}()

		fullURL, found, err := st.GetFullURL(ctx, tt.token)
		assert.False(t, found)
		assert.Empty(t, fullURL)
		require.NoError(t, err)

		token, found, err := st.AlreadyExists(ctx, tt.fullURL)
		assert.False(t, found)
		assert.Empty(t, token)
		require.NoError(t, err)

		err = st.CreateShortURL(ctx, tt.fullURL, tt.token)
		require.NoError(t, err)

		token, found, err = st.AlreadyExists(ctx, tt.fullURL)
		assert.True(t, found)
		assert.Equal(t, tt.token, token)
		require.NoError(t, err)

		fullURL, found, err = st.GetFullURL(ctx, tt.token)
		assert.True(t, found)
		assert.Equal(t, tt.fullURL, fullURL)
		require.NoError(t, err)
	})
}
