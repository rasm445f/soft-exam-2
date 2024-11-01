package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"soft/db"
	"soft-exam-2/models"
)

// GetAllTodos godoc
// @Summary Get all todos
// @Description Fetches a list of all todos from the database
// @Tags todos
// @Produce application/json
// @Success 200 {array} models.Todo
// @Router /api/todos [get]
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := db.FetchAllTodos()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := json.Marshal(todos)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetTodo godoc
// @Summary Get todo
// @Description Fetches a todo based on the id from the database
// @Tags todos
// @Produce application/json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.Todo
// @Router /api/todo/{id} [get]
func GetTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := db.GetTodoById(id)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := json.Marshal(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete a todo based on the id from the database
// @Tags todos
// @Produce application/json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /api/todo/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := db.DeleteTodoById(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(res)
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo entry in the database
// @Tags todos
// @Accept  application/json
// @Produce application/json
// @Param todo body models.Todo true "Todo object"
// @Success 201 {object} models.Todo
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/todo [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.CreateTodo(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}