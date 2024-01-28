package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	mock_hasher "github.com/ilyakharev/url-short/internal/hasher/mock"
	grpchandler "github.com/ilyakharev/url-short/internal/server/grpc/grpc_handler"
	"github.com/ilyakharev/url-short/internal/storage/inmemory"
)

func TestServer(t *testing.T) {
	t.Run("Create server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		hasher := mock_hasher.NewMockHasher(ctrl)
		memory := inmemory.New()
		handler := grpchandler.New(memory, hasher, zap.NewNop())
		srv := New("81", handler, zap.NewNop())
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		err := srv.Run(ctx)
		require.NoError(t, err)
	})
}
