package errs

import "errors"

var (
	ErrBadRequest         = errors.New("invalid request data")
	ErrInternalServer     = errors.New("internal server error")
	ErrDataNotFound       = errors.New("data not found")
	ErrInvalidIDParameter = errors.New("invalid id parameter")
)
