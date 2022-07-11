package exception

import "fmt"

type NegativeException struct {
	Err error
}

func NewNegativeException(model string) *NegativeException {
	return &NegativeException{Err: fmt.Errorf("negative value in: %s", model)}
}

func (e *NegativeException) Error() string {
	return e.Err.Error()
}
