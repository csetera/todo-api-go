package entities

import "time"

type ToDoItemEntity struct {
	ID          uint
	Description string
	Completed   bool
	DueDate     time.Time
	CompletedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
