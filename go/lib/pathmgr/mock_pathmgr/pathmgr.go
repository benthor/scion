// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/lib/pathmgr (interfaces: Querier,Resolver)

// Package mock_pathmgr is a generated GoMock package.
package mock_pathmgr

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/go/lib/addr"
	common "github.com/scionproto/scion/go/lib/common"
	path_mgmt "github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
	pathmgr "github.com/scionproto/scion/go/lib/pathmgr"
	pathpol "github.com/scionproto/scion/go/lib/pathpol"
	sciond "github.com/scionproto/scion/go/lib/sciond"
	spathmeta "github.com/scionproto/scion/go/lib/spath/spathmeta"
	reflect "reflect"
)

// MockQuerier is a mock of Querier interface
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// Query mocks base method
func (m *MockQuerier) Query(arg0 context.Context, arg1, arg2 addr.IA, arg3 sciond.PathReqFlags) spathmeta.AppPathSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(spathmeta.AppPathSet)
	return ret0
}

// Query indicates an expected call of Query
func (mr *MockQuerierMockRecorder) Query(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockQuerier)(nil).Query), arg0, arg1, arg2, arg3)
}

// MockResolver is a mock of Resolver interface
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// Query mocks base method
func (m *MockResolver) Query(arg0 context.Context, arg1, arg2 addr.IA, arg3 sciond.PathReqFlags) spathmeta.AppPathSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(spathmeta.AppPathSet)
	return ret0
}

// Query indicates an expected call of Query
func (mr *MockResolverMockRecorder) Query(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockResolver)(nil).Query), arg0, arg1, arg2, arg3)
}

// QueryFilter mocks base method
func (m *MockResolver) QueryFilter(arg0 context.Context, arg1, arg2 addr.IA, arg3 *pathpol.Policy) spathmeta.AppPathSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryFilter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(spathmeta.AppPathSet)
	return ret0
}

// QueryFilter indicates an expected call of QueryFilter
func (mr *MockResolverMockRecorder) QueryFilter(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryFilter", reflect.TypeOf((*MockResolver)(nil).QueryFilter), arg0, arg1, arg2, arg3)
}

// Revoke mocks base method
func (m *MockResolver) Revoke(arg0 context.Context, arg1 *path_mgmt.SignedRevInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Revoke", arg0, arg1)
}

// Revoke indicates an expected call of Revoke
func (mr *MockResolverMockRecorder) Revoke(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockResolver)(nil).Revoke), arg0, arg1)
}

// RevokeRaw mocks base method
func (m *MockResolver) RevokeRaw(arg0 context.Context, arg1 common.RawBytes) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RevokeRaw", arg0, arg1)
}

// RevokeRaw indicates an expected call of RevokeRaw
func (mr *MockResolverMockRecorder) RevokeRaw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeRaw", reflect.TypeOf((*MockResolver)(nil).RevokeRaw), arg0, arg1)
}

// Sciond mocks base method
func (m *MockResolver) Sciond() sciond.Connector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sciond")
	ret0, _ := ret[0].(sciond.Connector)
	return ret0
}

// Sciond indicates an expected call of Sciond
func (mr *MockResolverMockRecorder) Sciond() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sciond", reflect.TypeOf((*MockResolver)(nil).Sciond))
}

// Watch mocks base method
func (m *MockResolver) Watch(arg0 context.Context, arg1, arg2 addr.IA) (*pathmgr.SyncPaths, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1, arg2)
	ret0, _ := ret[0].(*pathmgr.SyncPaths)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockResolverMockRecorder) Watch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockResolver)(nil).Watch), arg0, arg1, arg2)
}

// WatchCount mocks base method
func (m *MockResolver) WatchCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// WatchCount indicates an expected call of WatchCount
func (mr *MockResolverMockRecorder) WatchCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCount", reflect.TypeOf((*MockResolver)(nil).WatchCount))
}

// WatchFilter mocks base method
func (m *MockResolver) WatchFilter(arg0 context.Context, arg1, arg2 addr.IA, arg3 *pathpol.Policy) (*pathmgr.SyncPaths, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchFilter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*pathmgr.SyncPaths)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchFilter indicates an expected call of WatchFilter
func (mr *MockResolverMockRecorder) WatchFilter(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchFilter", reflect.TypeOf((*MockResolver)(nil).WatchFilter), arg0, arg1, arg2, arg3)
}
