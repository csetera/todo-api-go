package testsupport

import (
	"todo-api-go/entities"

	"time"
)

// CollectIds collects the IDs from a list of ToDoItemEntity.
//
// It takes a parameter `items` of type []entities.ToDoItemEntity, which is a list of ToDoItemEntity objects.
// It returns a []uint, which is a list of IDs extracted from the ToDoItemEntity objects.
func CollectIds(items []entities.ToDoItemEntity) []uint {
	var ids []uint
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	return ids
}

// Parses the given date string and returns the corresponding time.Time value.
// This function ignores errors that may occur when parsing the date string.
//
// It takes a single parameter:
// - date: a string representing the date in the format "2006-01-02".
//
// It returns a time.Time value.
func ParseTestDate(date string) time.Time {
	parsed, _ := time.Parse("2006-01-02", date)
	return parsed
}
