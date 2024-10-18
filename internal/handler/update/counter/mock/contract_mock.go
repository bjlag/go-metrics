// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	logger "github.com/bjlag/go-metrics/internal/logger"
	gomock "github.com/golang/mock/gomock"
)

// Mockrepo is a mock of repo interface.
type Mockrepo struct {
	ctrl     *gomock.Controller
	recorder *MockrepoMockRecorder
}

// MockrepoMockRecorder is the mock recorder for Mockrepo.
type MockrepoMockRecorder struct {
	mock *Mockrepo
}

// NewMockrepo creates a new mock instance.
func NewMockrepo(ctrl *gomock.Controller) *Mockrepo {
	mock := &Mockrepo{ctrl: ctrl}
	mock.recorder = &MockrepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepo) EXPECT() *MockrepoMockRecorder {
	return m.recorder
}

// AddCounter mocks base method.
func (m *Mockrepo) AddCounter(ctx context.Context, name string, value int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCounter", ctx, name, value)
}

// AddCounter indicates an expected call of AddCounter.
func (mr *MockrepoMockRecorder) AddCounter(ctx, name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCounter", reflect.TypeOf((*Mockrepo)(nil).AddCounter), ctx, name, value)
}

// Mockbackup is a mock of backup interface.
type Mockbackup struct {
	ctrl     *gomock.Controller
	recorder *MockbackupMockRecorder
}

// MockbackupMockRecorder is the mock recorder for Mockbackup.
type MockbackupMockRecorder struct {
	mock *Mockbackup
}

// NewMockbackup creates a new mock instance.
func NewMockbackup(ctrl *gomock.Controller) *Mockbackup {
	mock := &Mockbackup{ctrl: ctrl}
	mock.recorder = &MockbackupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockbackup) EXPECT() *MockbackupMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockbackup) Create(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockbackupMockRecorder) Create(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockbackup)(nil).Create), ctx)
}

// Mocklog is a mock of log interface.
type Mocklog struct {
	ctrl     *gomock.Controller
	recorder *MocklogMockRecorder
}

// MocklogMockRecorder is the mock recorder for Mocklog.
type MocklogMockRecorder struct {
	mock *Mocklog
}

// NewMocklog creates a new mock instance.
func NewMocklog(ctrl *gomock.Controller) *Mocklog {
	mock := &Mocklog{ctrl: ctrl}
	mock.recorder = &MocklogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklog) EXPECT() *MocklogMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Mocklog) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MocklogMockRecorder) Error(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklog)(nil).Error), msg)
}

// Info mocks base method.
func (m *Mocklog) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MocklogMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklog)(nil).Info), msg)
}

// WithField mocks base method.
func (m *Mocklog) WithField(key string, value interface{}) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithField", key, value)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithField indicates an expected call of WithField.
func (mr *MocklogMockRecorder) WithField(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithField", reflect.TypeOf((*Mocklog)(nil).WithField), key, value)
}
