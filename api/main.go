package main

import (
	"dinero/api/config"
	"dinero/api/models"
	"dinero/api/routes"
	"net/http"

	"github.com/urfave/negroni"
)

const (
	dbName = "./dinero.db"
)

func main() {
	log := config.Log

	// Get database reference
	db, err := models.InitDB(dbName)
	if err != nil {
		log.Error(err)
	}

	// Set up environment
	env := &config.Env{DB: db, Log: log}

	// Register negroni middleware(s)
	n := negroni.New()

	// Getting all middlewares and using them
	mws := config.Middlewares()
	for _, mw := range mws {
		n.Use(mw)
	}

	// Register mux
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/account/", routes.AccountHandler(env))
	mux.HandleFunc("/user", routes.UserHandler(env))
	mux.HandleFunc("/accounts", routes.AccountsHandler(env))
	mux.HandleFunc("/users", routes.UsersHandler(env))

	// Binding all middlewares to mux
	n.UseHandler(mux)

	// Serve
	port := ":3000"
	log.WithField("port", port).Info("Serving...")
	log.Fatal(http.ListenAndServe(port, n))
}
