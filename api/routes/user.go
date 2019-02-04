package routes

import (
	"dinero/api/config"
	"dinero/api/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
)

// UsersHandler handles the GET request for all users
func UsersHandler(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

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

// UserHandler handles the GET and POST requests for specific users
func UserHandler(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Get the user ID from the request path
			userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/user/"))
			if err != nil {
				env.Log.Error(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		} else if r.Method == "POST" {
			// Read POST request body
			newUser, err := ioutil.ReadAll(r.Body)
			if err != nil {
				env.Log.Error(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			// Read request body into User object
			var user models.User
			err = json.Unmarshal(newUser, &user)
			if err != nil {
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
					env.Log.Error(err)
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
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}
