package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func (app *App) routes() *mux.Router {
	// StrictSlash(true) makes a redirect from '/abc' to '/abc/' or vise versa,
	// depending on whether the route is registered with a trailing slash or not.
	// Default is set to false, thus '/abc' and '/abc/' are treated differently
	// and no redirect is issued. In this case if a '/abc' is registered but
	// '/abc/' is requested, the NotFound handler is called and 404 is returned.
	router := mux.NewRouter().StrictSlash(false)

	// CORS
	cr := cors.New(cors.Options{
		//AllowedOrigins: []string{"*"},
		AllowedOrigins: []string{"http://127.0.0.1", "http://localhost:3000"},
		//AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
	router.Use(cr.Handler)

	// add a global middleware
	router.Use(basicMiddleware)

	// static file server
	staticFs := http.FileServer(http.Dir("staticFs"))
	// add a middleware to log requests to staticFs
	router.PathPrefix("/static").Handler(staticFsMiddleware(http.StripPrefix("/static", staticFs)))

	// health check route
	router.HandleFunc("/health", HealthHandler).Methods("GET")

	// auth routes
	srAuth := router.PathPrefix("/auth").Subrouter()
	srAuth.HandleFunc("/authorize", AuthorizeHandler).Methods("POST")
	srAuth.Handle("/userinfo", TokenMiddleware(http.HandlerFunc(UserinfoHandler))).Methods("GET")
	srAuth.Handle("/refresh", TokenMiddleware(http.HandlerFunc(RefreshHandler))).Methods("POST")

	srTodos := router.PathPrefix("/todos").Subrouter()
	srTodos.Use(TokenMiddleware)
	srTodos.HandleFunc("", ListTodosHandler).Methods("GET")
	srTodos.HandleFunc("/{id}", GetTodoHandler).Methods("GET")
	srTodos.HandleFunc("", CreateTodoHandler).Methods("POST")
	srTodos.HandleFunc("/{id}", UpdateTodoHandler).Methods("PUT")
	srTodos.HandleFunc("/{id}", DeleteTodoHandler).Methods("DELETE")
	return router
}
