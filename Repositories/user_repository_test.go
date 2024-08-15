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

// UserRepositoryTestSuite defines the test suite for UserRepository
type UserRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	repo       UserRepository
	collection string
}

// SetupSuite runs once before the test suite
func (suite *UserRepositoryTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	suite.NoError(err)

	err = client.Ping(ctx, readpref.Primary())
	suite.NoError(err)

	suite.client = client
	suite.collection = "users_test"
	suite.db = client.Database("test_db")
	suite.repo = NewUserRepository(suite.db, suite.collection)
}

// TearDownSuite runs once after the test suite
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// drop the database at the end
	err := suite.client.Database("test_db").Drop(context.Background())
	suite.NoError(err)
	
	err = suite.client.Disconnect(context.TODO())
	suite.NoError(err)
}

// SetupTest runs before each test
func (suite *UserRepositoryTestSuite) SetupTest() {
	err := suite.db.Collection(suite.collection).Drop(context.TODO())
	suite.NoError(err)
}

// TestUserRepositorySuite runs the test suite
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}


// TestCreateUser_Success tests the CreateUser method with valid input
func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
	}

	err := suite.repo.CreateUser(user)
	assert.NoError(suite.T(), err)

	// Verify user exists in the database
	var storedUser domain.User
	err = suite.db.Collection(suite.collection).FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&storedUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, storedUser.Username)
}

// TestUpdateUser_Success tests the UpdateUser method with valid input
func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:	 "user",
	}

	// Insert a user to update
	insertedResult, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), user)
	assert.NoError(suite.T(), err)

	id := insertedResult.InsertedID.(primitive.ObjectID).Hex()

	// Update the user
	updatedUser := domain.User{
		Username: "updateduser",
	}

	err = suite.repo.UpdateUser(id, updatedUser)
	assert.NoError(suite.T(), err)

	// Verify user was updated
	var storedUser domain.User
	err = suite.db.Collection(suite.collection).FindOne(context.TODO(), bson.M{"_id": insertedResult.InsertedID}).Decode(&storedUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "updateduser", storedUser.Username)
}


// TestFindByUsername_Success tests the FindByUsername method with a valid username
func (suite *UserRepositoryTestSuite) TestFindByUsername_Success() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
	}

	// Insert a user to find
	_, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), user)
	assert.NoError(suite.T(), err)

	// Find the user
	storedUser, err := suite.repo.FindByUsername(user.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, storedUser.Username)
}

// TestFindByUsername_NotFound tests the FindByUsername method when the user is not found
func (suite *UserRepositoryTestSuite) TestFindByUsername_NotFound() {
	_, err := suite.repo.FindByUsername("nonexistentuser")
	assert.Error(suite.T(), err)
}

// TestCountUsers_Success tests the CountUsers method
func (suite *UserRepositoryTestSuite) TestCountUsers_Success() {
	user1 := domain.User{
		Username: "user1",
		Password: "password123",
	}
	user2 := domain.User{
		Username: "user2",
		Password: "password123",
	}

	_, err := suite.db.Collection(suite.collection).InsertOne(context.TODO(), user1)
	assert.NoError(suite.T(), err)
	_, err = suite.db.Collection(suite.collection).InsertOne(context.TODO(), user2)
	assert.NoError(suite.T(), err)

	count, err := suite.repo.CountUsers()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
}

// TestCountUsers_Empty tests the CountUsers method when there are no users
func (suite *UserRepositoryTestSuite) TestCountUsers_Empty() {
	count, err := suite.repo.CountUsers()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}
