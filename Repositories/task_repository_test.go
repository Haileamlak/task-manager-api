package repositories

import (
	"context"
	"testing"
	"time"

	domain "task-manager/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TaskRepositoryTestSuite defines the test suite for UserRepository
type TaskRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	repo       TaskRepository
	collection string
}

// SetupSuite runs once before the test suite
func (suite *TaskRepositoryTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	suite.NoError(err)

	err = client.Ping(ctx, readpref.Primary())
	suite.NoError(err)

	suite.client = client
	suite.collection = "tasks_test"
	suite.db = client.Database("test_db")
	suite.repo = NewTaskRepository(suite.db, suite.collection)
}

// TearDownSuite runs once after the test suite
func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	// drop the database at the end
	err := suite.client.Database("test_db").Drop(context.Background())
	suite.NoError(err)
	
	err = suite.client.Disconnect(context.TODO())
	suite.NoError(err)
}

// SetupTest runs before each test
func (suite *TaskRepositoryTestSuite) SetupTest() {
	err := suite.db.Collection(suite.collection).Drop(context.TODO())
	suite.NoError(err)
}

// TestUserRepositorySuite runs the test suite
func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}

// TestCreateTask_Success tests the CreateTask method with valid input
func (suite *TaskRepositoryTestSuite) TestCreateTask() {
	task := domain.Task{
		Title: "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status: "pending",
	}

	err := suite.repo.CreateTask(task)
	assert.NoError(suite.T(), err)

	var result domain.Task
	err = suite.db.Collection(suite.collection).FindOne(context.TODO(), bson.M{"title": "Test Task"}).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Task", result.Title)
}

// TestGetTask_Success tests the GetTask method with valid input
func (suite *TaskRepositoryTestSuite) TestGetTask_Success() {
	task := domain.Task{
		Title: "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status: "pending",
	}

	insertResult, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), task)
	assert.NoError(suite.T(), err)

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	
	result, err := suite.repo.GetTask(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Task", result.Title)
}

// TestGetTask_InvalidId tests the GetTask method with invalid input
func (suite *TaskRepositoryTestSuite) TestGetTask_InvalidId() {
	_, err := suite.repo.GetTask("invalid")
	assert.Error(suite.T(), err)
}

func (suite *TaskRepositoryTestSuite) TestGetTask_NotFound() {
	_, err := suite.repo.GetTask(primitive.NewObjectID().Hex())
	assert.Error(suite.T(), err)
}

// TestGetTasks_Success tests the GetTasks method with valid input
func (suite *TaskRepositoryTestSuite) TestGetTasks_Success() {
	task := domain.Task{
		Title: "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status: "pending",
	}

	_, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), task)
	assert.NoError(suite.T(), err)

	tasks, err := suite.repo.GetTasks()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(tasks))
}

// TestGetTasks_NotFound tests the GetTasks method with no tasks
func (suite *TaskRepositoryTestSuite) TestGetTasks_NotFound() {
	tasks, err := suite.repo.GetTasks()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tasks)
}

// TestUpdateTask_Success tests the UpdateTask method with valid input
func (suite *TaskRepositoryTestSuite) TestUpdateTask_Success() {
	task := domain.Task{
		Title: "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status: "pending",
	}

	insertResult, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), task)
	assert.NoError(suite.T(), err)

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()

	newTask := domain.Task{
		Title: "Updated Task",
	}

	err = suite.repo.UpdateTask(id, newTask)
	assert.NoError(suite.T(), err)

	var result domain.Task
	err = suite.db.Collection(suite.collection).FindOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Task", result.Title)
}

// TestUpdateTask_InvalidId tests the UpdateTask method with invalid input
func (suite *TaskRepositoryTestSuite) TestUpdateTask_InvalidId() {
	err := suite.repo.UpdateTask("invalid", domain.Task{Title: "Test Task", DueDate: time.Now().Add(24 * time.Hour), Status: "pending"})
	assert.Error(suite.T(), err)
}

func (suite *TaskRepositoryTestSuite) TestUpdateTask_NotFound() {
	err := suite.repo.UpdateTask(primitive.NewObjectID().Hex(), domain.Task{Title: "Test Task", DueDate: time.Now().Add(24 * time.Hour), Status: "pending"})
	assert.Error(suite.T(), err)
}

// TestDeleteTask_Success tests the DeleteTask method with valid input
func (suite *TaskRepositoryTestSuite) TestDeleteTask_Success() {
	task := domain.Task{
		Title: "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
		Status: "pending",
	}

	insertResult, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), task)
	assert.NoError(suite.T(), err)

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()

	err = suite.repo.DeleteTask(id)
	assert.NoError(suite.T(), err)

	var result domain.Task
	err = suite.db.Collection(suite.collection).FindOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}).Decode(&result)
	assert.Error(suite.T(), err)
}

// TestDeleteTask_InvalidId tests the DeleteTask method with invalid input
func (suite *TaskRepositoryTestSuite) TestDeleteTask_InvalidId() {
	err := suite.repo.DeleteTask("invalid")
	assert.Error(suite.T(), err)
}

func (suite *TaskRepositoryTestSuite) TestDeleteTask_NotFound() {
	err := suite.repo.DeleteTask(primitive.NewObjectID().Hex())
	assert.Error(suite.T(), err)
}


