package pipeline

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Step represents an individual operation within the pipeline
type Step[T any] interface {
	Execute(T) error
	Name() string
}

type NamedStep[T any] struct {
	StepName  string
	ExecuteFn func(T) error
}

func (n NamedStep[T]) Execute(state T) error {
	return n.ExecuteFn(state)
}

func (n NamedStep[T]) Name() string {
	return n.StepName
}

type ErrorHandler[T any] func(T, error) error

// Pipeline struct holds the pipeline state and a stack of steps
type Pipeline[T any] struct {
	state  T
	stack  []Steps[T]
	err    error
	errors []StepError
	logger Logger
}

type Steps[T any] struct {
	Step            Step[T]
	StepName        string
	OnError         ErrorHandler[T]
	OnSuccessEnd    bool
	OnErrorContinue bool
}

type Settings struct {
	Logger Logger
}
type ClientOption interface {
	Apply(settings *Settings)
}

func WithLogger(logger Logger) ClientOption {
	return loggerOption{l: logger}
}

type loggerOption struct{ l Logger }

func (w loggerOption) Apply(o *Settings) {
	o.Logger = w.l
}

// New initializes a new pipeline with the provided client and request
func New[T any](state T, opts ...ClientOption) *Pipeline[T] {
	settings := &Settings{}
	for _, opt := range opts {
		opt.Apply(settings)
	}
	return &Pipeline[T]{state: state, logger: settings.Logger}
}

// Then adds a new step to the pipeline
func (p *Pipeline[T]) Then(f Step[T]) *Pipeline[T] {
	s := Steps[T]{Step: f}
	p.stack = append(p.stack, s)
	p.info("Added step", "stepName", f.Name())
	return p
}

// OnError sets the error handler for the previous step
func (p *Pipeline[T]) OnError(onErr ErrorHandler[T]) *Pipeline[T] {
	p.stack[len(p.stack)-1].OnError = onErr
	p.info("Set error handler", "stepName", p.stack[len(p.stack)-1].Step.Name())
	return p
}

// OnSuccessEnd stops pipeline execution after the current step if set to true
func (p *Pipeline[T]) OnSuccessEnd() *Pipeline[T] {
	p.stack[len(p.stack)-1].OnSuccessEnd = true
	p.info("Step will end pipeline execution on success", "stepName", p.stack[len(p.stack)-1].Step.Name())
	return p
}

// OnErrorContinue continues pipeline execution even if the previous step failed
func (p *Pipeline[T]) OnErrorContinue() *Pipeline[T] {
	p.stack[len(p.stack)-1].OnErrorContinue = true
	p.info("Step will continue execution after failure", "stepName", p.stack[len(p.stack)-1].Step.Name())
	return p
}

// Errors returns a slice of all errors encountered during the pipeline execution
func (p *Pipeline[T]) Errors() []StepError {
	return p.errors
}

// Execute runs the pipeline and logs various steps
func (p *Pipeline[T]) Execute() error {
	for i, s := range p.stack {
		p.info("Executing step", "stepIndex", i, "stepName", s.Step.Name())
		err := s.Step.Execute(p.state)
		if err != nil {
			stepError := StepError{
				StepIndex: i,
				StepName:  s.Step.Name(),
				Err:       err,
			}
			p.errors = append(p.errors, stepError)
			p.error("Step failed", "stepIndex", i, "stepName", s.Step.Name(), "error", err)
			if s.OnError != nil {
				onErr := s.OnError(p.state, err)
				if onErr != nil && !s.OnErrorContinue {
					p.error("Error handler failed", "stepIndex", i, "error", onErr)
					return onErr
				}
			} else {
				if !s.OnErrorContinue {
					return err
				}
			}
		}

		if s.OnSuccessEnd {
			p.info("Pipeline execution stopped after step", "stepIndex", i, "stepName", s.Step.Name())
			break
		}
	}
	p.info("Pipeline execution completed")
	return nil
}

func (p *Pipeline[T]) info(msg string, args ...interface{}) {
	if p.logger == nil {
		return
	}
	p.logger.Info(msg, args...)
}

func (p *Pipeline[T]) error(msg string, args ...interface{}) {
	if p.logger == nil {
		return
	}
	p.logger.Error(msg, args...)
}

func (p *Pipeline[T]) debug(msg string, args ...interface{}) {
	if p.logger == nil {
		return
	}
	p.logger.Debug(msg, args...)
}
