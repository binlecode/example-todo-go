package main

import (
	"context"
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

	staticFs := http.FileServer(http.Dir("staticFs"))
	//router.Handle("/staticFs/", staticFs)
	// add a middleware to log requests to staticFs
	router.PathPrefix("/static").Handler(basicMiddleware(http.StripPrefix("/static", staticFs)))

	router.HandleFunc("/health", HealthHandler).Methods("GET")

	router.Use(basicMiddleware)

	// http://localhost:9000/authorize
	router.HandleFunc("/authorize", AuthorizeHandler).Methods("POST")

	router.HandleFunc("/todos", ListTodosHandler).Methods("GET")
	router.HandleFunc("/todos/{id}", GetTodoHandler).Methods("GET")
	//router.HandleFunc("/todos", CreateTodoHandler).Methods("POST")
	router.Handle("/todos", tokenMiddleware(http.HandlerFunc(CreateTodoHandler))).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateTodoHandler).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodoHandler).Methods("DELETE")

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

// https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
// https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/#jwt-authentication-in-golang

func tokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check header for token, return 401 if not found or not valid
		token := r.Header.Get("Authorization")

		// validate token
		claims, err := ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Failed authorization"))
			return
		}

		// add claims to context
		log.Info("claims: ", claims)
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		// call next handler function
		next.ServeHTTP(w, r)
	})
}

func basicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("basic middleware called on ", r.URL.Path)
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
