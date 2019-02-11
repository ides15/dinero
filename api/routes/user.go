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

// ContextUser is a wrapper for the string type to prevent reuse of context
// types from 3rd party libraries
type ContextUser string

// UserCtx provides a context for all user routes to have access to that user ID
func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userParam := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(userParam)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUser("userID"), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AllUsers gets all User records within the users tablein the database
func AllUsers(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := env.DB.AllUsers()
		if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(users)
	}
}

// GetUser gets a user from the databased based on the User ID in the URL and returns it
func GetUser(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(ContextUser("userID")).(int)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		// Get the user from the database based on the user ID
		user, err := env.DB.GetUser(userID)
		if err == models.ErrNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Send the found user JSON back in the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(user)
		return
	}
}

// CreateUser creates a user record in the database and returns that created record
func CreateUser(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read POST request body
		newUser, err := ioutil.ReadAll(r.Body)
		if err != nil {
			env.Log.Errorf("IOUTIL %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Read request body into User object
		var user models.User
		err = json.Unmarshal(newUser, &user)
		if err != nil {
			env.Log.Error(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate User fields
		valid := user.Validate()
		if valid != true {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Create User in database
		createdUser, err := env.DB.CreateUser(user)
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
		w.Write(createdUser)
		return
	}
}
