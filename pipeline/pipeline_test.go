package pipeline

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

type pipelineState struct {
	name string
}

func TestNewPipeline(t *testing.T) {
	logger := slog.Default()
	pipeline := New(pipelineState{}, WithLogger(logger))

	assert.NotNil(t, pipeline)
	assert.Equal(t, logger, pipeline.logger)
	assert.Equal(t, pipelineState{}, pipeline.state)
}

func TestPipelineThen(t *testing.T) {
	pipeline := New(pipelineState{})
	pipeline.Then(NamedStep[pipelineState]{
		StepName: "Step",
		ExecuteFn: func(state pipelineState) error {
			return nil
		},
	})

	assert.Len(t, pipeline.stack, 1)
	assert.Equal(t, "Step", pipeline.stack[0].Step.Name())
}

func TestPipelineOnError(t *testing.T) {
	pipeline := New(pipelineState{})

	pipeline.Then(NamedStep[pipelineState]{ExecuteFn: func(state pipelineState) error {
		return nil
	}})

	// Define mock error handler
	mockErrorHandler := func(state pipelineState, err error) error {
		return nil
	}
	pipeline.OnError(mockErrorHandler)

	assert.NotNil(t, pipeline.stack[0].OnError)
}

func TestPipelineExecute_Success(t *testing.T) {
	pipeline := New(pipelineState{})

	pipeline.Then(NamedStep[pipelineState]{ExecuteFn: func(state pipelineState) error {
		return nil
	}})

	err := pipeline.Execute()

	assert.NoError(t, err)
	assert.Len(t, pipeline.errors, 0)
}

func TestPipelineExecute_ErrorHandler(t *testing.T) {
	pipeline := New(pipelineState{})

	pipeline.Then(NamedStep[pipelineState]{
		StepName: "Step",
		ExecuteFn: func(state pipelineState) error {
			return errors.New("execution error")
		},
	})

	mockErrorHandler := func(state pipelineState, err error) error {
		return nil
	}
	pipeline.OnError(mockErrorHandler)

	err := pipeline.Execute()

	assert.NoError(t, err)
	assert.Len(t, pipeline.errors, 1)
	assert.Equal(t, "execution error", pipeline.errors[0].Err.Error())
}

func TestPipelineExecute_OnErrorContinue(t *testing.T) {
	pipeline := New(pipelineState{})

	pipeline.Then(NamedStep[pipelineState]{
		StepName: "Step 1",
		ExecuteFn: func(state pipelineState) error {
			return errors.New("execution error")
		},
	}).OnErrorContinue().Then(NamedStep[pipelineState]{
		StepName: "Step 2",
		ExecuteFn: func(state pipelineState) error {
			return errors.New("execution error 2")
		},
	})

	err := pipeline.Execute()

	assert.Error(t, err)
	assert.Len(t, pipeline.errors, 2)
	assert.Equal(t, "execution error", pipeline.errors[0].Err.Error())
	assert.Equal(t, "execution error 2", pipeline.errors[1].Err.Error())

}
