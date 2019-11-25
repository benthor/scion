// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/lib/sciond (interfaces: Service,Connector)

// Package mock_sciond is a generated GoMock package.
package mock_sciond

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/go/lib/addr"
	common "github.com/scionproto/scion/go/lib/common"
	path_mgmt "github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
	sciond "github.com/scionproto/scion/go/lib/sciond"
	snet "github.com/scionproto/scion/go/lib/snet"
	proto "github.com/scionproto/scion/go/proto"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockService) Connect(arg0 context.Context) (sciond.Connector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(sciond.Connector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect
func (mr *MockServiceMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockService)(nil).Connect), arg0)
}

// MockConnector is a mock of Connector interface
type MockConnector struct {
	ctrl     *gomock.Controller
	recorder *MockConnectorMockRecorder
}

// MockConnectorMockRecorder is the mock recorder for MockConnector
type MockConnectorMockRecorder struct {
	mock *MockConnector
}

// NewMockConnector creates a new mock instance
func NewMockConnector(ctrl *gomock.Controller) *MockConnector {
	mock := &MockConnector{ctrl: ctrl}
	mock.recorder = &MockConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnector) EXPECT() *MockConnectorMockRecorder {
	return m.recorder
}

// ASInfo mocks base method
func (m *MockConnector) ASInfo(arg0 context.Context, arg1 addr.IA) (*sciond.ASInfoReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ASInfo", arg0, arg1)
	ret0, _ := ret[0].(*sciond.ASInfoReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ASInfo indicates an expected call of ASInfo
func (mr *MockConnectorMockRecorder) ASInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ASInfo", reflect.TypeOf((*MockConnector)(nil).ASInfo), arg0, arg1)
}

// Close mocks base method
func (m *MockConnector) Close(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockConnectorMockRecorder) Close(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnector)(nil).Close), arg0)
}

// IFInfo mocks base method
func (m *MockConnector) IFInfo(arg0 context.Context, arg1 []common.IFIDType) (*sciond.IFInfoReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IFInfo", arg0, arg1)
	ret0, _ := ret[0].(*sciond.IFInfoReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IFInfo indicates an expected call of IFInfo
func (mr *MockConnectorMockRecorder) IFInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IFInfo", reflect.TypeOf((*MockConnector)(nil).IFInfo), arg0, arg1)
}

// Paths mocks base method
func (m *MockConnector) Paths(arg0 context.Context, arg1, arg2 addr.IA, arg3 uint16, arg4 sciond.PathReqFlags) ([]snet.Path, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Paths", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]snet.Path)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Paths indicates an expected call of Paths
func (mr *MockConnectorMockRecorder) Paths(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Paths", reflect.TypeOf((*MockConnector)(nil).Paths), arg0, arg1, arg2, arg3, arg4)
}

// RevNotification mocks base method
func (m *MockConnector) RevNotification(arg0 context.Context, arg1 *path_mgmt.SignedRevInfo) (*sciond.RevReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevNotification", arg0, arg1)
	ret0, _ := ret[0].(*sciond.RevReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevNotification indicates an expected call of RevNotification
func (mr *MockConnectorMockRecorder) RevNotification(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevNotification", reflect.TypeOf((*MockConnector)(nil).RevNotification), arg0, arg1)
}

// RevNotificationFromRaw mocks base method
func (m *MockConnector) RevNotificationFromRaw(arg0 context.Context, arg1 []byte) (*sciond.RevReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevNotificationFromRaw", arg0, arg1)
	ret0, _ := ret[0].(*sciond.RevReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevNotificationFromRaw indicates an expected call of RevNotificationFromRaw
func (mr *MockConnectorMockRecorder) RevNotificationFromRaw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevNotificationFromRaw", reflect.TypeOf((*MockConnector)(nil).RevNotificationFromRaw), arg0, arg1)
}

// SVCInfo mocks base method
func (m *MockConnector) SVCInfo(arg0 context.Context, arg1 []proto.ServiceType) (*sciond.ServiceInfoReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SVCInfo", arg0, arg1)
	ret0, _ := ret[0].(*sciond.ServiceInfoReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SVCInfo indicates an expected call of SVCInfo
func (mr *MockConnectorMockRecorder) SVCInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SVCInfo", reflect.TypeOf((*MockConnector)(nil).SVCInfo), arg0, arg1)
}
