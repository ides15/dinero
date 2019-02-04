package config

import (
	"dinero/api/models"

	log "github.com/sirupsen/logrus"
)

// Env contains all environment dependencies for use in route handlers
type Env struct {
	DB  models.Store
	Log *log.Logger
}
