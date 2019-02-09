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
	logger := config.Log

	// Get database reference
	db, err := models.InitDB(dbName)
	if err != nil {
		logger.Error(err)
	}

	// Set up environment
	env := &config.Env{DB: db, Log: logger}

	// Register negroni middleware(s)
	n := negroni.New()

	// Getting all middlewares and using them
	mws := config.Middlewares()
	for _, mw := range mws {
		n.Use(mw)
	}

	// Register chi router
	r := routes.NewRouter(env)

	// Binding all middlewares to chi router
	n.UseHandler(r)

	// Serve
	port := ":3000"
	logger.WithField("port", port).Info("Serving...")
	logger.Fatal(http.ListenAndServe(port, n))
}
