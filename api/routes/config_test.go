package routes_test

import (
	"dinero/api/config"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDB is a dinero/api/models Store implementation,
// used for mocking responses from a mock database
type MockDB struct {
	dbErr bool
}

// TestCase defines the structure for a route test case
type TestCase struct {
	name           string
	rec            *httptest.ResponseRecorder
	req            *http.Request
	env            *config.Env
	expectedBody   string
	expectedHeader string
	expectedStatus int
}

// ErrReader is a utility to produce an error from
// ioutil.ReadAll in tests
type ErrReader int

func (ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("ioutil.ReadAll error")
}

// Test runs test cases
func RunTest(c *TestCase, t *testing.T) {
	if c.expectedBody != c.rec.Body.String() {
		t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", c.rec.Body.String(), c.expectedBody)
	}

	if c.expectedHeader != c.rec.Header().Get("Content-Type") {
		t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", c.rec.Body.String(), c.expectedHeader)
	}

	if c.expectedStatus != c.rec.Code {
		t.Errorf("\nCode:\n\tGot: \t\t%d\n\tExpected: \t%d\n", c.rec.Code, c.expectedStatus)
	}
}
