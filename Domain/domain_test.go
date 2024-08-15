package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	task := Task{
		Title:   "Task 1",
		DueDate: time.Now().Add(24 * time.Hour),
		Status:  "pending",
	}

	assert.Equal(t, "Task 1", task.Title)
	assert.Equal(t, time.Now().Add(24*time.Hour).Day(), task.DueDate.Day())
	assert.Equal(t, "pending", task.Status)
}
func TestUser(t *testing.T) {
	user := User{
		Username: "testuser",
		Password: "hashedPassword",
	}

	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashedPassword", user.Password)
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected string
	}{
		{
			name: "valid pending task with future due date",
			task: Task{
				Title:   "Task 1",
				DueDate: time.Now().Add(24 * time.Hour),
				Status:  "pending",
			},
			expected: "",
		},
		{
			name: "valid completed task with past due date",
			task: Task{
				Title:   "Task 2",
				DueDate: time.Now().Add(-24 * time.Hour),
				Status:  "completed",
			},
			expected: "",
		},
		{
			name: "missing title",
			task: Task{
				Title:   "",
				DueDate: time.Now().Add(24 * time.Hour),
				Status:  "pending",
			},
			expected: "title is required",
		},
		{
			name: "missing due date",
			task: Task{
				Title:   "Task 3",
				DueDate: time.Time{},
				Status:  "pending",
			},
			expected: "due date is required",
		},
		{
			name: "missing status",
			task: Task{
				Title:   "Task 4",
				DueDate: time.Now().Add(24 * time.Hour),
				Status:  "",
			},
			expected: "status is required",
		},
		{
			name: "invalid status",
			task: Task{
				Title:   "Task 5",
				DueDate: time.Now().Add(24 * time.Hour),
				Status:  "invalid",
			},
			expected: "status must be either pending or completed",
		},
		{
			name: "completed task with future due date",
			task: Task{
				Title:   "Task 6",
				DueDate: time.Now().Add(24 * time.Hour),
				Status:  "completed",
			},
			expected: "due date must be in the past",
		},
		{
			name: "pending task with past due date",
			task: Task{
				Title:   "Task 7",
				DueDate: time.Now().Add(-24 * time.Hour),
				Status:  "pending",
			},
			expected: "due date must be in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected)
			}
		})
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Message: "Resource not found"}
	assert.EqualError(t, err, "Resource not found")
}

func TestUserAlreadyExistsError(t *testing.T) {
	err := &UserAlreadyExistsError{Message: "User already exists"}
	assert.EqualError(t, err, "User already exists")
}

func TestUnauthorizedError(t *testing.T) {
	err := &UnauthorizedError{Message: "Unauthorized access"}
	assert.EqualError(t, err, "Unauthorized access")
}

func TestForbiddenError(t *testing.T) {
	err := &ForbiddenError{Message: "Forbidden access"}
	assert.EqualError(t, err, "Forbidden access")
}

func TestInternalServerError(t *testing.T) {
	err := &InternalServerError{Message: "Internal server error"}
	assert.EqualError(t, err, "Internal server error")
}

func TestBadRequestError(t *testing.T) {
	err := &BadRequestError{Message: "Bad request"}
	assert.EqualError(t, err, "Bad request")
}