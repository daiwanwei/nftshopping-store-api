package sagas

import (
	"bytes"
	"context"
	"fmt"
)

type SagaContext interface {
	context.Context
	Saga
}

type sagaContext struct {
	context.Context
	Saga
}

func NewSagaContext(ctx context.Context, name string) SagaContext {
	return &sagaContext{
		Context: ctx,
		Saga:    NewSaga(name),
	}
}

type Saga interface {
	Abort(ctx context.Context) error
	AddStep(step *Step)
}

type saga struct {
	Name  string
	steps []*Step
}

func NewSaga(name string) Saga {
	return &saga{
		Name: name,
	}
}

func (s *saga) Abort(ctx context.Context) error {
	if len(s.steps) == 0 {
		return nil
	}
	var errs MultiError
	for i := len(s.steps) - 1; i >= 0; i-- {
		if s.steps[i] == nil {
			continue
		}
		err := s.steps[i].CompensateFunc()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (s *saga) AddStep(step *Step) {
	s.steps = append(s.steps, step)
	return
}

type Step struct {
	Name           string
	CompensateFunc func() error
}

type MultiError []error

func (e MultiError) Error() string {
	if len(e) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d error(s) occurred\n------------------\n%s", len(e), e[0])
	for i := 1; i < len(e); i++ {
		fmt.Fprintf(buf, "\n\n&&&&&\n\n%s", e[i])
	}
	return buf.String()
}
