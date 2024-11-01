package db

import (
	"database/sql"
	"log"

	"soft-exam-2/models"
)

func FetchAllTodos() ([]models.Todo, error) {
	todos := []models.Todo{}
	rows, err := DB.Query("SELECT id, title, text, isCompleted, category, deadline FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		var category sql.NullString
		var deadline sql.NullTime

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Body, &todo.Done, &category, &deadline)
		if err != nil {
			log.Fatal(err)
		}

		// Handle nullable fields (category and deadline)
		if category.Valid {
			todo.Category = &category.String
		} else {
			todo.Category = nil
		}

		if deadline.Valid {
			todo.Deadline = &deadline.Time
		} else {
			todo.Deadline = nil
		}

		todos = append(todos, todo)
	}
	if err != nil {
		log.Fatal(err)
	}
	return todos, err
}

func GetTodoById(id string) (models.Todo, error) {
	todo := models.Todo{}

	row := DB.QueryRow("SELECT id, title, text, isCompleted, category, deadline FROM todo WHERE id = $1", id)

	var category sql.NullString
	var deadline sql.NullTime

	err := row.Scan(&todo.ID, &todo.Title, &todo.Body, &todo.Done, &category, &deadline)
	if err != nil {
		log.Fatal(err)
	}

	// Handle nullable fields
	if category.Valid {
		todo.Category = &category.String
	} else {
		todo.Category = nil
	}

	if deadline.Valid {
		todo.Deadline = &deadline.Time
	} else {
		todo.Deadline = nil
	}
	return todo, err
}

func DeleteTodoById(id string) error {
	_, err := DB.Exec("DELETE FROM todo WHERE id=$1", id)
	return err
}

func CreateTodo(todo *models.Todo) error {
	var lastInsertId int
	query := `INSERT INTO todo (title, text, iscompleted, category, deadline)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := DB.QueryRow(query, todo.Title, todo.Body, todo.Done, todo.Category, todo.Deadline).Scan(&lastInsertId)
	return err
}