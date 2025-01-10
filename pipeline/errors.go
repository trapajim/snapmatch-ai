package pipeline

import "fmt"

type StepError struct {
	StepIndex int
	StepName  string
	Err       error
}

func (e *StepError) Error() string {
	return fmt.Sprintf("Step %d (%s) failed: %v", e.StepIndex, e.StepName, e.Err)
}

func (e *StepError) Unwrap() error {
	return e.Err
}
