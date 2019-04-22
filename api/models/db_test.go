package models

import (
	"database/sql"
	"fmt"
	"testing"
)

func must(d *DB, e error) *DB {
	return d
}

func TestInitDB(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		dbName   string
		expected error
	}{
		{
			// passes the test because the driver and dbName are both valid
			name:     "GOOD_DRIVER",
			driver:   "sqlite3",
			dbName:   ":memory:",
			expected: nil,
		},
		{
			// breaks the test because the driver 'bad' is an invalid sql driver
			name:     "BAD_DRIVER",
			driver:   "bad",
			dbName:   ":memory:",
			expected: fmt.Errorf("sql: unknown driver %q (forgotten import?)", "bad"),
		},
		{
			// breaks the test because the dbName address is invalid and conn can't be pinged
			name:     "BAD_PING_NONEXISTENT",
			driver:   "sqlite3",
			dbName:   "bad/path/db.db",
			expected: fmt.Errorf("unable to open database file"),
		},
		{
			// breaks the test because the dbName address is invalid and conn can't be pinged
			name:     "NO_PING_INVALID",
			driver:   "sqlite3",
			dbName:   ".",
			expected: fmt.Errorf("unable to open database file"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := InitDB(test.driver, test.dbName)
			if err != nil {
				if err.Error() != test.expected.Error() {
					t.Errorf("\nError\n\tExpected: \t%v\n\tGot: \t\t%v", test.expected.Error(), err.Error())
				}
			}
		})
	}
}

type sqliteMaster struct {
	resource  string
	name      string
	tableName string
	rootPage  int
	sql       string
}

func TestCreateUsersTable(t *testing.T) {
	db := must(InitDB("sqlite3", ":memory:"))

	err := createUsersTable(db.DB)
	if err != nil {
		t.Errorf("Error when initially creating the 'users' table: %v", err)
	}

	row := db.DB.QueryRow("SELECT * FROM sqlite_master WHERE name = 'users'")

	schema := new(sqliteMaster)
	err = row.Scan(
		&schema.resource,
		&schema.name,
		&schema.tableName,
		&schema.rootPage,
		&schema.sql,
	)

	if err == sql.ErrNoRows {
		t.Errorf("Received no rows when querying the 'users' table, meaning the 'users' table was not created.\n")
	} else if err != nil {
		t.Errorf("Received err when trying to create the 'users' table: %v\n", err)
	}

	if schema.sql != usersTableStmt {
		t.Errorf("Schema\n\tExpected: \n\t'%s'\n\tGot: \n\t'%s'", usersTableStmt, schema.sql)
	}
}

func TestCreateAccountsTable(t *testing.T) {
	db := must(InitDB("sqlite3", ":memory:"))

	err := createAccountsTable(db.DB)
	if err != nil {
		t.Errorf("Error when initially creating the 'accounts' table: %v", err)
	}

	row := db.DB.QueryRow("SELECT * FROM sqlite_master WHERE name = 'accounts'")

	schema := new(sqliteMaster)
	err = row.Scan(
		&schema.resource,
		&schema.name,
		&schema.tableName,
		&schema.rootPage,
		&schema.sql,
	)

	if err == sql.ErrNoRows {
		t.Errorf("Received no rows when querying the 'accounts' table, meaning the 'accounts' table was not created.\n")
	} else if err != nil {
		t.Errorf("Received err when trying to create the 'accounts' table: %v\n", err)
	}

	if schema.sql != accountsTableStmt {
		t.Errorf("Schema\n\tExpected: \n\t'%s'\n\tGot: \n\t'%s'", accountsTableStmt, schema.sql)
	}
}
