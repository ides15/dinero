package models

import "errors"

var (
	// ErrNotFound is an error creator for models where DB query results return nothing
	ErrNotFound = errors.New("error: no record(s) found")
)
