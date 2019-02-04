package models

import (
	"database/sql"
	"encoding/json"
	"regexp"
)

// User is a user of the applications
type User struct {
	ID             int     `json:"ID"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	FullName       string  `json:"fullName"`
	Email          string  `json:"email"`
	BiweeklyIncome float64 `json:"biweeklyIncome"`
}

// AllUsers retrieves all user rows from the users table
func (db *DB) AllUsers() ([]byte, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.FullName,
			&user.Email,
			&user.BiweeklyIncome)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return usersJSON, nil
}

// Validate validates the fields in a User object
func (u *User) Validate() bool {
	namePattern := regexp.MustCompile(`^[a-zA-Z ]+$`)
	emailPattern := regexp.MustCompile(`^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$`)

	if !namePattern.MatchString(u.FirstName) {
		return false
	}

	if !namePattern.MatchString(u.LastName) {
		return false
	}

	if !namePattern.MatchString(u.FullName) {
		return false
	}

	if !emailPattern.MatchString(u.Email) {
		return false
	}

	return true
}

// GetUser retrieves a user that matches the userID parameter
// from the users table, otherwise will return nothing.
func (db *DB) GetUser(userID int) ([]byte, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", userID)

	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.FullName,
		&user.Email,
		&user.BiweeklyIncome)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return userJSON, nil
}

// CreateUser creates a user in the database and returns the user in JSON in the response
func (db *DB) CreateUser(u User) ([]byte, error) {
	result, err := db.Exec(`
		INSERT INTO users (first_name, last_name, full_name, email, biweekly_income)
		VALUES (?, ?, ?, ?, ?)`,
		u.FirstName,
		u.LastName,
		u.FullName,
		u.Email,
		u.BiweeklyIncome)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	userJSON, err := db.GetUser(int(id))
	if err != nil {
		return nil, err
	}

	return userJSON, nil
}
