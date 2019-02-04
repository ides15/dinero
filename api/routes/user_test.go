package routes

import (
	"bytes"
	"dinero/api/config"
	"dinero/api/models"
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

	fmt.Println(u.Email)
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

func TestUsersHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		env            *config.Env
		expectedBody   string
		expectedHeader string
	}{
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
			req:            httptest.NewRequest("PUT", "/users", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/users", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			http.HandlerFunc(UsersHandler(test.env)).ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}

func TestUserHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		env            *config.Env
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "GET_OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/user/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"firstName":"Luke","lastName":"Toth","fullName":"Luke Toth","email":"lptoth55@gmail.com","biweeklyIncome":1400}`,
			expectedHeader: "application/json",
		},
		{
			name:           "GET_BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/user/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "GET_NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/user/3", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "GET_BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/user/test", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "GET_DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/user/1", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "POST_OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`,
			expectedHeader: "application/json",
		},
		{
			name:           "POST_BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", ErrReader(0)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "POST_BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{"firstName":123,"lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "POST_INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"invalid.email","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "POST_SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"already-here@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "POST_DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{"firstName":"John","lastName":"Ide","fullName":"John Ide","email":"ide.johnc@gmail.com","biweeklyIncome":1860.99}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			http.HandlerFunc(UserHandler(test.env)).ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}
