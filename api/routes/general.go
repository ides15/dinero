package routes

import (
	"dinero/api/config"
	"net/http"
)

// MethodNotAllowed is a route handler for catching requests in unallowed methods
func MethodNotAllowed(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
