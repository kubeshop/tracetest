package expression

import (
	"errors"
	"fmt"
)

var ErrExpressionResolution error = errors.New("resolution error")

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
