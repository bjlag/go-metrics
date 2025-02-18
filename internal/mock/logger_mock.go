// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	logger "github.com/bjlag/go-metrics/internal/logger"
	gomock "github.com/golang/mock/gomock"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", msg)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), msg)
}

// Error mocks base method.
func (m *MockLogger) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), msg)
}

// Info mocks base method.
func (m *MockLogger) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), msg)
}

// WithError mocks base method.
func (m *MockLogger) WithError(err error) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithError", err)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithError indicates an expected call of WithError.
func (mr *MockLoggerMockRecorder) WithError(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithError", reflect.TypeOf((*MockLogger)(nil).WithError), err)
}

// WithField mocks base method.
func (m *MockLogger) WithField(key string, value interface{}) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithField", key, value)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithField indicates an expected call of WithField.
func (mr *MockLoggerMockRecorder) WithField(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithField", reflect.TypeOf((*MockLogger)(nil).WithField), key, value)
}
