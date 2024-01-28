package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	mock_hasher "github.com/ilyakharev/url-short/internal/hasher/mock"
	"github.com/ilyakharev/url-short/internal/storage/inmemory"
	mock_storage "github.com/ilyakharev/url-short/internal/storage/mock"
)

func TestSaveHandler(t *testing.T) {
	cases := []*struct {
		name        string
		rawURL      string
		expectToken string
		hashToken   string
		method      string
		statusCode  int
		failHash    bool
		failStorage bool
		prepareMock func(ctx context.Context, mem *mock_storage.MockStorager)
	}{
		{
			name:       "Use GET method",
			rawURL:     "https://ya.ru/",
			method:     http.MethodGet,
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "Use bad URL",
			rawURL:     "http//wrong",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Exception in generate token",
			rawURL:     "http://wro.ng",
			hashToken:  "",
			method:     http.MethodPost,
			statusCode: http.StatusInternalServerError,
			failHash:   true,
		},
		{
			name:        "Success",
			rawURL:      "http://ya.ru",
			expectToken: "1234567890",
			hashToken:   "1234567890",
			method:      http.MethodPost,
			statusCode:  http.StatusCreated,
			failHash:    false,
		},
		{
			name:        "Expect same token",
			rawURL:      "http://ya.ru",
			expectToken: "1234567890",
			hashToken:   "9876543210",
			method:      http.MethodPost,
			statusCode:  http.StatusOK,
			failHash:    false,
		},
		{
			name:        "Check get error in storager AlreadyExists",
			rawURL:      "http://ya.ru",
			expectToken: "1234567890",
			hashToken:   "1234567890",
			method:      http.MethodPost,
			statusCode:  http.StatusInternalServerError,
			failHash:    false,
			failStorage: true,
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, errors.New("some"))
			},
		},
		{
			name:        "Check get error in storager GetFullURL",
			rawURL:      "http://ya.ru",
			expectToken: "1234567890",
			hashToken:   "1234567890",
			method:      http.MethodPost,
			statusCode:  http.StatusInternalServerError,
			failHash:    false,
			failStorage: true,
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, nil)
				mockMemory.EXPECT().GetFullURL(gomock.Any(), gomock.Any()).Return("", false, errors.New("some"))
			},
		},
		{
			name:        "Check get error in storager CreateShortURL",
			rawURL:      "http://ya.ru",
			expectToken: "1234567890",
			hashToken:   "1234567890",
			method:      http.MethodPost,
			statusCode:  http.StatusInternalServerError,
			failHash:    false,
			failStorage: true,
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().AlreadyExists(gomock.Any(), gomock.Any()).Return("", false, nil)
				mockMemory.EXPECT().GetFullURL(gomock.Any(), gomock.Any()).Return("",
					false, nil)
				mockMemory.EXPECT().CreateShortURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some"))
			},
		},
	}
	memory := inmemory.New()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			hasher := mock_hasher.NewMockHasher(ctrl)
			var b bytes.Buffer
			_, _ = fmt.Fprint(&b, tc.rawURL)

			if tc.failHash {
				hasher.EXPECT().GenerateToken().Return("", errors.New("any")).AnyTimes()
			} else {
				hasher.EXPECT().GenerateToken().Return(tc.hashToken, nil).AnyTimes()
			}

			ctx := context.Background()

			var handler *HTTPHandler
			if tc.failStorage {
				mockMemory := mock_storage.NewMockStorager(ctrl)
				tc.prepareMock(ctx, mockMemory)
				handler = New(mockMemory, hasher, zap.NewNop())
			} else {
				handler = New(memory, hasher, zap.NewNop())
			}
			req, err := http.NewRequestWithContext(ctx, tc.method, "/create", &b)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			handler.CreateShortURL(rr, req)
			status := rr.Code
			if status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.statusCode)
			} else if status == http.StatusCreated {
				expectResult, _ := json.Marshal(
					struct {
						Code    int    `json:"code"`    // token
						Message string `json:"message"` // short URL
					}{
						Code:    tc.statusCode,
						Message: tc.expectToken,
					})
				if !bytes.Equal(expectResult, rr.Body.Bytes()) {
					t.Error("handler returned wrong body")
				}
			}
		})
	}
}

func TestGetFullUrl(t *testing.T) {
	cases := []*struct {
		name        string
		method      string
		token       string
		expectURL   string
		statusCode  int
		failStorage bool
		prepareMock func(ctx context.Context, mem *mock_storage.MockStorager)
	}{
		{
			name:       "Use POST method",
			token:      "0123456789",
			method:     http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "Use bad token",
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Success",
			token:      "0123456789",
			expectURL:  "http://ya.ru",
			method:     http.MethodGet,
			statusCode: http.StatusFound,
		},
		{
			name:       "Not found",
			token:      "9876543210",
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
		},
		{
			name:        "Check get error in storager GetFullURL",
			token:       "9876543210",
			method:      http.MethodGet,
			statusCode:  http.StatusInternalServerError,
			failStorage: true,
			prepareMock: func(ctx context.Context, mockMemory *mock_storage.MockStorager) {
				mockMemory.EXPECT().GetFullURL(gomock.Any(), gomock.Any()).Return("", false, errors.New("some"))
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

			var handler *HTTPHandler
			if tc.failStorage {
				mockMemory := mock_storage.NewMockStorager(ctrl)
				tc.prepareMock(ctx, mockMemory)
				handler = New(mockMemory, hasher, zap.NewNop())
			} else {
				handler = New(memory, hasher, zap.NewNop())
			}

			req, err := http.NewRequestWithContext(context.Background(), tc.method, "/"+tc.token, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			handler.GetFullURL(rr, req)
			status := rr.Code
			if status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.statusCode)
			}
		})
	}
}
