package routes

import (
	"dinero/api/config"
	"dinero/api/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestAccountsHandler(t *testing.T) {
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
			req:            must(http.NewRequest("PUT", "/accounts", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/accounts", nil)),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			http.HandlerFunc(AccountsHandler(test.env)).ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}

func TestAccountHandler(t *testing.T) {
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
			req:            must(http.NewRequest("GET", "/account/1", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   `{"ID":1,"userID":1,"name":"Phone Payment","accountType":"monthly","minimumPayment":42.83,"currentPayment":100,"fullAmount":728,"dueDate":"10","URL":"https://www.synchronycredit.com/eService/AccountSummary/initiateAccSummaryAction.action"}`,
			expectedHeader: "application/json",
		},
		{
			name:           "BAD_METHOD",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("PUT", "/account/1", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusMethodNotAllowed)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "NOT_FOUND",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/account/3", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusNotFound)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "BAD_REQUEST",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/account/test", nil)),
			env:            &config.Env{DB: &MockDB{}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusBadRequest)),
			expectedHeader: "text/plain; charset=utf-8",
		},
		{
			name:           "DB_ERR",
			rec:            httptest.NewRecorder(),
			req:            must(http.NewRequest("GET", "/account/1", nil)),
			env:            &config.Env{DB: &MockDB{dbErr: true}, Log: config.Log},
			expectedBody:   fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)),
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			http.HandlerFunc(AccountHandler(test.env)).ServeHTTP(test.rec, test.req)

			if test.expectedBody != test.rec.Body.String() {
				t.Errorf("\nBody:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedBody)
			}

			if test.expectedHeader != test.rec.Header().Get("Content-Type") {
				t.Errorf("\nHeader:\n\tGot: \t\t%s\n\tExpected: \t%s\n", test.rec.Body.String(), test.expectedHeader)
			}
		})
	}
}
