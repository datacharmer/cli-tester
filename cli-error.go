package cli_tester

import "fmt"

// CliRunError implements the error interface
// It is like a common error, with additional information about the set being tested
type CliRunError struct {
	ErrorMessage string
	RunSet       *CliSet
}

// SetCliRunError creates a new CliRunError from a CliSet and a formatted string
func SetCliRunError(cliSet *CliSet, format string, args ...interface{}) *CliRunError {
	return &CliRunError{
		RunSet:       cliSet,
		ErrorMessage: fmt.Sprintf(format, args...),
	}
}

// Error implements the error interface
func (cre *CliRunError) Error() string {
	return cre.ErrorMessage
}

// IsFailed returns true if any of the errors is not nil
func IsFailed(errors []*CliRunError) bool {
	if errors == nil {
		return false
	}
	for _, err := range errors {
		if err != nil {
			return true
		}
	}
	return false
}

func ErrorMessages(errors []*CliRunError) string {
	message := ""
	for i, err := range errors {
		if message != "" {
			message += "\n"
		}
		message += fmt.Sprintf("%2d - %s", i, err)
	}
	return message
}

// FullError returns the error message from a CliRunError, and additional error output, if found
func (cre *CliRunError) FullError() string {
	if cre == nil {
		return ""
	}
	message := ""
	if cre.Error() != "" {
		message += cre.Error()
	}
	if cre.RunSet.Command.ErrOut != "" {
		message += " - " + cre.RunSet.Command.ErrOut
	}
	return message
}
