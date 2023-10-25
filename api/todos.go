package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func (app *App) ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	log := log.WithField("action", "ListTodosHandler")
	params := mux.Vars(r)
	log = log.WithField("params", params)

	var todos []Todo
	err := app.DB.Find(&todos).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("finding todos, total: ", len(todos))
	respondWithJSON(w, http.StatusOK, todos)
}

func (app *App) GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log := log.WithField("params", params)
	id, _ := strconv.Atoi(params["id"])

	var todo Todo
	err := app.DB.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Infof("got todo: %v \n", todo)
	respondWithJSON(w, http.StatusOK, todo)
}

func (app *App) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// title := r.FormValue("title")
	// log.WithFields(log.Fields{"title": title}).Info("add new todo")
	// todo := &Todo{Title: title, Completed: false}

	rBodyJson, _ := io.ReadAll(r.Body)
	// log.Info("creating todo: ", rBodyJson)
	var todo Todo
	err := json.Unmarshal(rBodyJson, &todo)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := app.DB.Create(&todo)
	if result == nil {
		log.Error(fmt.Sprintf("failed to create todo: %+v", todo))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debug(fmt.Sprintf("created new todo: %+v", todo))
	log.Debug("db rows affected:", result.RowsAffected)
	log.WithFields(log.Fields{"Id": todo.ID, "Completed": todo.Completed}).Info("Creating todo")
	respondWithJSON(w, http.StatusCreated, todo)
}

func (app *App) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var todo Todo
	err := app.DB.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	title := r.FormValue("title")
	todo.Completed = completed
	todo.Title = title

	log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating todo")

	err = app.DB.Save(&todo).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (app *App) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var todo Todo
	err := app.DB.First(&todo, id).Error // gorm.ErrRecordNotFound if not found
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{"Id": id}).Info("Deleting todo")

	err = app.DB.Delete(&todo).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
