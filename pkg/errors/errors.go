package errors

import "errors"

var (
	ErrNotSupported                       = errors.New("not supported")
	ErrDDLTagGoAnnotationNotFoundInSource = errors.New("ddl-tag-go annotation not found in source")
)
