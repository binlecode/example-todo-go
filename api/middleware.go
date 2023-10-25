package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (app *App) BasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("action", "BasicMiddleware")
		log.Info("basic middleware called on ", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// StaticFsMiddleware is a middleware to static file server Routes
func (app *App) StaticFsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("static file server middleware called on ", r.URL.Path)
		// if valid, call next handler function
		next.ServeHTTP(w, r)
	})
}
