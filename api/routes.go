package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io"
	"net/http"
)

func (app *App) Routes() *mux.Router {
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
	router.Use(app.BasicMiddleware)

	// static file server
	//staticFs := http.FileServer(http.Dir("staticFs"))
	// add a middleware to log requests to staticFs
	//router.PathPrefix("/static").Handler(app.StaticFsMiddleware(http.StripPrefix("/static", staticFs)))

	// health check route
	router.HandleFunc("/health", app.HealthHandler).Methods("GET")

	// auth Routes
	srAuth := router.PathPrefix("/auth").Subrouter()
	srAuth.HandleFunc("/authorize", AuthorizeHandler).Methods("POST")
	srAuth.Handle("/userinfo", TokenMiddleware(http.HandlerFunc(UserinfoHandler))).Methods("GET")
	srAuth.Handle("/refresh", TokenMiddleware(http.HandlerFunc(RefreshHandler))).Methods("POST")

	srTodos := router.PathPrefix("/todos").Subrouter()
	srTodos.Use(TokenMiddleware)
	srTodos.HandleFunc("", app.ListTodosHandler).Methods("GET")
	srTodos.HandleFunc("/{id}", app.GetTodoHandler).Methods("GET")
	srTodos.HandleFunc("", app.CreateTodoHandler).Methods("POST")
	srTodos.HandleFunc("/{id}", app.UpdateTodoHandler).Methods("PUT")
	srTodos.HandleFunc("/{id}", app.DeleteTodoHandler).Methods("DELETE")
	return router
}

func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		//w.Header().Set("Allow", "GET")
		w.Header().Set("Allow", http.MethodGet)
		//w.WriteHeader(405)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
