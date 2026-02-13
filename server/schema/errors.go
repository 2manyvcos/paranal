package schema

type err string

func (e err) Error() string {
	return string(e)
}

const (
	ErrConflict       = err("conflict")
	ErrNotFound       = err("not found")
	ErrAlreadyRunning = err("already running")
)
