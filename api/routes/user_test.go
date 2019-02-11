package routes_test

import (
	"bytes"
	"dinero/api/config"
	"dinero/api/models"
	"dinero/api/routes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
)

func (mdb *MockDB) AllUsers() ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	users := make([]*models.User, 0)
	users = append(users, &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.00})
	users = append(users, &models.User{ID: 2, FirstName: "Luke", LastName: "Toth", FullName: "Luke Toth", Email: "lptoth55@gmail.com", BiweeklyIncome: 1400.00})

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return usersJSON, nil
}

func (mdb *MockDB) GetUser(userID int) ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	if userID != 1 {
		return nil, models.ErrNotFound
	}

	user := &models.User{ID: 1, FirstName: "Luke", LastName: "Toth", FullName: "Luke Toth", Email: "lptoth55@gmail.com", BiweeklyIncome: 1400.00}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return userJSON, nil
}

func (mdb *MockDB) CreateUser(u models.User) ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	if u.Email == "already-here@gmail.com" {
		return nil, sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintUnique,
		}
	}

	user := &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.99}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return userJSON, nil
}

func TestAllUsers(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `[{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860},{"ID":2,"firstName":"Luke","lastName":"Toth","fullName":"Luke Toth","email":"lptoth55@gmail.com","biweeklyIncome":1400}]`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users", nil), // breaks the test because the PUT method is not allowed
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := routes.NewRouter(test.env)
			r.ServeHTTP(test.rec, test.req)

			RunTest(&test, t)
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"firstName":"Luke","lastName":"Toth","fullName":"Luke Toth","email":"lptoth55@gmail.com","biweeklyIncome":1400}`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", nil), // breaks the test because the PUT method is not allowed
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/3", nil), // breaks the test because a user with the ID of 3 is not being found
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/test", nil), // breaks the test because "test" is not an integer
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "CTX_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "CTX_ERR" {
				http.HandlerFunc(routes.GetUser(test.env)).ServeHTTP(test.rec, test.req)

				RunTest(&test, t)
			} else {
				r := routes.NewRouter(test.env)
				r.ServeHTTP(test.rec, test.req)

				RunTest(&test, t)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", ErrReader(0)), // breaks the test because the request body is set to produce an error
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":123,"lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))), // breaks the test because the "firstName" key in the request body is not a string
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"invalid.email","biweeklyIncome":1860.99}`))), // breaks the test because the "email" key in the request body is not a valid email
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"already-here@gmail.com","biweeklyIncome":1860.99}`))), // breaks the test because the "email" key in the request body ("already-here@gmail.com") is set to cause a conflict
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := routes.NewRouter(test.env)
			r.ServeHTTP(test.rec, test.req)

			RunTest(&test, t)
		})
	}
}
