package grpchandler

import (
	"context"
	"errors"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/ilyakharev/url-short/internal/hasher"
	"github.com/ilyakharev/url-short/internal/storage"
	"github.com/ilyakharev/url-short/proto"
)

type GrpcHandler struct {
	proto.UnimplementedGrpcHandlerServer
	storage storage.Storager
	hasher  hasher.Hasher
	logger  *zap.Logger
}

func (handler GrpcHandler) CreateShortURL(ctx context.Context,
	request *proto.CreateShortURLRequest,
) (*proto.CreateShortURLResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	handler.logger.Debug(
		"GetFullURL grpc request",
		zap.Any("raw_full_URL", request.RawFullURL),
	)

	_, err := url.ParseRequestURI(request.RawFullURL)
	if err != nil {
		handler.logger.Error("error on check URL:", zap.Error(err))
		return nil, err
	}

	token, exists, err := handler.storage.AlreadyExists(ctx, request.RawFullURL)
	if err != nil {
		handler.logger.Error("error on check URL on exists:", zap.Error(err))
		return nil, err
	}
	if exists {
		return &proto.CreateShortURLResponse{
			Token: token,
		}, nil
	}
	exists = true
	for exists {
		token, err = handler.hasher.GenerateToken()
		if err != nil {
			handler.logger.Error("error on generate expectToken:", zap.Error(err))
			return nil, err
		}
		_, exists, err = handler.storage.GetFullURL(ctx, token)
		if err != nil {
			handler.logger.Error("error on check URL:", zap.Error(err))
			return nil, err
		}
	}
	err = handler.storage.CreateShortURL(ctx, request.RawFullURL, token)
	if err != nil {
		handler.logger.Error("error on save expectToken:", zap.Error(err))
		return nil, err
	}
	return &proto.CreateShortURLResponse{
		Token: token,
	}, nil
}

func (handler GrpcHandler) GetFullURL(ctx context.Context,
	request *proto.GetFullURLRequest,
) (*proto.GetFullURLResponse, error) {
	handler.logger.Debug(
		"GetFullURL grpc request",
		zap.Any("raw_token", request.RawToken),
	)
	fullURL, ok, err := handler.storage.GetFullURL(ctx, request.RawToken)
	if err != nil {
		handler.logger.Error("error on get full URL:", zap.Error(err))
		return nil, err
	}
	if !ok {
		return nil, errors.New("not found")
	}
	return &proto.GetFullURLResponse{
		FullURL: fullURL,
	}, nil
}

func New(storage storage.Storager, hasher hasher.Hasher,
	logger *zap.Logger,
) *GrpcHandler {
	return &GrpcHandler{
		storage: storage,
		hasher:  hasher,
		logger:  logger,
	}
}
