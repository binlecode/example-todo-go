package main

import (
	"github.com/binlecode/example-todo-go/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Error(err)
	}

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	// include calling method in the log
	log.SetReportCaller(true)

	db, err := api.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	app := api.App{
		DB: db,
	}

	serverAddr := api.GetEnv("SERVER_ADDR", "127.0.0.1:9000")

	// http.ListenAndServe(":9000", router)
	server := &http.Server{
		Handler: app.Routes(),
		Addr:    serverAddr,
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
