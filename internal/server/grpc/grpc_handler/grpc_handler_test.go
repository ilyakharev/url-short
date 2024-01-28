package grpchandler

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	mock_hasher "github.com/ilyakharev/url-short/internal/hasher/mock"
	"github.com/ilyakharev/url-short/internal/storage/inmemory"
	mock_storage "github.com/ilyakharev/url-short/internal/storage/mock"
	"github.com/ilyakharev/url-short/proto"
)

func TestSaveHandler(t *testing.T) {
	cases := []*struct {
		name        string
		request     *proto.CreateShortURLRequest
		expectToken string
		hashToken   string
		expectErr   bool
		failHash    bool
		failStorage bool
		prepareMock func(ctx context.Context, mem *mock_storage.MockStorager)
	}{
		{
			name: "Use bad URL",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http//wrong",
			},
			expectErr: true,
		},
		{
			name:      "Exception in generate expectToken",
			expectErr: true,
			failHash:  true,
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
		},
		{
			name:        "Success",
			expectErr:   false,
			expectToken: "1234567890",
			hashToken:   "1234567890",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
		},
		{
			name:        "Success",
			expectErr:   false,
			expectToken: "1234567890",
			hashToken:   "0123456789",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
		},
		{
			name:        "Check get error in storager AlreadyExists",
			expectErr:   true,
			failStorage: true,
			hashToken:   "0123456789",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, errors.New("some"))
			},
		},
		{
			name:        "Check get error in storager GetFullURL",
			expectErr:   true,
			failStorage: true,
			hashToken:   "0123456789",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, nil)
				mockMemory.EXPECT().GetFullURL(gomock.Any(), gomock.Any()).Return("", false, errors.New("some"))
			},
		},
		{
			name:        "Check get error in storager CreateShortURL",
			expectErr:   true,
			failStorage: true,
			hashToken:   "0123456789",
			request: &proto.CreateShortURLRequest{
				RawFullURL: "http://wro.ng",
			},
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, nil)
				mockMemory.EXPECT().GetFullURL(gomock.Any(), gomock.Any()).Return("", false, nil)
				mockMemory.EXPECT().CreateShortURL(gomock.Any(), gomock.Any(),
					gomock.Any()).Return(errors.New("some"))
			},
		},
	}
	memory := inmemory.New()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			hasher := mock_hasher.NewMockHasher(ctrl)
			if tc.failHash {
				hasher.EXPECT().GenerateToken().Return("", errors.New("any")).AnyTimes()
			} else {
				hasher.EXPECT().GenerateToken().Return(tc.hashToken, nil).AnyTimes()
			}

			var handler *GrpcHandler
			if tc.failStorage {
				mockMemory := mock_storage.NewMockStorager(ctrl)
				tc.prepareMock(ctx, mockMemory)
				handler = New(mockMemory, hasher, zap.NewNop())
			} else {
				handler = New(memory, hasher, zap.NewNop())
			}

			res, err := handler.CreateShortURL(ctx, tc.request)
			switch {
			case tc.expectErr && err == nil:
				t.Error("expect error")
			case !tc.expectErr && err != nil:
				t.Error("unexpected error: ", err)
			case !tc.expectErr && res.Token != tc.expectToken:
				t.Errorf("handler returned wrong token: got %v want %v",
					res.Token, tc.expectToken)
			}
		})
	}
}

func TestGetFullURL(t *testing.T) {
	cases := []*struct {
		name        string
		request     *proto.GetFullURLRequest
		expectURL   string
		expectErr   bool
		failStorage bool
		prepareMock func(ctx context.Context, mem *mock_storage.MockStorager)
	}{
		{
			name:      "Success",
			expectErr: false,
			expectURL: "http://ya.ru",
			request:   &proto.GetFullURLRequest{RawToken: "0123456789"},
		},
		{
			name:      "Not found",
			expectErr: true,
			request:   &proto.GetFullURLRequest{RawToken: "9876543210"},
		},
		{
			name:        "Check get error in storager GetFullURL",
			expectErr:   true,
			failStorage: true,
			request:     &proto.GetFullURLRequest{RawToken: "9876543210"},
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().GetFullURL(gomock.Any(), "9876543210").Return("", false, errors.New("some"))
			},
		},
	}
	ctx := context.Background()
	memory := inmemory.New()
	_ = memory.CreateShortURL(ctx, "http://ya.ru", "0123456789")
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			hasher := mock_hasher.NewMockHasher(ctrl)

			var handler *GrpcHandler
			if tc.failStorage {
				mockMemory := mock_storage.NewMockStorager(ctrl)
				tc.prepareMock(ctx, mockMemory)
				handler = New(mockMemory, hasher, zap.NewNop())
			} else {
				handler = New(memory, hasher, zap.NewNop())
			}

			res, err := handler.GetFullURL(ctx, tc.request)
			if tc.expectErr == false && err != nil {
				t.Errorf("Unexpected error")
			} else if tc.expectErr == false && res.FullURL != tc.expectURL {
				t.Errorf("handler returned wrong token: got %v want %v",
					res.FullURL, tc.expectURL)
			}
		})
	}
}
