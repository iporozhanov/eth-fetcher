// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

type Auth_Expecter struct {
	mock *mock.Mock
}

func (_m *Auth) EXPECT() *Auth_Expecter {
	return &Auth_Expecter{mock: &_m.Mock}
}

// AuthenticateRequest provides a mock function with given fields: r
func (_m *Auth) AuthenticateRequest(r *http.Request) (string, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for AuthenticateRequest")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (string, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) string); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Auth_AuthenticateRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AuthenticateRequest'
type Auth_AuthenticateRequest_Call struct {
	*mock.Call
}

// AuthenticateRequest is a helper method to define mock.On call
//   - r *http.Request
func (_e *Auth_Expecter) AuthenticateRequest(r interface{}) *Auth_AuthenticateRequest_Call {
	return &Auth_AuthenticateRequest_Call{Call: _e.mock.On("AuthenticateRequest", r)}
}

func (_c *Auth_AuthenticateRequest_Call) Run(run func(r *http.Request)) *Auth_AuthenticateRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *Auth_AuthenticateRequest_Call) Return(_a0 string, _a1 error) *Auth_AuthenticateRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Auth_AuthenticateRequest_Call) RunAndReturn(run func(*http.Request) (string, error)) *Auth_AuthenticateRequest_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateToken provides a mock function with given fields: subject
func (_m *Auth) GenerateToken(subject string) (string, error) {
	ret := _m.Called(subject)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(subject)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(subject)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subject)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Auth_GenerateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateToken'
type Auth_GenerateToken_Call struct {
	*mock.Call
}

// GenerateToken is a helper method to define mock.On call
//   - subject string
func (_e *Auth_Expecter) GenerateToken(subject interface{}) *Auth_GenerateToken_Call {
	return &Auth_GenerateToken_Call{Call: _e.mock.On("GenerateToken", subject)}
}

func (_c *Auth_GenerateToken_Call) Run(run func(subject string)) *Auth_GenerateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Auth_GenerateToken_Call) Return(_a0 string, _a1 error) *Auth_GenerateToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Auth_GenerateToken_Call) RunAndReturn(run func(string) (string, error)) *Auth_GenerateToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuth creates a new instance of Auth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuth(t interface {
	mock.TestingT
	Cleanup(func())
}) *Auth {
	mock := &Auth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}