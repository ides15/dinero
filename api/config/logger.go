package config

import (
	"github.com/sirupsen/logrus"
)

// Log is the default logger for the environment.
// You can add new logrus instances here for logging
// to different output sources.
var Log = logrus.New()

func init() {
	// Log.SetReportCaller(true)

	// Set Logger options here:
	// Log.Out
	// Log.Hooks
	// Log.Formatter
	// Log.ReportCaller
	// Log.Level
	// Log.ExitFunc
}
