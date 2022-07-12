package exception

type NegativeException struct {
	Err error
}

func (e *NegativeException) Error() string {
	return e.Err.Error()
}
