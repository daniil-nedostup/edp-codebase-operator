// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/epam/edp-codebase-operator/v2/api/v1"
)

// MockCodebaseHandler is an autogenerated mock type for the CodebaseHandler type
type MockCodebaseHandler struct {
	mock.Mock
}

type MockCodebaseHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCodebaseHandler) EXPECT() *MockCodebaseHandler_Expecter {
	return &MockCodebaseHandler_Expecter{mock: &_m.Mock}
}

// ServeRequest provides a mock function with given fields: _a0, _a1
func (_m *MockCodebaseHandler) ServeRequest(_a0 context.Context, _a1 *v1.Codebase) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ServeRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Codebase) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCodebaseHandler_ServeRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServeRequest'
type MockCodebaseHandler_ServeRequest_Call struct {
	*mock.Call
}

// ServeRequest is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *v1.Codebase
func (_e *MockCodebaseHandler_Expecter) ServeRequest(_a0 interface{}, _a1 interface{}) *MockCodebaseHandler_ServeRequest_Call {
	return &MockCodebaseHandler_ServeRequest_Call{Call: _e.mock.On("ServeRequest", _a0, _a1)}
}

func (_c *MockCodebaseHandler_ServeRequest_Call) Run(run func(_a0 context.Context, _a1 *v1.Codebase)) *MockCodebaseHandler_ServeRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Codebase))
	})
	return _c
}

func (_c *MockCodebaseHandler_ServeRequest_Call) Return(_a0 error) *MockCodebaseHandler_ServeRequest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCodebaseHandler_ServeRequest_Call) RunAndReturn(run func(context.Context, *v1.Codebase) error) *MockCodebaseHandler_ServeRequest_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCodebaseHandler creates a new instance of MockCodebaseHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCodebaseHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCodebaseHandler {
	mock := &MockCodebaseHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
