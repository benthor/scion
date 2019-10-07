// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/beacon_srv/internal/ifstate (interfaces: RevInserter)

// Package mock_ifstate is a generated GoMock package.
package mock_ifstate

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	path_mgmt "github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
)

// MockRevInserter is a mock of RevInserter interface
type MockRevInserter struct {
	ctrl     *gomock.Controller
	recorder *MockRevInserterMockRecorder
}

// MockRevInserterMockRecorder is the mock recorder for MockRevInserter
type MockRevInserterMockRecorder struct {
	mock *MockRevInserter
}

// NewMockRevInserter creates a new mock instance
func NewMockRevInserter(ctrl *gomock.Controller) *MockRevInserter {
	mock := &MockRevInserter{ctrl: ctrl}
	mock.recorder = &MockRevInserterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRevInserter) EXPECT() *MockRevInserterMockRecorder {
	return m.recorder
}

// InsertRevocations mocks base method
func (m *MockRevInserter) InsertRevocations(arg0 context.Context, arg1 ...*path_mgmt.SignedRevInfo) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertRevocations", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertRevocations indicates an expected call of InsertRevocations
func (mr *MockRevInserterMockRecorder) InsertRevocations(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertRevocations", reflect.TypeOf((*MockRevInserter)(nil).InsertRevocations), varargs...)
}
