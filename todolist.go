package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

type Todo struct {
	gorm.Model        // id, timestamping, and soft delete!
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
}

var db *gorm.DB
var err error

func main() {

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	// include calling method in the log
	log.SetReportCaller(true)

	if err := initDatabase(); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting TodoList API server")
	// StrictSlash(true) routes '/abc' to '/abc/'
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", HealthHandler).Methods("GET")
	router.HandleFunc("/todos/", ListTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", GetTodo).Methods("GET")
	router.HandleFunc("/todos/", CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")

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

func ListTodos(w http.ResponseWriter, r *http.Request) {
	log := log.WithField("action", "ListTodos")
	params := mux.Vars(r)
	log = log.WithField("params", params)

	var todos []Todo
	err := db.Find(&todos).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("finding todos, total: ", len(todos))
	respondWithJSON(w, http.StatusOK, todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log := log.WithField("params", params)
	id, _ := strconv.Atoi(params["id"])

	var todo Todo
	err := db.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Infof("got todo: %v \n", todo)
	respondWithJSON(w, http.StatusOK, todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	// title := r.FormValue("title")
	// log.WithFields(log.Fields{"title": title}).Info("add new todo")
	// todo := &Todo{Title: title, Completed: false}

	rBodyJson, _ := io.ReadAll(r.Body)
	// log.Info("creating todo: ", rBodyJson)
	var todo Todo
	json.Unmarshal(rBodyJson, &todo)
	result := db.Create(&todo)
	log.Debug("created new todo:", todo.ID)
	log.Debug("db rows affected:", result.RowsAffected)
	log.WithFields(log.Fields{"Id": todo.ID, "Completed": todo.Completed}).Info("Creating todo")
	respondWithJSON(w, http.StatusCreated, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var todo Todo
	err := db.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	title := r.FormValue("title")
	todo.Completed = completed
	todo.Title = title

	log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating todo")

	err = db.Save(&todo).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var todo Todo
	err := db.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{"Id": id}).Info("Deleting todo")

	err = db.Delete(&todo).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func initDatabase() error {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to sqlite3 database")
	}

	// create table and load data
	db.AutoMigrate(&Todo{})
	db.Create(&Todo{Title: "Test todo 1", Completed: false})
	db.Create(&Todo{Title: "Test todo 2", Completed: false})

	// sanity check
	var todo Todo
	db.First(&todo, 1)
	log.Infof("first loaded todo: %v \n", todo)

	return nil
}
