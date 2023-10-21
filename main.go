package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func main() {
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	// include calling method in the log
	log.SetReportCaller(true)

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Error(err)
	}

	if err := initDatabase(); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting TodoList API server")
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
	srAuth.Handle("/userInfo", TokenMiddleware(http.HandlerFunc(UserInfoHandler))).Methods("GET")

	srTodos := router.PathPrefix("/todos").Subrouter()
	srTodos.Use(TokenMiddleware)
	srTodos.HandleFunc("", ListTodosHandler).Methods("GET")
	srTodos.HandleFunc("/{id}", GetTodoHandler).Methods("GET")
	srTodos.HandleFunc("", CreateTodoHandler).Methods("POST")
	srTodos.HandleFunc("/{id}", UpdateTodoHandler).Methods("PUT")
	srTodos.HandleFunc("/{id}", DeleteTodoHandler).Methods("DELETE")

	// http.ListenAndServe(":9000", router)
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:9000",
		// good practice: always set timeout
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run server in a goroutine, so it doesn't block main thread.
	// This is NOT needed if this is the last part of the main() function.
	//go func() {
	//	if err := server.ListenAndServe(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	// any error returned by http.ListenAndServe() is always non-nil
	err := server.ListenAndServe()
	log.Fatal(err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func basicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("basic middleware called on ", r.URL.Path)
		// if valid, call next handler function
		next.ServeHTTP(w, r)
	})
}

// StaticFsMiddleware is a middleware to static file server routes
func staticFsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("static file server middleware called on ", r.URL.Path)
		// if valid, call next handler function
		next.ServeHTTP(w, r)
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
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
