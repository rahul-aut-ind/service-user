// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Context is an autogenerated mock type for the Context type
type Context struct {
	mock.Mock
}

// Deadline provides a mock function with given fields:
func (_m *Context) Deadline() (time.Time, bool) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Deadline")
	}

	var r0 time.Time
	var r1 bool
	if rf, ok := ret.Get(0).(func() (time.Time, bool)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Done provides a mock function with given fields:
func (_m *Context) Done() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Err provides a mock function with given fields:
func (_m *Context) Err() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Err")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHeader provides a mock function with given fields: key
func (_m *Context) GetHeader(key string) string {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for GetHeader")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// JSON provides a mock function with given fields: code, obj
func (_m *Context) JSON(code int, obj interface{}) {
	_m.Called(code, obj)
}

// Param provides a mock function with given fields: key
func (_m *Context) Param(key string) string {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Param")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Query provides a mock function with given fields: key
func (_m *Context) Query(key string) string {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ShouldBindJSON provides a mock function with given fields: obj
func (_m *Context) ShouldBindJSON(obj interface{}) error {
	ret := _m.Called(obj)

	if len(ret) == 0 {
		panic("no return value specified for ShouldBindJSON")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Value provides a mock function with given fields: key
func (_m *Context) Value(key any) any {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Value")
	}

	var r0 any
	if rf, ok := ret.Get(0).(func(any) any); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(any)
		}
	}

	return r0
}

// NewContext creates a new instance of Context. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContext(t interface {
	mock.TestingT
	Cleanup(func())
}) *Context {
	mock := &Context{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}