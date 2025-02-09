// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	snapmatchai "github.com/trapajim/snapmatch-ai/snapmatchai"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository[T snapmatchai.Entity] struct {
	mock.Mock
}

type MockRepository_Expecter[T snapmatchai.Entity] struct {
	mock *mock.Mock
}

func (_m *MockRepository[T]) EXPECT() *MockRepository_Expecter[T] {
	return &MockRepository_Expecter[T]{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, entity
func (_m *MockRepository[T]) Create(ctx context.Context, entity T) error {
	ret := _m.Called(ctx, entity)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, T) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepository_Create_Call[T snapmatchai.Entity] struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - entity T
func (_e *MockRepository_Expecter[T]) Create(ctx interface{}, entity interface{}) *MockRepository_Create_Call[T] {
	return &MockRepository_Create_Call[T]{Call: _e.mock.On("Create", ctx, entity)}
}

func (_c *MockRepository_Create_Call[T]) Run(run func(ctx context.Context, entity T)) *MockRepository_Create_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(T))
	})
	return _c
}

func (_c *MockRepository_Create_Call[T]) Return(_a0 error) *MockRepository_Create_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Create_Call[T]) RunAndReturn(run func(context.Context, T) error) *MockRepository_Create_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockRepository[T]) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRepository_Delete_Call[T snapmatchai.Entity] struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockRepository_Expecter[T]) Delete(ctx interface{}, id interface{}) *MockRepository_Delete_Call[T] {
	return &MockRepository_Delete_Call[T]{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockRepository_Delete_Call[T]) Run(run func(ctx context.Context, id string)) *MockRepository_Delete_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockRepository_Delete_Call[T]) Return(_a0 error) *MockRepository_Delete_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Delete_Call[T]) RunAndReturn(run func(context.Context, string) error) *MockRepository_Delete_Call[T] {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, filters
func (_m *MockRepository[T]) List(ctx context.Context, filters map[string]interface{}) ([]T, error) {
	ret := _m.Called(ctx, filters)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) ([]T, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) []T); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockRepository_List_Call[T snapmatchai.Entity] struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - filters map[string]interface{}
func (_e *MockRepository_Expecter[T]) List(ctx interface{}, filters interface{}) *MockRepository_List_Call[T] {
	return &MockRepository_List_Call[T]{Call: _e.mock.On("List", ctx, filters)}
}

func (_c *MockRepository_List_Call[T]) Run(run func(ctx context.Context, filters map[string]interface{})) *MockRepository_List_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *MockRepository_List_Call[T]) Return(_a0 []T, _a1 error) *MockRepository_List_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_List_Call[T]) RunAndReturn(run func(context.Context, map[string]interface{}) ([]T, error)) *MockRepository_List_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx, id
func (_m *MockRepository[T]) Read(ctx context.Context, id string) (T, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (T, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) T); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(T)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockRepository_Read_Call[T snapmatchai.Entity] struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockRepository_Expecter[T]) Read(ctx interface{}, id interface{}) *MockRepository_Read_Call[T] {
	return &MockRepository_Read_Call[T]{Call: _e.mock.On("Read", ctx, id)}
}

func (_c *MockRepository_Read_Call[T]) Run(run func(ctx context.Context, id string)) *MockRepository_Read_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockRepository_Read_Call[T]) Return(_a0 T, _a1 error) *MockRepository_Read_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Read_Call[T]) RunAndReturn(run func(context.Context, string) (T, error)) *MockRepository_Read_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, entity
func (_m *MockRepository[T]) Update(ctx context.Context, entity T) error {
	ret := _m.Called(ctx, entity)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, T) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockRepository_Update_Call[T snapmatchai.Entity] struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - entity T
func (_e *MockRepository_Expecter[T]) Update(ctx interface{}, entity interface{}) *MockRepository_Update_Call[T] {
	return &MockRepository_Update_Call[T]{Call: _e.mock.On("Update", ctx, entity)}
}

func (_c *MockRepository_Update_Call[T]) Run(run func(ctx context.Context, entity T)) *MockRepository_Update_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(T))
	})
	return _c
}

func (_c *MockRepository_Update_Call[T]) Return(_a0 error) *MockRepository_Update_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Update_Call[T]) RunAndReturn(run func(context.Context, T) error) *MockRepository_Update_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository[T snapmatchai.Entity](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository[T] {
	mock := &MockRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
