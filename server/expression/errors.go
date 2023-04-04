package expression

import (
	"errors"
	"fmt"
)

var ErrExpressionResolution error = errors.New("resolution error")
var ErrInvalidSyntax error = errors.New("invalid syntax")

type ResolutionError struct {
	innerErr error
}

func (e *ResolutionError) Error() string {
	return fmt.Sprintf("%s: %s", ErrExpressionResolution.Error(), e.innerErr.Error())
}

func (e *ResolutionError) Is(target error) bool {
	// all Resolution errors are ErrExpressionResolution
	if errors.Is(target, ErrExpressionResolution) {
		return true
	}

	return errors.Is(e.innerErr, target)
}

func (e *ResolutionError) Unwrap() error {
	return e.innerErr
}

func resolutionError(innerErr error) error {
	return &ResolutionError{innerErr: innerErr}
}

type InvalidSyntaxError struct {
	innerErr error
}

func (e *InvalidSyntaxError) Error() string {
	return e.innerErr.Error()
}

func (e *InvalidSyntaxError) Is(target error) bool {
	if errors.Is(target, ErrInvalidSyntax) {
		return true
	}

	return errors.Is(e.innerErr, target)
}

func (e *InvalidSyntaxError) Unwrap() error {
	return errors.Unwrap(e.innerErr)
}

func invalidSyntaxError(err error, expression string) error {
	innerErr := fmt.Errorf(`invalid syntax "%s": %w`, expression, err)
	return &InvalidSyntaxError{innerErr: innerErr}
}
