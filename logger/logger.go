package logger

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// For creates a logger for a specific package and function.
func For(pkg string, fn string) *log.Entry {
	return log.
		WithField("pkg", pkg).
		WithField("fn", fn)
}
