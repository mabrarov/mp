package panicerr

import "fmt"

type Wrapper struct {
	Value any
}

func New(value any) *Wrapper {
	return &Wrapper{Value: value}
}

func (e *Wrapper) Error() string {
	if err, ok := e.Value.(error); ok {
		return err.Error()
	}
	if s, ok := e.Value.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", e.Value)
}
