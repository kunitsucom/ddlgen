package errors

import "errors"

var (
	ErrUnknownError                  = errors.New("unknown error")
	ErrNotSupported                  = errors.New("not supported")
	ErrUnformattedFileIsNotSupported = errors.New("unformatted file is not supported")
	ErrDDLTagGoNotFoundInSource      = errors.New("ddl-tag-go not found in source")
)
