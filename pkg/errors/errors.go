package errors

import "errors"

var (
	ErrNotSupported             = errors.New("not supported")
	ErrDDLTagGoNotFoundInSource = errors.New("ddl-tag-go not found in source")
)
