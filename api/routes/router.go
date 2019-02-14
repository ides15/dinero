package routes

import (
	"dinero/api/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter sets up a chi Mux router
func NewRouter(env *config.Env) *chi.Mux {
	r := chi.NewRouter()

	// Middleware to log each route using Logrus
	r.Use(config.RouteLogger(env))
	// Middleware to recover gracefully from panics
	r.Use(middleware.Recoverer)

	// Define routes
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", AllAccounts(env))    // GET /accounts
		r.Post("/", CreateAccount(env)) // POST /accounts

		r.Route("/{accountID}", func(r chi.Router) {
			r.Use(AccountCtx(env))
			r.Get("/", GetAccount(env))       // GET /accounts/123
			r.Put("/", UpdateAccount(env))    // PUT /accounts/123
			r.Delete("/", DeleteAccount(env)) // DELETE /accounts/123
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", AllUsers(env))    // GET /users
		r.Post("/", CreateUser(env)) // POST /users

		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx(env))
			r.Get("/", GetUser(env))       // GET /users/123
			r.Put("/", UpdateUser(env))    // PUT /users/123
			r.Delete("/", DeleteUser(env)) // DELETE /users/123
		})
	})

	r.MethodNotAllowed(MethodNotAllowed(env))

	return r
}
