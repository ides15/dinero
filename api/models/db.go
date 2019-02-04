package models

import (
	"database/sql"

	// SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Store is a general interface for a datastore (real vs mock)
type Store interface {
	AllAccounts() ([]byte, error)
	GetAccount(int) ([]byte, error)
	// CreateAccount(struct {
	// 	UserID         int
	// 	Name           string
	// 	AccountType    string
	// 	MinimumPayment float64
	// 	CurrentPayment float64
	// 	FullAmount     *float64
	// 	DueDate        string
	// 	URL            *string
	// }) ([]byte, error)
	AllUsers() ([]byte, error)
	GetUser(int) ([]byte, error)
	CreateUser(User) ([]byte, error)
}

// DB is a general DB type for actual DB connections (vs mock DBs)
type DB struct {
	*sql.DB
}

// InitDB initializes a database
func InitDB(dbName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbName)
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
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "users" (
		"id" INTEGER,
		"first_name" TEXT NOT NULL,
		"last_name" TEXT NOT NULL,
		"full_name" TEXT NOT NULL,
		"email" TEXT NOT NULL,
		"biweekly_income" REAL NOT NULL,
		
		PRIMARY KEY("id")
	)`)
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
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "accounts" (
		"id" INTEGER,
		"user_id" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"account_type" TEXT NOT NULL,
		"minimum_payment" REAL NOT NULL,
		"current_payment" REAL NOT NULL,
		"full_amount" REAL NOT NULL,
		"due_date" TEXT NOT NULL,
		"url" TEXT NOT NULL,

		PRIMARY KEY("id")
	)`)
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
