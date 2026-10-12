package pipeline

// CLIError carries a specific exit code for the CLI
// (contracts/cli.md §generate Exit Codes).
type CLIError struct {
	Code    int
	Message string
}

func (e *CLIError) Error() string { return e.Message }

// Exit codes from contracts/cli.md.
const (
	ExitSuccess          = 0
	ExitValidationError  = 1
	ExitDestinationError = 2
	ExitTemplateError    = 3
	ExitCancelled        = 4
)

// ValidationError wraps a configuration validation failure.
func ValidationError(message string) *CLIError {
	return &CLIError{Code: ExitValidationError, Message: message}
}

// DestinationError wraps a destination protection failure.
func DestinationError(message string) *CLIError {
	return &CLIError{Code: ExitDestinationError, Message: message}
}

// TemplateError wraps a template or IO failure.
func TemplateError(message string) *CLIError {
	return &CLIError{Code: ExitTemplateError, Message: message}
}

// CancelledError wraps a user-cancellation.
func CancelledError(message string) *CLIError {
	return &CLIError{Code: ExitCancelled, Message: message}
}
