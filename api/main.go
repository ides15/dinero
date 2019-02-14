package main

import (
	"dinero/api/config"
	"dinero/api/models"
	"dinero/api/routes"
	"net/http"
)

const (
	dbName = "../dinero.db"
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

	// Register chi router
	r := routes.NewRouter(env)

	// Serve
	port := ":3000"
	logger.WithField("port", port).Info("Serving...")
	logger.Fatal(http.ListenAndServe(port, r))
}
