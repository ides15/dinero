package config

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// statusWriter extends http.ResponseWriter to get access to the http.Response.StatusCode
type statusWriter struct {
	http.ResponseWriter
	status int
}

// Log is the default logger for the environment.
// You can add new logrus instances here for logging
// to different outputs.
var (
	Log = logrus.New()
)

// Interface method for extending http.ResponseWriter
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Interface method for extending http.ResponseWriter
func (w *statusWriter) Write(b []byte) (int, error) {
	// defaults to 200
	if w.status == 0 {
		w.status = 200
	}

	n, err := w.ResponseWriter.Write(b)
	return n, err
}

// RouteLogger is a middleware to log all requests to the server using logrus
func RouteLogger(env *Env) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := statusWriter{ResponseWriter: w}

			next.ServeHTTP(&sw, r)

			env.Log.WithFields(logrus.Fields{
				"date":      start.Format(time.RFC1123),
				"duration":  time.Since(start),
				"method":    r.Method,
				"path":      r.URL.Path,
				"requester": r.RemoteAddr,
				"status":    sw.status,
			}).Info()
		})
	}
}

func init() {
	// // Logging in JSON
	// Log.SetFormatter(&logrus.JSONFormatter{})

	// // Logging the function where the log was called
	// Log.SetReportCaller(true)

	// // Sets the logging output to stdout
	// Log.SetOutput(os.Stdout)

	// // Sets the logging level to warn or higher
	// Log.SetLevel(logrus.WarnLevel)
}
