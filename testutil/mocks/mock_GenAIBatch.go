// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	snapmatchai "github.com/trapajim/snapmatch-ai/snapmatchai"
)

// MockGenAIBatch is an autogenerated mock type for the GenAIBatch type
type MockGenAIBatch struct {
	mock.Mock
}

type MockGenAIBatch_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGenAIBatch) EXPECT() *MockGenAIBatch_Expecter {
	return &MockGenAIBatch_Expecter{mock: &_m.Mock}
}

// CreateBatchPredictionJob provides a mock function with given fields: _a0, _a1
func (_m *MockGenAIBatch) CreateBatchPredictionJob(_a0 context.Context, _a1 snapmatchai.BatchPredictionRequest) (snapmatchai.BatchPredictionJobConfig, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateBatchPredictionJob")
	}

	var r0 snapmatchai.BatchPredictionJobConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, snapmatchai.BatchPredictionRequest) (snapmatchai.BatchPredictionJobConfig, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, snapmatchai.BatchPredictionRequest) snapmatchai.BatchPredictionJobConfig); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(snapmatchai.BatchPredictionJobConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, snapmatchai.BatchPredictionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenAIBatch_CreateBatchPredictionJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBatchPredictionJob'
type MockGenAIBatch_CreateBatchPredictionJob_Call struct {
	*mock.Call
}

// CreateBatchPredictionJob is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 snapmatchai.BatchPredictionRequest
func (_e *MockGenAIBatch_Expecter) CreateBatchPredictionJob(_a0 interface{}, _a1 interface{}) *MockGenAIBatch_CreateBatchPredictionJob_Call {
	return &MockGenAIBatch_CreateBatchPredictionJob_Call{Call: _e.mock.On("CreateBatchPredictionJob", _a0, _a1)}
}

func (_c *MockGenAIBatch_CreateBatchPredictionJob_Call) Run(run func(_a0 context.Context, _a1 snapmatchai.BatchPredictionRequest)) *MockGenAIBatch_CreateBatchPredictionJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(snapmatchai.BatchPredictionRequest))
	})
	return _c
}

func (_c *MockGenAIBatch_CreateBatchPredictionJob_Call) Return(_a0 snapmatchai.BatchPredictionJobConfig, _a1 error) *MockGenAIBatch_CreateBatchPredictionJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenAIBatch_CreateBatchPredictionJob_Call) RunAndReturn(run func(context.Context, snapmatchai.BatchPredictionRequest) (snapmatchai.BatchPredictionJobConfig, error)) *MockGenAIBatch_CreateBatchPredictionJob_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGenAIBatch creates a new instance of MockGenAIBatch. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGenAIBatch(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGenAIBatch {
	mock := &MockGenAIBatch{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
