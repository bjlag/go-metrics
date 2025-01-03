// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bjlag/go-metrics/internal/logger (interfaces: Logger)

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
func (m *MockLogger) Debug(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", arg0)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), arg0)
}

// Error mocks base method.
func (m *MockLogger) Error(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), arg0)
}

// Info mocks base method.
func (m *MockLogger) Info(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", arg0)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), arg0)
}

// WithError mocks base method.
func (m *MockLogger) WithError(arg0 error) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithError", arg0)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithError indicates an expected call of WithError.
func (mr *MockLoggerMockRecorder) WithError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithError", reflect.TypeOf((*MockLogger)(nil).WithError), arg0)
}

// WithField mocks base method.
func (m *MockLogger) WithField(arg0 string, arg1 interface{}) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithField", arg0, arg1)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithField indicates an expected call of WithField.
func (mr *MockLoggerMockRecorder) WithField(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithField", reflect.TypeOf((*MockLogger)(nil).WithField), arg0, arg1)
}
