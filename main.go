package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

// App is application object to hold the dependencies
type App struct {
	DB *gorm.DB
}

// a global error variable to hold any error
var err error

func main() {
	// load .env file
	if err = godotenv.Load(".env"); err != nil {
		log.Error(err)
	}

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	// include calling method in the log
	log.SetReportCaller(true)

	db, err := initDatabase()
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		DB: db,
	}

	// http.ListenAndServe(":9000", router)
	server := &http.Server{
		Handler: app.routes(),
		Addr:    "127.0.0.1:9000",
		// good practice: always set timeout
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288, // 512Kb
	}

	// Run server in a goroutine, so it doesn't block main thread.
	// This is NOT needed if this is the last part of the main() function.
	//go func() {
	//	if err := server.ListenAndServe(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	log.Info("starting server on ", server.Addr)
	// any error returned by http.ListenAndServe() is always non-nil
	err = server.ListenAndServe()
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
