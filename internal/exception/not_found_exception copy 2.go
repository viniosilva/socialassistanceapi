package exception

import "fmt"

type NotFoundException struct {
	Err error
}

func NewNotFoundException(model string) *NotFoundException {
	return &NotFoundException{Err: fmt.Errorf("not found: %s", model)}
}

func (e *NotFoundException) Error() string {
	return e.Err.Error()
}
