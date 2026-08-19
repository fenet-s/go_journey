package mocks

import (
	"context"
	"refactored_task_management_api/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepositoryMock struct {
	mock.Mock
}

func (m *TaskRepositoryMock) Create(ctx context.Context,
	task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepositoryMock) GetAll(ctx context.Context) ([]domain.Task, error) {
	args := m.Called(ctx)

	var tasks []domain.Task

	if args.Get(0) != nil {
		tasks = args.Get(0).([]domain.Task)
	}

	return tasks, args.Error(1)
}
func (m *TaskRepositoryMock) GetByID(
	ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	args := m.Called(ctx, id)

	var task *domain.Task

	if args.Get(0) != nil {
		task = args.Get(0).(*domain.Task)
	}
	return task, args.Error(1)
}
func (m *TaskRepositoryMock) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepositoryMock) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
