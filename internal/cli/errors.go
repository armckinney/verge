package cli

import "fmt"

// Exit codes as per spec
const (
	ExitOK           = 0
	ExitError        = 1
	ExitUsageError   = 2
	ExitCompareLeft  = 10 // left < right
	ExitCompareRight = 11 // left > right
	ExitNotFound     = 20 // version not found in source
	ExitParseError   = 21 // invalid version format
	ExitConfigError  = 22 // config file error
	ExitNetworkError = 30 // network/provider error
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
