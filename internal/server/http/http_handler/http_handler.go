package httphandler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/ilyakharev/url-short/internal/hasher"
	"github.com/ilyakharev/url-short/internal/storage"
)

//go:generate mockgen -source=httpHandler.go -destination=./mock/httpHandler.go
type HTTPHandler struct {
	storager storage.Storager
	hasher   hasher.Hasher
	logger   *zap.Logger
}

func New(st storage.Storager, h hasher.Hasher, logger *zap.Logger) *HTTPHandler {
	return &HTTPHandler{storager: st, hasher: h, logger: logger}
}

func (handler *HTTPHandler) CreateRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", handler.CreateShortURL)
	mux.HandleFunc("/", handler.GetFullURL)
	return mux
}

func (handler *HTTPHandler) CreateShortURL(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), time.Second)
	defer cancel()

	handler.logger.Debug(
		"CreateShortURL http request",
		zap.Any("address", request.RemoteAddr),
		zap.Any("method", request.Method),
		zap.Any("url", request.URL),
	)

	if request.Method != http.MethodPost {
		handler.sendResponse(http.StatusMethodNotAllowed, writer, "Method is not allowed")
		return
	}
	var err error
	var token string

	writer.Header().Add("Content-Type", "application/json")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		handler.logger.Error("error on read full URL", zap.Error(err))
		handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
	}

	rawURL := string(body)
	_, err = url.ParseRequestURI(rawURL)
	if err != nil {
		handler.sendResponse(http.StatusBadRequest, writer, "Invalid URL")
		return
	}

	token, exists, err := handler.storager.AlreadyExists(ctx, rawURL)
	if err != nil {
		handler.logger.Error("error on check url on exists", zap.Error(err))
		handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
		return
	}
	if exists {
		handler.sendResponse(http.StatusOK, writer, token)
		return
	}
	exists = true
	for exists {
		token, err = handler.hasher.GenerateToken()
		if err != nil {
			handler.logger.Error("error on generate token", zap.Error(err))
			handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
			return
		}
		_, exists, err = handler.storager.GetFullURL(ctx, token)
		if err != nil {
			handler.logger.Error("error on check url:", zap.Error(err))
			handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
			return
		}
	}

	err = handler.storager.CreateShortURL(ctx, rawURL, token)
	if err != nil {
		handler.logger.Error("error on save token", zap.Error(err))
		handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
		return
	}
	handler.sendResponse(http.StatusCreated, writer, token)
}

func (handler *HTTPHandler) GetFullURL(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), time.Second)
	defer cancel()

	handler.logger.Debug(
		"GetFullURL http request",
		zap.Any("address", request.RemoteAddr),
		zap.Any("method", request.Method),
		zap.Any("url", request.URL),
	)

	if request.Method != http.MethodGet {
		handler.sendResponse(http.StatusMethodNotAllowed, writer, "Method is not allowed")
		return
	}
	writer.Header().Add("Content-Type", "application/json")

	rawShortURL := request.URL.Path[1:]

	fullURL, ok, err := handler.storager.GetFullURL(ctx, rawShortURL)
	if err != nil {
		handler.logger.Error("error on get full url", zap.Error(err))
		handler.sendResponse(http.StatusInternalServerError, writer, err.Error())
		return
	}
	if !ok {
		handler.sendResponse(http.StatusNotFound, writer, "Not found")
		return
	}

	http.Redirect(writer, request, fullURL, http.StatusFound)
}

func (handler *HTTPHandler) sendResponse(code int, w http.ResponseWriter, message string) {
	w.WriteHeader(code)
	resp, err := json.Marshal(
		struct {
			Code    int    `json:"code"`    // token
			Message string `json:"message"` // short URL
		}{
			Code:    code,
			Message: message,
		})
	if err != nil {
		handler.logger.Error("error while marshal", zap.Error(err))
	}
	_, err = w.Write(resp)
	handler.logger.Error("error while write response", zap.Error(err))
}
