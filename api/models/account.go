package models

import (
	"database/sql"
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
func (db *DB) AllAccounts() ([]*Account, error) {
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

	return accounts, nil
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
func (db *DB) GetAccount(accountID int) (*Account, error) {
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

	return account, nil
}

// CreateAccount creates an account in the database and returns the account in JSON in the response
func (db *DB) CreateAccount(a Account) (*Account, error) {
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

	account, err := db.GetAccount(int(id))
	if err != nil {
		return nil, err
	}

	return account, nil
}

// UpdateAccount updates a full resource in the database and returns an error if something goes wrong
func (db *DB) UpdateAccount(accountID int, a *Account) error {
	_, err := db.Exec(`
		UPDATE accounts
		SET
			user_id = ?,
			name = ?,
			account_type = ?,
			minimum_payment = ?,
			current_payment = ?,
			full_amount = ?,
			due_date = ?,
			url = ?
		WHERE id = ?`,
		a.UserID,
		a.Name,
		a.AccountType,
		a.MinimumPayment,
		a.CurrentPayment,
		a.FullAmount,
		a.DueDate,
		a.URL,
		accountID)

	if err != nil {
		return err
	}

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return err
	// }

	// rows, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// DeleteAccount removes a resource from the database and returns an error if something goes wrong
func (db *DB) DeleteAccount(userID int) error {
	result, err := db.Exec(`
		DELETE
		FROM accounts
		WHERE id = ?`,
		userID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return ErrNotFound
	}

	return nil
}
