package model

import (
	"time"
)

// Task task
type Task struct {
	ID        uint      `json:"id"`         // id
	Title     string    `json:"title"`      // title
	Done      bool      `json:"done"`       // done
	Sort      int       `json:"sort"`       // sort
	DelFlag   bool      `json:"del_flag"`   // del_flag
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at
}
