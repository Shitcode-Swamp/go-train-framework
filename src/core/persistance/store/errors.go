package store

import "errors"

var (
	ErrNotFound     = errors.New("store: not found")
	ErrDuplicateKey = errors.New("store: duplicate key")

	ErrConstraintViolation = errors.New("store: constraint violation")
	ErrConflict            = errors.New("store: write conflict")

	ErrInvalidFilter   = errors.New("store: invalid filter")
	ErrInvalidDocument = errors.New("store: invalid document")

	ErrUnsupportedOperation = errors.New("store: unsupported operation")
	ErrConnection           = errors.New("store: connection failure")
)
