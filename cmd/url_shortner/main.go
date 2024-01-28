package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ilyakharev/url-short/internal/hasher"
	"github.com/ilyakharev/url-short/internal/server"
	grpchandler "github.com/ilyakharev/url-short/internal/server/grpc/grpc_handler"
	grpcserver "github.com/ilyakharev/url-short/internal/server/grpc/grpc_server"
	httphandler "github.com/ilyakharev/url-short/internal/server/http/http_handler"
	httpserver "github.com/ilyakharev/url-short/internal/server/http/http_server"
	"github.com/ilyakharev/url-short/internal/storage"
	"github.com/ilyakharev/url-short/internal/storage/inmemory"
	"github.com/ilyakharev/url-short/internal/storage/postgres"
)

var logger *zap.Logger

func initLogger() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	logger, err = cfg.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(fmt.Errorf("error create logger: %w", err))
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	initLogger()

	portFlag, found := os.LookupEnv("PORT")
	if !found {
		portFlag = "80"
	}

	var err error
	var storager storage.Storager
	storageType, _ := os.LookupEnv("STORAGE_TYPE")
	switch storageType {
	case "postgres":
		logger.Info("Create postgres storager")
		storager, err = postgres.New(os.Getenv("POSTGRES_URL"))
		if err != nil {
			logger.Panic("unable to use postgres", zap.Error(err))
		}
	case "inmemory":
		logger.Info("Create in memory storager")
		storager = inmemory.New()
	default:
		logger.Panic("'STORAGE_TYPE' must be 'postgres' or 'inmemory'")
	}
	defer func() {
		err = storager.Close()
		if err != nil {
			logger.Error("unable to close storage", zap.Error(err))
		}
	}()

	hash := hasher.New()
	var srv server.Server
	transportType, _ := os.LookupEnv("TRANSPORT_TYPE")
	switch transportType {
	case "grpc":
		logger.Info("Create gRPC handler")
		handler := grpchandler.New(storager, hash, logger)
		logger.Info("Create gRPC server")
		srv = grpcserver.New(portFlag, handler, logger)
	case "http":
		logger.Info("Create HTTP handler")
		handler := httphandler.New(storager, hash, logger)
		logger.Info("Create HTTP server")
		srv = httpserver.New(portFlag, handler, logger)
	default:
		logger.Panic("'TRANSPORT_TYPE' must be 'grpc' or 'http'")
	}
	err = srv.Run(ctx)
	if err != nil {
		logger.Error("error in server", zap.Error(err))
	}
}
