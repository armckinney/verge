package cli

import "fmt"

// Exit codes
const (
	ExitOK           = 0
	ExitError        = 1
	ExitUsageError   = 2
	ExitNotFound     = 3
	ExitParseError   = 4
	ExitCompareLeft  = 10 // left < right
	ExitCompareRight = 11 // left > right
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
