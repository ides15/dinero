package routes

import (
	"dinero/api/config"
	"dinero/api/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// AllAccounts gets all Account records within the accounts table in the database
func AllAccounts(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := env.DB.AllAccounts()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		accountsJSON, _ := json.Marshal(accounts)

		w.Header().Set("Content-Type", "application/json")
		w.Write(accountsJSON)
	}
}

// GetAccount gets an account from the database based on the Account ID in the URL and returns it
func GetAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.Atoi(vars["accountID"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := env.DB.GetAccount(accountID)
		if err == models.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		accountJSON, _ := json.Marshal(account)

		// Send the found account JSON back in the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(accountJSON)
		return
	}
}

// CreateAccount creates an account record in the database and returns that created record
func CreateAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read POST request body
		newAccount, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Read request body into Account object
		var account models.Account
		err = json.Unmarshal(newAccount, &account)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate Account fields
		valid := account.Validate()
		if !valid {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		// Create User in database
		createdAccount, err := env.DB.CreateAccount(account)
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		createdAccountJSON, _ := json.Marshal(createdAccount)

		// Send the created user JSON back in the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(createdAccountJSON)
		return
	}
}

// UpdateAccount updates an account record in the database and returns that created record
func UpdateAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.Atoi(vars["accountID"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Read PUT request body
		editedAccount, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Read request body into User object
		var newAccount models.Account
		err = json.Unmarshal(editedAccount, &newAccount)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate Account fields
		valid := newAccount.Validate()
		if !valid {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		// Check if Account is already in database and if not, create it
		_, err = env.DB.GetAccount(accountID)
		if err == models.ErrNotFound {
			_, err := env.DB.CreateAccount(newAccount)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
					return
				}
			} else if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Send a Status Created response
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			return
		}

		// Update user in database
		err = env.DB.UpdateAccount(accountID, &newAccount)
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Send a Status No Content response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// DeleteAccount deletes an account record in the database
func DeleteAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID, err := strconv.Atoi(vars["accountID"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = env.DB.DeleteAccount(accountID)
		if err == models.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
