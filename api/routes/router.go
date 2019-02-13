package routes

import (
	"dinero/api/config"

	"github.com/go-chi/chi"
)

// NewRouter sets up a chi Mux router
func NewRouter(env *config.Env) *chi.Mux {
	r := chi.NewRouter()

	// Define routes
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", AllAccounts(env))    // GET /accounts
		r.Post("/", CreateAccount(env)) // POST /accounts

		r.Route("/{accountID}", func(r chi.Router) {
			r.Use(AccountCtx(env))
			r.Get("/", GetAccount(env))    // GET /accounts/123
			r.Put("/", UpdateAccount(env)) // PUT /accounts/123
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", AllUsers(env))    // GET /users
		r.Post("/", CreateUser(env)) // POST /users

		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx(env))
			r.Get("/", GetUser(env))    // GET /users/123
			r.Put("/", UpdateUser(env)) // PUT /users/123
		})
	})

	r.MethodNotAllowed(MethodNotAllowed(env))

	return r
}
