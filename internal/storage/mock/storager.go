// Code generated by MockGen. DO NOT EDIT.
// Source: storager.go
//
// Generated by this command:
//
//	mockgen -source=storager.go -destination=./mock/storager.go
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// AlreadyExists mocks base method.
func (m *MockStorager) AlreadyExists(ctx context.Context, fullURL string) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AlreadyExists", ctx, fullURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AlreadyExists indicates an expected call of AlreadyExists.
func (mr *MockStoragerMockRecorder) AlreadyExists(ctx, fullURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlreadyExists", reflect.TypeOf((*MockStorager)(nil).AlreadyExists), ctx, fullURL)
}

// Close mocks base method.
func (m *MockStorager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoragerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorager)(nil).Close))
}

// CreateShortURL mocks base method.
func (m *MockStorager) CreateShortURL(ctx context.Context, fullURL, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortURL", ctx, fullURL, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateShortURL indicates an expected call of CreateShortURL.
func (mr *MockStoragerMockRecorder) CreateShortURL(ctx, fullURL, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortURL", reflect.TypeOf((*MockStorager)(nil).CreateShortURL), ctx, fullURL, token)
}

// GetFullURL mocks base method.
func (m *MockStorager) GetFullURL(ctx context.Context, token string) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFullURL", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFullURL indicates an expected call of GetFullURL.
func (mr *MockStoragerMockRecorder) GetFullURL(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullURL", reflect.TypeOf((*MockStorager)(nil).GetFullURL), ctx, token)
}
