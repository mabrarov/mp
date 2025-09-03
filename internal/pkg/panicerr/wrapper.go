package panicerr

import "fmt"

type Wrapper struct {
	value any
}

func New(value any) *Wrapper {
	return &Wrapper{value: value}
}

func (e *Wrapper) Error() string {
	if err, ok := e.value.(error); ok {
		return err.Error()
	}
	if s, ok := e.value.(fmt.Stringer); ok {
		return s.String()
	}
	if s, ok := e.value.(string); ok {
		return s
	}
	return "panic"
}
