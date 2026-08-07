package services

import (
	"task_manager/data"
	"task_manager/models"
)

func GetAllTasks() []models.Task {
	return data.Tasks
}

func GetTaskByID(id int) (models.Task, bool) {
	for _, task := range data.Tasks {
		if task.ID == id {
			return task, true
		}
	}
	return models.Task{}, false
}

func AddTask(task models.Task) models.Task {
	nextID := 1
	for _, existingTask := range data.Tasks {
		if existingTask.ID >= nextID {
			nextID = existingTask.ID + 1
		}

	}
	task.ID = nextID
	data.Tasks = append(data.Tasks, task)
	return task

}

func IsValidStatus(status string) bool {
	return status == "Pending" ||
		status == "In Progress" ||
		status == "Completed"
}

func UpdateTask(id int, updatedTask models.Task) (models.Task, bool) {

	for i, task := range data.Tasks {

		if task.ID == id {

			updatedTask.ID = id

			data.Tasks[i] = updatedTask

			return updatedTask, true
		}
	}

	return models.Task{}, false
}

func DeleteTask(id int) bool {
	for i, task := range data.Tasks {
		if task.ID == id {
			data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)
			return true
		}
	}
	return false
}
