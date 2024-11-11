package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/oTuff/go-startkode/db/generated"
)

// GetAllTodos godoc
// @Summary Get all todos
// @Description Fetches a list of all todos from the database
// @Tags todos
// @Produce application/json
// @Success 200 {array} generated.Todo
// @Router /api/todos [get]
func GetAllTodos(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use request context
		ctx := r.Context()

		// Call the FetchAllTodos method with context
		todos, err := queries.FetchAllTodos(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Marshal and write response
		res, _ := json.Marshal(todos)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetTodo godoc
// @Summary Get todo
// @Description Fetches a todo based on the id from the database
// @Tags todos
// @Produce application/json
// @Param id path string true "Todo ID"
// @Success 200 {object} generated.Todo
// @Router /api/todo/{id} [get]
func GetTodo(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		todo, err := queries.GetTodoById(ctx, id)
		if err != nil {
			log.Fatal(err)
		}

		res, _ := json.Marshal(todo)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete a todo based on the id from the database
// @Tags todos
// @Produce application/json
// @Param id path string true "Todo ID"
// @Success 200 {object} generated.Todo
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /api/todo/{id} [delete]
func DeleteTodo(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		err = queries.DeleteTodoById(ctx, id)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.Write(res)
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo entry in the database
// @Tags todos
// @Accept  application/json
// @Produce application/json
// @Param todo body generated.CreateTodoParams true "Todo object"
// @Success 201 {object} generated.CreateTodoParams
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/todo [post]
func CreateTodo(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var todo generated.CreateTodoParams

		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = queries.CreateTodo(ctx, todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
