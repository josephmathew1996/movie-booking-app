package errors

type ErrorCode int

const (
	ErrBadRequest ErrorCode = iota + 1001
	ErrValidationFailed
	ErrUserAlreadyExists
	ErrInternalServer
)

type CustomError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return e.Message
}

func New(code ErrorCode, message string, err error) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
