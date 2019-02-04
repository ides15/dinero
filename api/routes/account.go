package routes

import (
	"dinero/api/config"
	"dinero/api/models"
	"net/http"
	"strconv"
	"strings"
)

// AccountsHandler handles the GET request for all accounts
func AccountsHandler(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

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

// AccountHandler handles the GET request for specific accounts
func AccountHandler(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		accountID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/account/"))
		if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := env.DB.GetAccount(accountID)
		if err == models.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(account)
	}
}
