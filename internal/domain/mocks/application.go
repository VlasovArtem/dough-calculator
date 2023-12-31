// Code generated by MockGen. DO NOT EDIT.
// Source: application.go
//
// Generated by this command:
//
//	mockgen -source=application.go -destination=mocks/application.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	config "dough-calculator/internal/config"
	domain "dough-calculator/internal/domain"
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockApplicationInitializer is a mock of ApplicationInitializer interface.
type MockApplicationInitializer struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationInitializerMockRecorder
}

// MockApplicationInitializerMockRecorder is the mock recorder for MockApplicationInitializer.
type MockApplicationInitializerMockRecorder struct {
	mock *MockApplicationInitializer
}

// NewMockApplicationInitializer creates a new mock instance.
func NewMockApplicationInitializer(ctrl *gomock.Controller) *MockApplicationInitializer {
	mock := &MockApplicationInitializer{ctrl: ctrl}
	mock.recorder = &MockApplicationInitializerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationInitializer) EXPECT() *MockApplicationInitializerMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockApplicationInitializer) Initialize() (domain.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(domain.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Initialize indicates an expected call of Initialize.
func (mr *MockApplicationInitializerMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockApplicationInitializer)(nil).Initialize))
}

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Config mocks base method.
func (m *MockApplication) Config() config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(config.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockApplicationMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockApplication)(nil).Config))
}

// Server mocks base method.
func (m *MockApplication) Server() *http.Server {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Server")
	ret0, _ := ret[0].(*http.Server)
	return ret0
}

// Server indicates an expected call of Server.
func (mr *MockApplicationMockRecorder) Server() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Server", reflect.TypeOf((*MockApplication)(nil).Server))
}
