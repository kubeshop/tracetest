// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeshop/tracetest/server/go (interfaces: TestDB)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	openapi "github.com/kubeshop/tracetest/server/go"
)

// MockTestDB is a mock of TestDB interface.
type MockTestDB struct {
	ctrl     *gomock.Controller
	recorder *MockTestDBMockRecorder
}

// MockTestDBMockRecorder is the mock recorder for MockTestDB.
type MockTestDBMockRecorder struct {
	mock *MockTestDB
}

// NewMockTestDB creates a new mock instance.
func NewMockTestDB(ctrl *gomock.Controller) *MockTestDB {
	mock := &MockTestDB{ctrl: ctrl}
	mock.recorder = &MockTestDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTestDB) EXPECT() *MockTestDBMockRecorder {
	return m.recorder
}

// CreateAssertion mocks base method.
func (m *MockTestDB) CreateAssertion(arg0 context.Context, arg1 string, arg2 *openapi.Assertion) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAssertion", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAssertion indicates an expected call of CreateAssertion.
func (mr *MockTestDBMockRecorder) CreateAssertion(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAssertion", reflect.TypeOf((*MockTestDB)(nil).CreateAssertion), arg0, arg1, arg2)
}

// CreateResult mocks base method.
func (m *MockTestDB) CreateResult(arg0 context.Context, arg1 string, arg2 *openapi.TestRunResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateResult", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateResult indicates an expected call of CreateResult.
func (mr *MockTestDBMockRecorder) CreateResult(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateResult", reflect.TypeOf((*MockTestDB)(nil).CreateResult), arg0, arg1, arg2)
}

// CreateTest mocks base method.
func (m *MockTestDB) CreateTest(arg0 context.Context, arg1 *openapi.Test) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTest", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTest indicates an expected call of CreateTest.
func (mr *MockTestDBMockRecorder) CreateTest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTest", reflect.TypeOf((*MockTestDB)(nil).CreateTest), arg0, arg1)
}

// GetAssertion mocks base method.
func (m *MockTestDB) GetAssertion(arg0 context.Context, arg1 string) (*openapi.Assertion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssertion", arg0, arg1)
	ret0, _ := ret[0].(*openapi.Assertion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssertion indicates an expected call of GetAssertion.
func (mr *MockTestDBMockRecorder) GetAssertion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssertion", reflect.TypeOf((*MockTestDB)(nil).GetAssertion), arg0, arg1)
}

// GetAssertionsByTestID mocks base method.
func (m *MockTestDB) GetAssertionsByTestID(arg0 context.Context, arg1 string) ([]openapi.Assertion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssertionsByTestID", arg0, arg1)
	ret0, _ := ret[0].([]openapi.Assertion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssertionsByTestID indicates an expected call of GetAssertionsByTestID.
func (mr *MockTestDBMockRecorder) GetAssertionsByTestID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssertionsByTestID", reflect.TypeOf((*MockTestDB)(nil).GetAssertionsByTestID), arg0, arg1)
}

// GetResult mocks base method.
func (m *MockTestDB) GetResult(arg0 context.Context, arg1 string) (*openapi.TestRunResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResult", arg0, arg1)
	ret0, _ := ret[0].(*openapi.TestRunResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResult indicates an expected call of GetResult.
func (mr *MockTestDBMockRecorder) GetResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResult", reflect.TypeOf((*MockTestDB)(nil).GetResult), arg0, arg1)
}

// GetResultsByTestID mocks base method.
func (m *MockTestDB) GetResultsByTestID(arg0 context.Context, arg1 string) ([]openapi.TestRunResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResultsByTestID", arg0, arg1)
	ret0, _ := ret[0].([]openapi.TestRunResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResultsByTestID indicates an expected call of GetResultsByTestID.
func (mr *MockTestDBMockRecorder) GetResultsByTestID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResultsByTestID", reflect.TypeOf((*MockTestDB)(nil).GetResultsByTestID), arg0, arg1)
}

// GetResultsByTraceID mocks base method.
func (m *MockTestDB) GetResultsByTraceID(arg0 context.Context, arg1, arg2 string) (openapi.TestRunResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResultsByTraceID", arg0, arg1, arg2)
	ret0, _ := ret[0].(openapi.TestRunResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResultsByTraceID indicates an expected call of GetResultsByTraceID.
func (mr *MockTestDBMockRecorder) GetResultsByTraceID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResultsByTraceID", reflect.TypeOf((*MockTestDB)(nil).GetResultsByTraceID), arg0, arg1, arg2)
}

// GetTest mocks base method.
func (m *MockTestDB) GetTest(arg0 context.Context, arg1 string) (*openapi.Test, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTest", arg0, arg1)
	ret0, _ := ret[0].(*openapi.Test)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTest indicates an expected call of GetTest.
func (mr *MockTestDBMockRecorder) GetTest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTest", reflect.TypeOf((*MockTestDB)(nil).GetTest), arg0, arg1)
}

// GetTests mocks base method.
func (m *MockTestDB) GetTests(arg0 context.Context) ([]openapi.Test, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTests", arg0)
	ret0, _ := ret[0].([]openapi.Test)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTests indicates an expected call of GetTests.
func (mr *MockTestDBMockRecorder) GetTests(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTests", reflect.TypeOf((*MockTestDB)(nil).GetTests), arg0)
}

// UpdateResult mocks base method.
func (m *MockTestDB) UpdateResult(arg0 context.Context, arg1 *openapi.TestRunResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateResult", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResult indicates an expected call of UpdateResult.
func (mr *MockTestDBMockRecorder) UpdateResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResult", reflect.TypeOf((*MockTestDB)(nil).UpdateResult), arg0, arg1)
}

// UpdateTest mocks base method.
func (m *MockTestDB) UpdateTest(arg0 context.Context, arg1 *openapi.Test) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTest indicates an expected call of UpdateTest.
func (mr *MockTestDBMockRecorder) UpdateTest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTest", reflect.TypeOf((*MockTestDB)(nil).UpdateTest), arg0, arg1)
}
