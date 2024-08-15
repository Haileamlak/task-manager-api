package usecases

import (
	// "errors"
	"testing"
	// "time"

	domain "task-manager/Domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mocked dependencies
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(id string, user domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePasswords(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(username, role string) (string, error) {
	args := m.Called(username, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

// UserUsecaseTestSuite defines the test suite for UserUsecase
type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo        *MockUserRepository
	passwordService *MockPasswordService
	jwtService      *MockJWTService
	usecase         UserUsecase
}

// SetupTest runs before the test runs
func (suite *UserUsecaseTestSuite) SetupSuite() {
	suite.userRepo = new(MockUserRepository)
	suite.passwordService = new(MockPasswordService)
	suite.jwtService = new(MockJWTService)
	suite.usecase = NewUserUsecase(suite.userRepo, suite.passwordService, suite.jwtService)
}

func (suite *UserUsecaseTestSuite) TearDownSuite() {
	suite.userRepo.AssertExpectations(suite.T())
	suite.passwordService.AssertExpectations(suite.T())
	suite.jwtService.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo.ExpectedCalls = nil
	suite.passwordService.ExpectedCalls = nil
	suite.jwtService.ExpectedCalls = nil
}

func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.userRepo.AssertExpectations(suite.T())
	suite.passwordService.AssertExpectations(suite.T())
	suite.jwtService.AssertExpectations(suite.T())
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

// TestRegister_Success tests the Register method with valid input
func (suite *UserUsecaseTestSuite) TestRegister_Success() {
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, &domain.NotFoundError{})
	suite.passwordService.On("HashPassword", password).Return(hashedPassword, nil)
	suite.userRepo.On("CountUsers").Return(int64(0), nil)
	suite.userRepo.On("CreateUser", mock.AnythingOfType("domain.User")).Return(nil)

	err := suite.usecase.Register(username, password)
	assert.NoError(suite.T(), err)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "HashPassword", password)
	suite.userRepo.AssertCalled(suite.T(), "CountUsers")
	suite.userRepo.AssertCalled(suite.T(), "CreateUser", mock.AnythingOfType("domain.User"))
}

// TestRegister_ExistingUser tests the Register method when the username already exists
func (suite *UserUsecaseTestSuite) TestRegister_ExistingUser() {
	username := "testuser"
	password := "password123"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, nil)

	err := suite.usecase.Register(username, password)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "username already exists", err.Error())

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
}

// TestRegister_EmptyUsername tests the Register method with an empty username
func (suite *UserUsecaseTestSuite) TestRegister_EmptyUsernameAndPassword() {
	username := ""
	password := ""

	err := suite.usecase.Register(username, password)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "username and password are required", err.Error())
}

// TestRegister_Error tests the Register method when an error occurs
func (suite *UserUsecaseTestSuite) TestRegister_CountError() {
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, &domain.NotFoundError{})
	suite.passwordService.On("HashPassword", password).Return(hashedPassword, nil)
	suite.userRepo.On("CountUsers").Return(int64(0), &domain.InternalServerError{})

	err := suite.usecase.Register(username, password)
	assert.Error(suite.T(), err)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "HashPassword", password)
	suite.userRepo.AssertCalled(suite.T(), "CountUsers")
}

func (suite *UserUsecaseTestSuite) TestRegister_HashError() {
	username := "testuser"
	password := "password123"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, &domain.NotFoundError{})
	suite.passwordService.On("HashPassword", password).Return("", &domain.InternalServerError{})

	err := suite.usecase.Register(username, password)
	assert.Error(suite.T(), err)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "HashPassword", password)
}

// TestLogin_Success tests the Login method with valid input
func (suite *UserUsecaseTestSuite) TestLogin_Success() {
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"
	token := "token"

	user := domain.User{
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}

	suite.userRepo.On("FindByUsername", username).Return(user, nil)
	suite.passwordService.On("ComparePasswords", hashedPassword, password).Return(nil)
	suite.jwtService.On("GenerateToken", username, user.Role).Return(token, nil)

	t, err := suite.usecase.Login(username, password)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), token, t)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "ComparePasswords", hashedPassword, password)
	suite.jwtService.AssertCalled(suite.T(), "GenerateToken", username, user.Role)
}

// TestLogin_UserNotFound tests the Login method when the user is not found
func (suite *UserUsecaseTestSuite) TestLogin_UserNotFound() {
	username := "testuser"
	password := "password123"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, &domain.NotFoundError{})

	_, err := suite.usecase.Login(username, password)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "invalid username or password", err.Error())

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
}

// TestLogin_PasswordMismatch tests the Login method when the password does not match
func (suite *UserUsecaseTestSuite) TestLogin_PasswordMismatch() {
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"

	user := domain.User{
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}

	suite.userRepo.On("FindByUsername", username).Return(user, nil)
	suite.passwordService.On("ComparePasswords", hashedPassword, password).Return(&domain.BadRequestError{})

	_, err := suite.usecase.Login(username, password)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "invalid username or password", err.Error())

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "ComparePasswords", hashedPassword, password)
}

// TestLogin_Error tests the Login method when an error occurs
func (suite *UserUsecaseTestSuite) TestLogin_GenerateTokenError() {
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"

	user := domain.User{
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}

	suite.userRepo.On("FindByUsername", username).Return(user, nil)
	suite.passwordService.On("ComparePasswords", hashedPassword, password).Return(nil)
	suite.jwtService.On("GenerateToken", username, user.Role).Return("", &domain.InternalServerError{})

	_, err := suite.usecase.Login(username, password)
	assert.Error(suite.T(), err)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.passwordService.AssertCalled(suite.T(), "ComparePasswords", hashedPassword, password)
	suite.jwtService.AssertCalled(suite.T(), "GenerateToken", username, user.Role)
}

// TestPromoteUser_Success tests the PromoteUser method with valid input
func (suite *UserUsecaseTestSuite) TestPromoteUser_Success() {
	username := "testuser"
	user := domain.User{
		ID:       "test_id",
		Username: username,
		Role:     "user",
	}

	suite.userRepo.On("FindByUsername", username).Return(user, nil)

	user.Role = "admin"
	suite.userRepo.On("UpdateUser", user.ID, user).Return(nil)

	err := suite.usecase.PromoteUser(username)
	assert.NoError(suite.T(), err)

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
	suite.userRepo.AssertCalled(suite.T(), "UpdateUser", user.ID, user)
}

// TestPromoteUser_UserNotFound tests the PromoteUser method when the user is not found
func (suite *UserUsecaseTestSuite) TestPromoteUser_UserNotFound() {
	username := "testuser"

	suite.userRepo.On("FindByUsername", username).Return(domain.User{}, &domain.NotFoundError{Message: "user not found"})

	err := suite.usecase.PromoteUser(username)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
}

// TestPromoteUser_AlreadyAdmin tests the PromoteUser method when the user is already an admin
func (suite *UserUsecaseTestSuite) TestPromoteUser_AlreadyAdmin() {
	username := "testuser"
	user := domain.User{
		Username: username,
		Role:     "admin",
	}

	suite.userRepo.On("FindByUsername", username).Return(user, nil)

	err := suite.usecase.PromoteUser(username)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user is already an admin", err.Error())

	suite.userRepo.AssertCalled(suite.T(), "FindByUsername", username)
}
