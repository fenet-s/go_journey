package mocks

import (
	"context"
	"refactored_task_management_api/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecaseMock struct {
	mock.Mock
}

func (m *TaskUsecaseMock) CreateTask(ctx context.Context, task *domain.Task) error {

	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskUsecaseMock) GetAllTasks(ctx context.Context, requesterID string, requesterRole string) ([]domain.Task, error) {
	args := m.Called(ctx, requesterID, requesterRole)
	var tasks []domain.Task
	if args.Get(0) != nil {
		tasks = args.Get(0).([]domain.Task)
	}
	return tasks, args.Error(1)
}

func (m *TaskUsecaseMock) GetTask(ctx context.Context, id primitive.ObjectID, requesterID string, requesterRole string) (*domain.Task, error) {
	args := m.Called(ctx, id, requesterID, requesterRole)
	var task *domain.Task
	if args.Get(0) != nil {
		task = args.Get(0).(*domain.Task)
	}
	return task, args.Error(1)
}

func (m *TaskUsecaseMock) UpdateTask(ctx context.Context, task *domain.Task, requesterID string, requesterRole string) error {
	args := m.Called(ctx, task, requesterID, requesterRole)
	return args.Error(0)
}

func (m *TaskUsecaseMock) DeleteTask(ctx context.Context, id primitive.ObjectID, requesterID string, requesterRole string) error {
	args := m.Called(ctx, id, requesterID, requesterRole)
	return args.Error(0)
}
