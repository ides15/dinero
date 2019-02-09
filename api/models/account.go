package models

import (
	"database/sql"
	"encoding/json"
	"regexp"
)

// Account is an account a User wants to track
type Account struct {
	ID             int     `json:"ID"`
	UserID         int     `json:"userID"`
	Name           string  `json:"name"`
	AccountType    string  `json:"accountType"`
	MinimumPayment float64 `json:"minimumPayment"`
	CurrentPayment float64 `json:"currentPayment"`
	FullAmount     float64 `json:"fullAmount"`
	DueDate        string  `json:"dueDate"`
	URL            string  `json:"URL"`
}

// AllAccounts retrieves all account rows from the accounts table
func (db *DB) AllAccounts() ([]byte, error) {
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*Account, 0)
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.AccountType,
			&account.MinimumPayment,
			&account.CurrentPayment,
			&account.FullAmount,
			&account.DueDate,
			&account.URL)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	accountsJSON, err := json.Marshal(accounts)
	if err != nil {
		return nil, err
	}

	return accountsJSON, nil
}

// Validate validates the fields in an Account object
func (a *Account) Validate() bool {
	namePattern := regexp.MustCompile(`^[a-zA-Z ]+$`)
	typePattern := regexp.MustCompile(`^(daily|weekly|biweekly|monthly|yearly)$`)
	datePattern := regexp.MustCompile(`^([1-9]|[12]\d|3[01])$`)

	if a.UserID < 1 {
		return false
	}

	if !namePattern.MatchString(a.Name) {
		return false
	}

	if !typePattern.MatchString(a.AccountType) {
		return false
	}

	if a.MinimumPayment > a.FullAmount {
		return false
	}

	if a.CurrentPayment > a.FullAmount {
		return false
	}

	if !datePattern.MatchString(a.DueDate) {
		return false
	}

	return true
}

// GetAccount retrieves an account that matches the accountID parameter
// from the accounts table, otherwise will return nothing.
func (db *DB) GetAccount(accountID int) ([]byte, error) {
	row := db.QueryRow("SELECT * FROM accounts WHERE id = ?", accountID)

	account := new(Account)
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.AccountType,
		&account.MinimumPayment,
		&account.CurrentPayment,
		&account.FullAmount,
		&account.DueDate,
		&account.URL)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	return accountJSON, nil
}

// CreateAccount creates an account in the database and returns the account in JSON in the response
func (db *DB) CreateAccount(a Account) ([]byte, error) {
	result, err := db.Exec(`
		INSERT INTO accounts (user_id, name, account_type, minimum_payment, current_payment, full_amount, due_date, url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.UserID,
		a.Name,
		a.AccountType,
		a.MinimumPayment,
		a.CurrentPayment,
		a.FullAmount,
		a.DueDate,
		a.URL,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	accountJSON, err := db.GetAccount(int(id))
	if err != nil {
		return nil, err
	}

	return accountJSON, nil
}
