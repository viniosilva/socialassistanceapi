package exception

import "fmt"

type EmptyAddressModelException struct {
	Err error
}

func NewEmptyAddressModelException() *EmptyAddressModelException {
	return &EmptyAddressModelException{Err: fmt.Errorf("empty address model")}
}

func (e *EmptyAddressModelException) Error() string {
	return e.Err.Error()
}
