package routes

import (
	"context"
	"dinero/api/config"
	"dinero/api/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// ContextAccount is a wrapper for the string type to prevent reuse of context
// types from 3rd party libraries
type ContextAccount string

// AccountCtx provides a context for all account routes to have access to the account ID
func AccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountParam := chi.URLParam(r, "accountID")
		accountID, err := strconv.Atoi(accountParam)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAccount("accountID"), accountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AllAccounts gets all Account records within the accounts table in the database
func AllAccounts(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := env.DB.AllAccounts()
		if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(accounts)
	}
}

// GetAccount gets an account from the database based on the Account ID in the URL and returns it
func GetAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accountID, ok := ctx.Value(ContextAccount("accountID")).(int)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		// Get the account from the database based on the account ID
		account, err := env.DB.GetAccount(accountID)
		if err == models.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Send the found account JSON back in the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(account)
		return
	}
}

// CreateAccount creates an account record in the database and returns that created record
func CreateAccount(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read POST request body
		newAccount, err := ioutil.ReadAll(r.Body)
		if err != nil {
			env.Log.Errorf("IOUTIL: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Read request body into Account object
		var account models.Account
		err = json.Unmarshal(newAccount, &account)
		if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate Account fields
		valid := account.Validate()
		if valid != true {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Send the created user JSON back in the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(createdAccount)
		return
	}
}
