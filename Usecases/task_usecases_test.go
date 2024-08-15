package usecases

import (
	domain "task-manager/Domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)


type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) CreateTask(task domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetTask(id string) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(id string, task domain.Task) error {
	args := m.Called(id, task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type TaskUsecaseTestSuite struct {
	suite.Suite
	taskRepo *MockTaskRepository
	usecase  TaskUsecase
}

func (suite *TaskUsecaseTestSuite) SetupSuite() {
	suite.taskRepo = new(MockTaskRepository)
	suite.usecase = NewTaskUsecase(suite.taskRepo)
}

func (suite *TaskUsecaseTestSuite) TearDownSuite() {
	suite.taskRepo.AssertExpectations(suite.T())

}

func (suite *TaskUsecaseTestSuite) SetupTest() {
	suite.taskRepo.ExpectedCalls = nil
}

func (suite *TaskUsecaseTestSuite) TearDownTest() {
	suite.taskRepo.AssertExpectations(suite.T())
}

func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

func (suite *TaskUsecaseTestSuite) TestCreateTask() {
	task := domain.Task{
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	suite.taskRepo.On("GetTasks").Return([]domain.Task{}, nil)

	suite.taskRepo.On("CreateTask", task).Return(nil)

	err := suite.usecase.CreateTask(task)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseTestSuite) TestCreateTaskWithExistingTitle() {
	task := domain.Task{
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	tasks := []domain.Task{
		{
			Title:   "Test Task",
			DueDate: time.Now().Add(24 * time.Hour),
			Status:  "pending",
		},
	}

	suite.taskRepo.On("GetTasks").Return(tasks, nil)

	err := suite.usecase.CreateTask(task)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Task already exists", err.Error())
}

func (suite *TaskUsecaseTestSuite) TestCreateTaskWithInvalidTask() {
	task := domain.Task{
		Title:   "",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	err := suite.usecase.CreateTask(task)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "title is required", err.Error())
}

func (suite *TaskUsecaseTestSuite) TestGetTask() {
	task := domain.Task{
		ID:      "1",
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	suite.taskRepo.On("GetTask", "1").Return(task, nil)

	result, err := suite.usecase.GetTask("1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), task, result)
}

func (suite *TaskUsecaseTestSuite) TestGetTasks() {
	tasks := []domain.Task{
		{
			ID:      "1",
			Title:   "Test Task 1",
			DueDate: time.Now().Add(24 * time.Hour),
			Status:  "pending",
		},
		{
			ID:      "2",
			Title:   "Test Task 2",
			DueDate: time.Now().Add(48 * time.Hour),
			Status:  "completed",
		},
	}

	suite.taskRepo.On("GetTasks").Return(tasks, nil)

	result, err := suite.usecase.GetTasks()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tasks, result)
}

func (suite *TaskUsecaseTestSuite) TestUpdateTask() {
	task := domain.Task{
		ID:      "1",
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	suite.taskRepo.On("UpdateTask", "1", task).Return(nil)

	err := suite.usecase.UpdateTask("1", task)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseTestSuite) TestUpdateTask_InvalidTask(){
	task := domain.Task{
		ID:      "1",
		Title:   "",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	err := suite.usecase.UpdateTask("1", task)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "title is required", err.Error())
}

func (suite *TaskUsecaseTestSuite) TestDeleteTask() {
	suite.taskRepo.On("DeleteTask", "1").Return(nil)

	err := suite.usecase.DeleteTask("1")
	assert.NoError(suite.T(), err)
}