package models

import (
	"database/sql"
	"encoding/json"
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
