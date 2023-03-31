package DB

type MigratorErrorType int64

type MigratorError struct {
	_type MigratorErrorType
}

const (
	CouldNotCreate MigratorErrorType = -1
	CouldNotFind   MigratorErrorType = -1
)

func (m *MigratorError) Error() string {
	return "boom"
}
