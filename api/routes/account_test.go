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

func (mdb *MockDB) AllAccounts() ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	accounts := make([]*models.Account, 0)
	accounts = append(accounts, &models.Account{ID: 1, UserID: 1, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 217.99, CurrentPayment: 217.99, DueDate: "12"})
	accounts = append(accounts, &models.Account{ID: 2, UserID: 1, Name: "Phone Payment", AccountType: "monthly", MinimumPayment: 42.83, CurrentPayment: 100.00, FullAmount: 728.00, DueDate: "10", URL: "https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"})

	accountsJSON, err := json.Marshal(accounts)
	if err != nil {
		return nil, err
	}

	return accountsJSON, nil
}

func (mdb *MockDB) GetAccount(accountID int) ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	if accountID != 1 {
		return nil, models.ErrNotFound
	}

	account := &models.Account{ID: 1, UserID: 1, Name: "Phone Payment", AccountType: "monthly", MinimumPayment: 42.83, CurrentPayment: 100.00, FullAmount: 728.00, DueDate: "10", URL: "https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	return accountJSON, nil
}

func (mdb *MockDB) CreateAccount(a models.Account) ([]byte, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	if a.Name == "Already here" {
		return nil, sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintUnique,
		}
	}

	account := &models.Account{ID: 1, UserID: 1, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 217.99, CurrentPayment: 217.99, FullAmount: 21000, DueDate: "10", URL: "ford.com"}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	return accountJSON, nil
}

func TestAllAccounts(t *testing.T) {
	t.Parallel()

	must := func(req *http.Request, err error) *http.Request {
		return req
	}

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
			req:            must(http.NewRequest("GET", "/accounts", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `[{"ID":1,"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":0,"dueDate":"12","URL":""},{"ID":2,"userID":1,"name":"Phone Payment","accountType":"monthly","minimumPayment":42.83,"currentPayment":100,"fullAmount":728,"dueDate":"10","URL":"https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}]`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("PUT", "/accounts", nil)), // breaks the test because the PUT method is not allowed
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/accounts", nil)),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewRouter(test.env)
			r.ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
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
			req:            httptest.NewRequest("GET", "/accounts/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"userID":1,"name":"Phone Payment","accountType":"monthly","minimumPayment":42.83,"currentPayment":100,"fullAmount":728,"dueDate":"10","URL":"https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", nil), // breaks the test because the PUT method is not allowed
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/3", nil), // breaks the test because an account with the ID of 3 is not being found
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/test", nil), // breaks the test because "test" is not an integer
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/1", nil),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log}, // breaks the test because the env.DB is set to have a dbErr
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewRouter(test.env)
			r.ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
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
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`,
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
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))), // breaks the test because the "accountType" key in the request body is not present
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":"bad","currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))), // breaks the test because the "minimumPayment" key in the request body is not a float64
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Already here","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))), // breaks the test because the "name" key in the request body ("Already here") is set to cause a conflict
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))), // breaks the test because the env.DB is set to have a dbErr
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewRouter(test.env)
			r.ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}
