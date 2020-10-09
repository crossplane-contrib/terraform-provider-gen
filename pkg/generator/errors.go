package generator

import "fmt"

// MultiError gives client code the ability to type assert a possible
// multiError into a form where the individual errors can be examined.
type MultiError interface {
	Append(error)
	Errors() []error
	Error() string
}

// multiError allows Validation methods to collect multiple errors
// into a single type for deeper inspection
type multiError struct {
	errors  []error
	heading string
}

// Append adds a new error to the list of validation errors encountered
func (e *multiError) Append(err error) {
	e.errors = append(e.errors, err)
}

// Errors allows access to the underlying errors,
// partially satisfying the MultiError interface
func (e *multiError) Errors() []error {
	return e.errors
}

// Error() combines the underlying set of validation errors into a single error message
func (e *multiError) Error() string {
	// do not print heading if there are no errors to report
	if len(e.errors) == 0 {
		return ""
	}
	full := e.heading
	for _, err := range e.errors {
		full = fmt.Sprintf("%s\n - %s", full, err.Error())
	}
	return full
}

// NewMultiError handles initializing the errors slice and allows
// a heading to be set. The Error() method will print this heading
// before listing out the individual errors.
func NewMultiError(heading string) MultiError {
	return &multiError{
		heading: heading,
		errors:  make([]error, 0),
	}
}
