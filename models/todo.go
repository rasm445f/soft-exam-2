package models

import "time"

type Todo struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	Done     bool       `json:"done"`
	Category *string    `json:"category"`
	Deadline *time.Time `json:"deadline"`
}