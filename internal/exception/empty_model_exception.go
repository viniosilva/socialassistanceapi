package exception

type EmptyModelException struct {
	Err error
}

func (e *EmptyModelException) Error() string {
	return e.Err.Error()
}
