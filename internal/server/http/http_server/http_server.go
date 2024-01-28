package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	httphandler "github.com/ilyakharev/url-short/internal/server/http/http_handler"
)

type HTTPServer struct {
	server *http.Server
	logger *zap.Logger
}

func New(port string, handler *httphandler.HTTPHandler,
	logger *zap.Logger,
) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:              ":" + port,
			Handler:           handler.CreateRouter(),
			ReadHeaderTimeout: time.Second,
		},
		logger: logger,
	}
}

func (server *HTTPServer) Run(ctx context.Context) error {
	ch := make(chan error)
	go func(chan error) {
		if err := server.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("failed to listen", zap.Error(err))
			ch <- err
		}
	}(ch)

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		server.logger.Info("Shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(ctx, time.Minute*1)
		defer cancel()

		if err := server.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return nil
	}
}
