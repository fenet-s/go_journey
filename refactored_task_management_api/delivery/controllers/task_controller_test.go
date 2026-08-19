package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"refactored_task_management_api/domain"
	"refactored_task_management_api/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerTestSuite struct {
	suite.Suite
	taskUsecase *mocks.TaskUsecaseMock
	controller  *TaskController
}

func (s *TaskControllerTestSuite) SetupTest() {
	s.taskUsecase = new(mocks.TaskUsecaseMock)
	s.controller = NewTaskController(s.taskUsecase)
}

func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}

func (s *TaskControllerTestSuite) TestCreateTask_Success() {

	body := `{
	"title":"Learn Go",
	"description":"Learn Go programming language",
	"status":"pending",
	"due_date":"2024-12-31T23:59:59Z"
	}`

	req := httptest.NewRequest(
		http.MethodPost, "/tasks", strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = req
	c.Set("userID", "60b8d295f1d2c4a3b8e4d2c1") // Mock user ID

	s.taskUsecase.
		On(
			"CreateTask",
			mock.Anything,
			mock.AnythingOfType("*domain.Task"),
		).
		Return(nil)
	s.controller.CreateTask(c)
	s.Equal(http.StatusCreated, w.Code)
	s.T().Log("Response:", w.Body.String())
	s.taskUsecase.AssertExpectations(s.T())

}

func (s *TaskControllerTestSuite) TestCreateTask_InvalidJSON() {

	body := `{
	"title":"Learn Go",
	"description":"Learn Go programming language",
	
	}`
	// Missing closing brace
	req := httptest.NewRequest(
		http.MethodPost, "/tasks", strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", "60b8d295f1d2c4a3b8e4d2c1") // Mock user ID
	s.controller.CreateTask(c)
	s.Equal(http.StatusBadRequest, w.Code)
	s.T().Log("Response:", w.Body.String())
	s.taskUsecase.AssertNotCalled(
		s.T(),
		"CreateTask",
		mock.Anything,
		mock.Anything,
	)
}

func (s *TaskControllerTestSuite) TestCreateTask_EmptyTitle() {

	body := `{
	"description":"Learn Go programming language",
	"status":"pending",
	"due_date":"2024-12-31T23:59:59Z"
	}`

	req := httptest.NewRequest(
		http.MethodPost, "/tasks", strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", "60b8d295f1d2c4a3b8e4d2c1") // Mock user ID
	s.controller.CreateTask(c)

	s.Equal(http.StatusBadRequest, w.Code)
	s.taskUsecase.AssertNotCalled(
		s.T(),
		"CreateTask",
		mock.Anything,
		mock.Anything,
	)

}

func (s *TaskControllerTestSuite) TestCreateTask_UsecaseError() {

	body := `{
	"title":"Learn Go",
	"description":"Learn Go programming language",
	"status":"pending",
	"due_date":"2024-12-31T23:59:59Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", "60b8d295f1d2c4a3b8e4d2c1") // Mock user ID

	s.taskUsecase.
		On(
			"CreateTask",
			mock.Anything,
			mock.AnythingOfType("*domain.Task"),
		).
		Return(domain.ErrTaskNotFound)

	s.controller.CreateTask(c)

	s.NotEqual(http.StatusCreated, w.Code)
	s.T().Log("Response:", w.Body.String())
	s.taskUsecase.AssertExpectations(s.T())
}
