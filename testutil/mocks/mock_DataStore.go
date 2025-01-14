// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	snapmatchai "github.com/trapajim/snapmatch-ai/snapmatchai"
)

// MockDataStore is an autogenerated mock type for the DataStore type
type MockDataStore struct {
	mock.Mock
}

type MockDataStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataStore) EXPECT() *MockDataStore_Expecter {
	return &MockDataStore_Expecter{mock: &_m.Mock}
}

// Query provides a mock function with given fields: ctx, queryString, parameters, target
func (_m *MockDataStore) Query(ctx context.Context, queryString string, parameters map[string]interface{}, target interface{}) error {
	ret := _m.Called(ctx, queryString, parameters, target)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, queryString, parameters, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataStore_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockDataStore_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - queryString string
//   - parameters map[string]interface{}
//   - target interface{}
func (_e *MockDataStore_Expecter) Query(ctx interface{}, queryString interface{}, parameters interface{}, target interface{}) *MockDataStore_Query_Call {
	return &MockDataStore_Query_Call{Call: _e.mock.On("Query", ctx, queryString, parameters, target)}
}

func (_c *MockDataStore_Query_Call) Run(run func(ctx context.Context, queryString string, parameters map[string]interface{}, target interface{})) *MockDataStore_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(map[string]interface{}), args[3].(interface{}))
	})
	return _c
}

func (_c *MockDataStore_Query_Call) Return(_a0 error) *MockDataStore_Query_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataStore_Query_Call) RunAndReturn(run func(context.Context, string, map[string]interface{}, interface{}) error) *MockDataStore_Query_Call {
	_c.Call.Return(run)
	return _c
}

// Schema provides a mock function with given fields: ctx, dataset, tableName
func (_m *MockDataStore) Schema(ctx context.Context, dataset string, tableName string) ([]snapmatchai.DBSchema, error) {
	ret := _m.Called(ctx, dataset, tableName)

	if len(ret) == 0 {
		panic("no return value specified for Schema")
	}

	var r0 []snapmatchai.DBSchema
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]snapmatchai.DBSchema, error)); ok {
		return rf(ctx, dataset, tableName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []snapmatchai.DBSchema); ok {
		r0 = rf(ctx, dataset, tableName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]snapmatchai.DBSchema)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, dataset, tableName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataStore_Schema_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Schema'
type MockDataStore_Schema_Call struct {
	*mock.Call
}

// Schema is a helper method to define mock.On call
//   - ctx context.Context
//   - dataset string
//   - tableName string
func (_e *MockDataStore_Expecter) Schema(ctx interface{}, dataset interface{}, tableName interface{}) *MockDataStore_Schema_Call {
	return &MockDataStore_Schema_Call{Call: _e.mock.On("Schema", ctx, dataset, tableName)}
}

func (_c *MockDataStore_Schema_Call) Run(run func(ctx context.Context, dataset string, tableName string)) *MockDataStore_Schema_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDataStore_Schema_Call) Return(_a0 []snapmatchai.DBSchema, _a1 error) *MockDataStore_Schema_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataStore_Schema_Call) RunAndReturn(run func(context.Context, string, string) ([]snapmatchai.DBSchema, error)) *MockDataStore_Schema_Call {
	_c.Call.Return(run)
	return _c
}

// TableExists provides a mock function with given fields: ctx, dataset, table
func (_m *MockDataStore) TableExists(ctx context.Context, dataset string, table string) error {
	ret := _m.Called(ctx, dataset, table)

	if len(ret) == 0 {
		panic("no return value specified for TableExists")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, dataset, table)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataStore_TableExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TableExists'
type MockDataStore_TableExists_Call struct {
	*mock.Call
}

// TableExists is a helper method to define mock.On call
//   - ctx context.Context
//   - dataset string
//   - table string
func (_e *MockDataStore_Expecter) TableExists(ctx interface{}, dataset interface{}, table interface{}) *MockDataStore_TableExists_Call {
	return &MockDataStore_TableExists_Call{Call: _e.mock.On("TableExists", ctx, dataset, table)}
}

func (_c *MockDataStore_TableExists_Call) Run(run func(ctx context.Context, dataset string, table string)) *MockDataStore_TableExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDataStore_TableExists_Call) Return(_a0 error) *MockDataStore_TableExists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataStore_TableExists_Call) RunAndReturn(run func(context.Context, string, string) error) *MockDataStore_TableExists_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDataStore creates a new instance of MockDataStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDataStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDataStore {
	mock := &MockDataStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
