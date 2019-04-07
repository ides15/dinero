package routes

import (
	"dinero/api/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewRouter sets up a chi Mux router
func NewRouter(env *config.Env) *mux.Router {
	r := mux.NewRouter()

	// // Middleware to log each route using Logrus
	r.Use(config.RouteLogger(env))
	// // Middleware to recover gracefully from panics
	r.Use(handlers.RecoveryHandler())

	// Define routes
	r.HandleFunc("/accounts", AllAccounts(env)).Methods("GET")
	r.HandleFunc("/accounts", CreateAccount(env)).Methods("POST")
	r.HandleFunc("/accounts/{accountID}", GetAccount(env)).Methods("GET")
	r.HandleFunc("/accounts/{accountID}", UpdateAccount(env)).Methods("PUT")
	r.HandleFunc("/accounts/{accountID}", DeleteAccount(env)).Methods("DELETE")

	r.HandleFunc("/users", AllUsers(env)).Methods("GET")
	r.HandleFunc("/users", CreateUser(env)).Methods("POST")
	r.HandleFunc("/users/{userID}", GetUser(env)).Methods("GET")
	r.HandleFunc("/users/{userID}", UpdateUser(env)).Methods("PUT")
	r.HandleFunc("/users/{userID}", DeleteUser(env)).Methods("DELETE")

	return r
}
