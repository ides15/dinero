package config

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

var (
	// Recovery is a middleware that gracefully recovers from panics with a 500 code
	Recovery = negroni.NewRecovery()

	// RouteLogger is a middleware that logs each server route
	RouteLogger = &negroni.Logger{ALogger: log.New()}
)

func init() {
	Recovery.PrintStack = false

	RouteLogger.SetDateFormat(time.RFC1123)
	RouteLogger.SetFormat("{{.StartTime}} --- ({{.Status}}) {{.Method}} {{.Path}}")
}

// Middlewares returns the middlewares that the server will use
func Middlewares() []negroni.Handler {
	return []negroni.Handler{
		Recovery,
		RouteLogger,
	}
}
