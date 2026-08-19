package usecases

import (
	"context"
	"errors"
	"refactored_task_management_api/domain"
	"refactored_task_management_api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
	taskRepo    *mocks.TaskRepositoryMock
	taskUsecase domain.TaskUsecase
}

func (s *TaskUsecaseTestSuite) SetupTest() {
	s.taskRepo = new(mocks.TaskRepositoryMock)
	s.taskUsecase = NewTaskUsecase(s.taskRepo)
}

func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

func (s *TaskUsecaseTestSuite) TestCreateTask_Success() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}
	s.taskRepo.On("Create", mock.Anything, task).Return(nil)

	err := s.taskUsecase.CreateTask(context.Background(), task)
	assert.NoError(s.T(), err)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestCreateTask_EmptyTitle() {
	task := &domain.Task{
		Title:       "",
		Description: "This is a test task",
		Status:      "pending",
	}
	err := s.taskUsecase.CreateTask(context.Background(), task)

	assert.ErrorIs(s.T(), err, domain.ErrEmptyTitle)

	s.taskRepo.AssertNotCalled(s.T(),
		"Create", mock.Anything, mock.Anything)

}

func (s *TaskUsecaseTestSuite) TestCreateTask_DefaultStatus() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "This is a test task"}
	s.taskRepo.On("Create", mock.Anything, task).Return(nil)

	err := s.taskUsecase.CreateTask(context.Background(), task)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "pending", task.Status)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestCreateTask_RepositoryError() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}

	expectedErr := errors.New("database error")

	s.taskRepo.On("Create", mock.Anything, task).Return(expectedErr)

	err := s.taskUsecase.CreateTask(context.Background(), task)
	assert.ErrorIs(s.T(), err, expectedErr)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestGetTask_AdminAccess() {
	taskID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()

	task := &domain.Task{
		ID:     taskID,
		UserID: ownerID,
		Title:  "Admin Task",
	}

	s.taskRepo.
		On("GetByID", mock.Anything, taskID).
		Return(task, nil)

	result, err := s.taskUsecase.GetTask(
		context.Background(),
		taskID,
		"some-requester",
		domain.RoleAdmin,
	)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, result)

	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestGetTask_Unauthorized() {
	taskID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	task := &domain.Task{
		ID:     taskID,
		UserID: ownerID,
		Title:  "Private Task",
	}

	s.taskRepo.
		On("GetByID", mock.Anything, taskID).
		Return(task, nil)

	result, err := s.taskUsecase.GetTask(
		context.Background(),
		taskID,
		otherUserID.Hex(),
		"user",
	)

	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, domain.ErrUnauthorized)

	s.taskRepo.AssertExpectations(s.T())
}
