// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/lib/snet/internal/ctxmonitor (interfaces: Monitor)

// Package mock_ctxmonitor is a generated GoMock package.
package mock_ctxmonitor

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockMonitor is a mock of Monitor interface
type MockMonitor struct {
	ctrl     *gomock.Controller
	recorder *MockMonitorMockRecorder
}

// MockMonitorMockRecorder is the mock recorder for MockMonitor
type MockMonitorMockRecorder struct {
	mock *MockMonitor
}

// NewMockMonitor creates a new mock instance
func NewMockMonitor(ctrl *gomock.Controller) *MockMonitor {
	mock := &MockMonitor{ctrl: ctrl}
	mock.recorder = &MockMonitorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMonitor) EXPECT() *MockMonitorMockRecorder {
	return m.recorder
}

// Count mocks base method
func (m *MockMonitor) Count() int {
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	return ret0
}

// Count indicates an expected call of Count
func (mr *MockMonitorMockRecorder) Count() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockMonitor)(nil).Count))
}

// SetDeadline mocks base method
func (m *MockMonitor) SetDeadline(arg0 time.Time) {
	m.ctrl.Call(m, "SetDeadline", arg0)
}

// SetDeadline indicates an expected call of SetDeadline
func (mr *MockMonitorMockRecorder) SetDeadline(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockMonitor)(nil).SetDeadline), arg0)
}

// WithDeadline mocks base method
func (m *MockMonitor) WithDeadline(arg0 context.Context, arg1 time.Time) (context.Context, context.CancelFunc) {
	ret := m.ctrl.Call(m, "WithDeadline", arg0, arg1)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(context.CancelFunc)
	return ret0, ret1
}

// WithDeadline indicates an expected call of WithDeadline
func (mr *MockMonitorMockRecorder) WithDeadline(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithDeadline", reflect.TypeOf((*MockMonitor)(nil).WithDeadline), arg0, arg1)
}

// WithTimeout mocks base method
func (m *MockMonitor) WithTimeout(arg0 context.Context, arg1 time.Duration) (context.Context, context.CancelFunc) {
	ret := m.ctrl.Call(m, "WithTimeout", arg0, arg1)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(context.CancelFunc)
	return ret0, ret1
}

// WithTimeout indicates an expected call of WithTimeout
func (mr *MockMonitorMockRecorder) WithTimeout(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTimeout", reflect.TypeOf((*MockMonitor)(nil).WithTimeout), arg0, arg1)
}
