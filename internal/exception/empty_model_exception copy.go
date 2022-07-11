package exception

import "fmt"

type EmptyModelException struct {
	Err error
}

func NewEmptyModelException(model string) *EmptyModelException {
	return &EmptyModelException{Err: fmt.Errorf("empty model: %s", model)}
}

func (e *EmptyModelException) Error() string {
	return e.Err.Error()
}
