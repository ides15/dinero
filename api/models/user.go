package models

import (
	"database/sql"
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
func (db *DB) AllUsers() ([]*User, error) {
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

	return users, nil
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
func (db *DB) GetUser(userID int) (*User, error) {
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

	return user, nil
}

// CreateUser creates a user in the database and returns the user in JSON in the response
func (db *DB) CreateUser(u User) (*User, error) {
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

	user, err := db.GetUser(int(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a full resource in the database and returns an error if something goes wrong
func (db *DB) UpdateUser(userID int, u User) error {
	_, err := db.Exec(`
		UPDATE users
		SET
			first_name = ?,
			last_name = ?,
			full_name = ?,
			email = ?,
			biweekly_income = ?
		WHERE id = ?`,
		u.FirstName,
		u.LastName,
		u.FullName,
		u.Email,
		u.BiweeklyIncome,
		userID)

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
