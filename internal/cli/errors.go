package cli

import "fmt"

const (
	ExitOK          = 0
	ExitError       = 1
	ExitConfigError = 2
)

// CLIError is an error with an exit code.
type CLIError struct {
	Code    int
	Message string
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewError creates a CLIError.
func NewError(code int, msg string, args ...interface{}) *CLIError {
	return &CLIError{Code: code, Message: fmt.Sprintf(msg, args...)}
}
