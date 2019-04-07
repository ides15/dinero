package models

import "errors"

var (
	// ErrNotFound is an error creator for models where DB query results return nothing
	ErrNotFound = errors.New("error: no record(s) found")
	// ErrBadPing is an error creator for the DB model where the application errors in
	// pinging the database
	ErrBadPing = errors.New("error: cannot ping database")
)
