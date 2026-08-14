package usecases

import (
	"context"
	"refactored_task_management_api/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUsecase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
	}
}

func (tu *taskUsecase) CreateTask(ctx context.Context, task *domain.Task) error {
	if task.Title == "" {
		return domain.ErrEmptyTitle
	}
	if task.Status == "" {
		task.Status = "pending"
	}

	return tu.taskRepo.Create(ctx, task)
}

func (tu *taskUsecase) GetTask(ctx context.Context, id primitive.ObjectID, requesterID, requesterRole string) (*domain.Task, error) {
	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canAccess(task.UserID, requesterID, requesterRole) {
		return nil, domain.ErrUnauthorized
	}
	return task, nil
}

func (tu *taskUsecase) GetAllTasks(ctx context.Context, requesterID, requesterRole string) ([]domain.Task, error) {
	allTasks, err := tu.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if requesterRole == domain.RoleAdmin {
		return allTasks, nil
	}

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	userTasks := make([]domain.Task, 0)
	for _, task := range allTasks {
		if task.UserID == requesterObjectID {
			userTasks = append(userTasks, task)
		}
	}
	return userTasks, nil
}

func (tu *taskUsecase) UpdateTask(ctx context.Context, task *domain.Task, requesterID, requesterRole string) error {
	existingTask, err := tu.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}
	if !canAccess(existingTask.UserID, requesterID, requesterRole) {
		return domain.ErrUnauthorized
	}
	return tu.taskRepo.Update(ctx, task)
}

func (tu *taskUsecase) DeleteTask(ctx context.Context, id primitive.ObjectID, requesterID, requesterRole string) error {
	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !canAccess(task.UserID, requesterID, requesterRole) {
		return domain.ErrUnauthorized
	}
	return tu.taskRepo.Delete(ctx, id)
}

func canAccess(taskUserID primitive.ObjectID, requesterID string, requesterRole string) bool {
	if requesterRole == domain.RoleAdmin {
		return true
	}
	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return false
	}
	return taskUserID == requesterObjectID
}
