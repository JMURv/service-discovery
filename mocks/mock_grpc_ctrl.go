// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/hdl/grpc/grpc.go
//
// Generated by this command:
//
//	mockgen -source=./internal/hdl/grpc/grpc.go -destination=mocks/mock_grpc_ctrl.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCtrl is a mock of Ctrl interface.
type MockCtrl struct {
	ctrl     *gomock.Controller
	recorder *MockCtrlMockRecorder
}

// MockCtrlMockRecorder is the mock recorder for MockCtrl.
type MockCtrlMockRecorder struct {
	mock *MockCtrl
}

// NewMockCtrl creates a new mock instance.
func NewMockCtrl(ctrl *gomock.Controller) *MockCtrl {
	mock := &MockCtrl{ctrl: ctrl}
	mock.recorder = &MockCtrlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCtrl) EXPECT() *MockCtrlMockRecorder {
	return m.recorder
}

// Deregister mocks base method.
func (m *MockCtrl) Deregister(ctx context.Context, name, addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deregister", ctx, name, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deregister indicates an expected call of Deregister.
func (mr *MockCtrlMockRecorder) Deregister(ctx, name, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deregister", reflect.TypeOf((*MockCtrl)(nil).Deregister), ctx, name, addr)
}

// FindServiceByName mocks base method.
func (m *MockCtrl) FindServiceByName(ctx context.Context, name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindServiceByName", ctx, name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindServiceByName indicates an expected call of FindServiceByName.
func (mr *MockCtrlMockRecorder) FindServiceByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindServiceByName", reflect.TypeOf((*MockCtrl)(nil).FindServiceByName), ctx, name)
}

// ListAddrs mocks base method.
func (m *MockCtrl) ListAddrs(ctx context.Context, name string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAddrs", ctx, name)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAddrs indicates an expected call of ListAddrs.
func (mr *MockCtrlMockRecorder) ListAddrs(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAddrs", reflect.TypeOf((*MockCtrl)(nil).ListAddrs), ctx, name)
}

// ListServices mocks base method.
func (m *MockCtrl) ListServices(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServices", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServices indicates an expected call of ListServices.
func (mr *MockCtrlMockRecorder) ListServices(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServices", reflect.TypeOf((*MockCtrl)(nil).ListServices), ctx)
}

// Register mocks base method.
func (m *MockCtrl) Register(ctx context.Context, name, addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, name, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockCtrlMockRecorder) Register(ctx, name, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCtrl)(nil).Register), ctx, name, addr)
}
