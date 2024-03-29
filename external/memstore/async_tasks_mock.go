// Code generated by MockGen. DO NOT EDIT.
// Source: ./external/memstore/async_tasks.go

// Package memstore is a generated GoMock package.
package memstore

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskPublisher is a mock of TaskPublisher interface.
type MockTaskPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockTaskPublisherMockRecorder
}

// MockTaskPublisherMockRecorder is the mock recorder for MockTaskPublisher.
type MockTaskPublisherMockRecorder struct {
	mock *MockTaskPublisher
}

// NewMockTaskPublisher creates a new mock instance.
func NewMockTaskPublisher(ctrl *gomock.Controller) *MockTaskPublisher {
	mock := &MockTaskPublisher{ctrl: ctrl}
	mock.recorder = &MockTaskPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskPublisher) EXPECT() *MockTaskPublisherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTaskPublisher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTaskPublisherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTaskPublisher)(nil).Close))
}

// PublishAsyncTask mocks base method.
func (m *MockTaskPublisher) PublishAsyncTask(task AsyncTask, opt TaskOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishAsyncTask", task, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishAsyncTask indicates an expected call of PublishAsyncTask.
func (mr *MockTaskPublisherMockRecorder) PublishAsyncTask(task, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishAsyncTask", reflect.TypeOf((*MockTaskPublisher)(nil).PublishAsyncTask), task, opt)
}

// MockTaskServer is a mock of TaskServer interface.
type MockTaskServer struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServerMockRecorder
}

// MockTaskServerMockRecorder is the mock recorder for MockTaskServer.
type MockTaskServerMockRecorder struct {
	mock *MockTaskServer
}

// NewMockTaskServer creates a new mock instance.
func NewMockTaskServer(ctrl *gomock.Controller) *MockTaskServer {
	mock := &MockTaskServer{ctrl: ctrl}
	mock.recorder = &MockTaskServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskServer) EXPECT() *MockTaskServerMockRecorder {
	return m.recorder
}

// RegisterTaskWorker mocks base method.
func (m *MockTaskServer) RegisterTaskWorker(taskType TaskType, tw ProcessTaskFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterTaskWorker", taskType, tw)
}

// RegisterTaskWorker indicates an expected call of RegisterTaskWorker.
func (mr *MockTaskServerMockRecorder) RegisterTaskWorker(taskType, tw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTaskWorker", reflect.TypeOf((*MockTaskServer)(nil).RegisterTaskWorker), taskType, tw)
}

// Run mocks base method.
func (m *MockTaskServer) Run() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockTaskServerMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockTaskServer)(nil).Run))
}

// StopAndShutdown mocks base method.
func (m *MockTaskServer) StopAndShutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopAndShutdown")
}

// StopAndShutdown indicates an expected call of StopAndShutdown.
func (mr *MockTaskServerMockRecorder) StopAndShutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopAndShutdown", reflect.TypeOf((*MockTaskServer)(nil).StopAndShutdown))
}

// MockTaskIdentifier is a mock of TaskIdentifier interface.
type MockTaskIdentifier struct {
	ctrl     *gomock.Controller
	recorder *MockTaskIdentifierMockRecorder
}

// MockTaskIdentifierMockRecorder is the mock recorder for MockTaskIdentifier.
type MockTaskIdentifierMockRecorder struct {
	mock *MockTaskIdentifier
}

// NewMockTaskIdentifier creates a new mock instance.
func NewMockTaskIdentifier(ctrl *gomock.Controller) *MockTaskIdentifier {
	mock := &MockTaskIdentifier{ctrl: ctrl}
	mock.recorder = &MockTaskIdentifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskIdentifier) EXPECT() *MockTaskIdentifierMockRecorder {
	return m.recorder
}

// GetPayload mocks base method.
func (m *MockTaskIdentifier) GetPayload() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayload")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetPayload indicates an expected call of GetPayload.
func (mr *MockTaskIdentifierMockRecorder) GetPayload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayload", reflect.TypeOf((*MockTaskIdentifier)(nil).GetPayload))
}

// GetTaskType mocks base method.
func (m *MockTaskIdentifier) GetTaskType() TaskType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskType")
	ret0, _ := ret[0].(TaskType)
	return ret0
}

// GetTaskType indicates an expected call of GetTaskType.
func (mr *MockTaskIdentifierMockRecorder) GetTaskType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskType", reflect.TypeOf((*MockTaskIdentifier)(nil).GetTaskType))
}

// TaskID mocks base method.
func (m *MockTaskIdentifier) TaskID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TaskID")
	ret0, _ := ret[0].(string)
	return ret0
}

// TaskID indicates an expected call of TaskID.
func (mr *MockTaskIdentifierMockRecorder) TaskID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TaskID", reflect.TypeOf((*MockTaskIdentifier)(nil).TaskID))
}
