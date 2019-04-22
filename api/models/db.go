package models

import (
	"database/sql"

	// SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	usersTableStmt = `CREATE TABLE "users" (
		"id" INTEGER,
		"first_name" TEXT NOT NULL,
		"last_name" TEXT NOT NULL,
		"full_name" TEXT NOT NULL,
		"email" TEXT NOT NULL UNIQUE,
		"biweekly_income" REAL NOT NULL,
		
		PRIMARY KEY("id")
	)`
	accountsTableStmt = `CREATE TABLE "accounts" (
		"id" INTEGER,
		"user_id" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"account_type" TEXT NOT NULL,
		"minimum_payment" REAL NOT NULL,
		"current_payment" REAL NOT NULL,
		"full_amount" REAL NOT NULL,
		"due_date" TEXT NOT NULL,
		"url" TEXT NOT NULL,

		UNIQUE("user_id", "name")
		PRIMARY KEY("id")
	)`
)

// Store is a general interface for a datastore (real vs mock)
type Store interface {
	AllAccounts() ([]*Account, error)
	GetAccount(int) (*Account, error)
	CreateAccount(Account) (*Account, error)
	UpdateAccount(int, *Account) error
	DeleteAccount(int) error
	AllUsers() ([]*User, error)
	GetUser(int) (*User, error)
	CreateUser(User) (*User, error)
	UpdateUser(int, *User) error
	DeleteUser(int) error
}

// DB is a general DB type for actual DB connections (vs mock DBs)
type DB struct {
	*sql.DB
}

// InitDB initializes a database
func InitDB(driver string, name string) (*DB, error) {
	db, err := sql.Open(driver, name)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = createUsersTable(db)
	if err != nil {
		panic(err)
	}
	err = createAccountsTable(db)
	if err != nil {
		panic(err)
	}

	return &DB{db}, nil
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS 'users'")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(usersTableStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func createAccountsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS 'accounts'")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(accountsTableStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
