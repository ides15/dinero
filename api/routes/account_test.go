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

func (mdb *MockDB) AllAccounts() ([]*models.Account, error) {
	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	accounts := make([]*models.Account, 0)
	accounts = append(accounts, &models.Account{ID: 1, UserID: 1, Name: "Car Payment", AccountType: "monthly", MinimumPayment: 217.99, CurrentPayment: 217.99, DueDate: "12"})
	accounts = append(accounts, &models.Account{ID: 2, UserID: 1, Name: "Phone Payment", AccountType: "monthly", MinimumPayment: 42.83, CurrentPayment: 100.00, FullAmount: 728.00, DueDate: "10", URL: "https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"})

	return accounts, nil
}

func (mdb *MockDB) GetAccount(accountID int) (*models.Account, error) {
	if accountID != 1 {
		return nil, models.ErrNotFound
	}

	if mdb.dbErr {
		return nil, errors.New("Database error")
	}

	account := &models.Account{ID: 1, UserID: 1, Name: "Phone Payment", AccountType: "monthly", MinimumPayment: 42.83, CurrentPayment: 100.00, FullAmount: 728.00, DueDate: "10", URL: "https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}

	return account, nil
}

func (mdb *MockDB) CreateAccount(a models.Account) (*models.Account, error) {
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

	return account, nil
}

func (mdb *MockDB) UpdateAccount(accountID int, a *models.Account) error {
	if mdb.dbErr {
		return errors.New("Database error")
	}

	if a.Name == "Already here" {
		return sqlite3.Error{
			Code:         sqlite3.ErrConstraint,
			ExtendedCode: sqlite3.ErrConstraintUnique,
		}
	}

	return nil
}

func (mdb *MockDB) DeleteAccount(accountID int) error {
	if accountID != 1 {
		return models.ErrNotFound
	}

	if mdb.dbErr {
		return errors.New("Database error")
	}

	return nil
}

func TestAllAccounts(t *testing.T) {
	t.Parallel()

	must := func(req *http.Request, err error) *http.Request {
		return req
	}

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/accounts", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `[{"ID":1,"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":0,"dueDate":"12","URL":""},{"ID":2,"userID":1,"name":"Phone Payment","accountType":"monthly","minimumPayment":42.83,"currentPayment":100,"fullAmount":728,"dueDate":"10","URL":"https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}]`,
			expectedHeader: "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/accounts", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/accounts", nil)),
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

func TestGetAccount(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"userID":1,"name":"Phone Payment","accountType":"monthly","minimumPayment":42.83,"currentPayment":100,"fullAmount":728,"dueDate":"10","URL":"https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}`,
			expectedHeader: "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/accounts/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because an account with the ID of 3 is not being found
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/3", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNotFound,
		},
		{
			// breaks the test because "test" is not an integer
			name:           "BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/test", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("GET", "/accounts/1", nil),
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

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`,
			expectedHeader: "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/accounts", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the request body is set to produce an error
			name:           "BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", ErrReader(0)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "accountType" key in the request body is not a string
			name:           "BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":123,minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "minimumPayment" key in the request body is not a float64
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"bad","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			// breaks the test because the "name" key in the request body ("Already here") is set to cause a conflict
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Already here","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
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

func TestUpdateAccount(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK_NO_CONTENT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "OK_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/3", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusCreated,
		},
		{
			// breaks the test because the BAD method is not allowed
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("BAD", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			// breaks the test because the request body is set to produce an error
			name:           "BAD_REQUEST_IOUTIL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", ErrReader(0)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "Name" key in the request body is not a string
			name:           "BAD_REQUEST_UNMARSHAL",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":12345,"accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because "test" is not an integer
			name:           "BAD_REQUEST_STRING",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/test", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			// breaks the test because the "accountType" key in the request body is not a valid account type
			name:           "INVALID",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"bad","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusUnprocessableEntity)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			// breaks the test because the "name" key in the request body ("Already here") is set to cause a conflict
			name:           "SQLITE_CONFLICT",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":"Already here","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the "name" key in the request body ("Already here") is set to cause a conflict
			name:           "SQLITE_CONFLICT_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/3", bytes.NewBuffer([]byte(`{"userID":1,"name":"Already here","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusConflict)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusConflict,
		},
		{
			// breaks the test because the env.DB is set to have a dbErr
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/1", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DB_ERR_CREATED",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("PUT", "/accounts/3", bytes.NewBuffer([]byte(`{"userID":1,"name":"Car Payment","accountType":"monthly","minimumPayment":217.99,"currentPayment":217.99,"fullAmount":21000,"dueDate":"10","URL":"ford.com"}`))),
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

func TestDeleteAccount(t *testing.T) {
	t.Parallel()

	tests := []TestCase{
		{
			name:           "OK",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("DELETE", "/accounts/1", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   "",
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNoContent,
		},
		{
			// breaks the test because "test" is not an integer
			name:           "BAD_REQUEST_STRING",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("DELETE", "/accounts/test", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("DELETE", "/accounts/3", nil),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            httptest.NewRequest("DELETE", "/accounts/1", nil),
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
