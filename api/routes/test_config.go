package routes

import "errors"

// MockDB is a dinero/api/models Store implementation,
// used for mocking responses from a mock database
type MockDB struct {
	dbErr bool
}

// ErrReader is a utility to produce an error from
// ioutil.ReadAll in tests
type ErrReader int

func (ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("ioutil.ReadAll error")
}
