package errs

import "errors"

var (
	ErrNotFound     = errors.New("data not found")
	ErrSaveFailed   = errors.New("failed to save data")
	ErrDeleteFailed = errors.New("failed to delete data")
	ErrUpdateFailed = errors.New("failed to update data")
)
