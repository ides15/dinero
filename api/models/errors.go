package models

import "errors"

var (
	// ErrNotFound is an error creator for models where DB query results return nothing
	ErrNotFound = errors.New("Error: no record(s) found")
)
