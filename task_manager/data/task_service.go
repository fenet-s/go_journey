package data

import "task_manager/models"

var Tasks = []models.Task{
	{
		ID:          1,
		Title:       "Learn Go",
		Description: "Complete Gin tutorial",
		DueDate:     "2026-08-10",
		Status:      "Pending",
	},
	{
		ID:          2,
		Title:       "Build API",
		Description: "Create Task Manager API",
		DueDate:     "2026-08-15",
		Status:      "In Progress",
	},
}

func IsValidStatus(status string) bool {
	return status == "Pending" ||
		status == "In Progress" ||
		status == "Completed"
}
