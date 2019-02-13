package routes_test

import (
	"bytes"
	"dinero/api/config"
	"dinero/api/models"
	"dinero/api/routes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
)

func (mdb *MockDB) AllUsers() ([]*models.User, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	users := make([]*models.User, 0)
	users = append(users, &models.User{ID: 1, FirstName: "John", LastName: "Ide", FullName: "John Ide", Email: "ide.johnc@gmail.com", BiweeklyIncome: 1860.00})
	users = append(users, &models.User{ID: 2, FirstName: "Luke", LastName: "Toth", FullName: "Luke Toth", Email: "lptoth55@gmail.com", BiweeklyIncome: 1400.00})

	return users, nil
}

func (mdb *MockDB) GetUser(userID int) (*models.User, error) {
	if userID != 1 {
		return nil, models.ErrNotFound
	}

	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	user := &models.User{ID: 1, FirstName: "Luke", LastName: "Toth", FullName: "Luke Toth", Email: "lptoth55@gmail.com", BiweeklyIncome: 1400.00}

	return user, nil
}

func (mdb *MockDB) CreateUser(u models.User) (*models.User, error) {
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

	return user, nil
}

func (mdb *MockDB) UpdateUser(userID int, u *models.User) error {
	if mdb.dbErr {
		return errors.New("Database error")
	}

	if u.Email == "already-here@gmail.com" {
		return sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintUnique,
		}
	}

	return nil
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
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/users", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
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
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because a user with the ID of 3 is not being found
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/3", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNotFound,
		},
		{
			// breaks the test because "test" is not an integer
			name:           "BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/test", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "CTX_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
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
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/users", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the request body is set to produce an error
			name:           "BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", ErrReader(0)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "firstName" key in the request body is not a string
			name:           "BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":123,"lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "email" key in the request body is not a valid email
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"invalid.email","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			// breaks the test because the "email" key in the request body ("already-here@gmail.com") is set to cause a conflict
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"already-here@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
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

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK_NO_CONTENT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(`{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "OK_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/3", bytes.NewBuffer([]byte(`{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusCreated,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the request body is set to produce an error
			name:           "BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", ErrReader(0)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "firstName" key in the request body is not a string
			name:           "BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(`{"firstName":123,"lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "email" key in the request body is not a valid email
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"invalid.email","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			// breaks the test because the "email" key in the request body ("already-here@gmail.com") is set to cause a conflict
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"already-here@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the "email" key in the request body ("already-here@gmail.com") is set to cause a conflict
			name:           "SQLITE_CONFLICT_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/3", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"already-here@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(`{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/3", bytes.NewBuffer([]byte(`{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "CTX_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/users/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "CTX_ERR" {
				http.HandlerFunc(routes.UpdateUser(test.env)).ServeHTTP(test.rec, test.req)

				RunTest(&test, t)
			} else {
				r := routes.NewRouter(test.env)
				r.ServeHTTP(test.rec, test.req)

				RunTest(&test, t)
			}
		})
	}
}
