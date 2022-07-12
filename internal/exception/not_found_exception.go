package exception

type NotFoundException struct {
	Err error
}

func (e *NotFoundException) Error() string {
	return e.Err.Error()
}
