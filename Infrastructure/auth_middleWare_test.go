package infrastructure

import (
	// "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/dgrijalva/jwt-go"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(username string, role string) (string, error) {
	args := m.Called(username, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

type AuthMiddlewareTestSuite struct {
	suite.Suite
	jwtService    *MockJWTService
	authMiddleware AuthMiddleware
	router        *gin.Engine
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.jwtService = new(MockJWTService)
	suite.authMiddleware = NewAuthMiddleware(suite.jwtService)
	suite.router = gin.Default()
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (suite *AuthMiddlewareTestSuite) TestAuthenticate_Success() {
	token := &jwt.Token{
		Valid: true,
		Claims: jwt.MapClaims{
			"username": "testuser",
			"role":     "user",
		},
	}
	suite.jwtService.On("ValidateToken", "valid_token").Return(token, nil)

	suite.router.Use(suite.authMiddleware.Authenticate())

	suite.router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Authenticated")
	suite.jwtService.AssertExpectations(suite.T())
}

// func (suite *AuthMiddlewareTestSuite) TestAuthenticate_InvalidToken() {
// 	suite.jwtService.On("ValidateToken", "invalid_token").Return(nil, errors.New("invalid token"))

// 	suite.router.Use(suite.authMiddleware.Authenticate())

// 	suite.router.GET("/test", func(ctx *gin.Context) {
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
// 	})

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/test", nil)
// 	req.Header.Set("Authorization", "Bearer invalid_token")
// 	suite.router.ServeHTTP(w, req)

// 	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
// 	assert.Contains(suite.T(), w.Body.String(), "invalid token")
// 	suite.jwtService.AssertExpectations(suite.T())
// }

func (suite *AuthMiddlewareTestSuite) TestAuthorize_Success() {
	token := &jwt.Token{
		Valid: true,
		Claims: jwt.MapClaims{
			"username": "testuser",
			"role":     "admin",
		},
	}
	suite.jwtService.On("ValidateToken", "valid_token").Return(token, nil)

	suite.router.Use(suite.authMiddleware.Authenticate())
	suite.router.Use(suite.authMiddleware.Authorize("admin"))

	suite.router.GET("/admin", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Authorized"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Authorized")
	suite.jwtService.AssertExpectations(suite.T())
}

func (suite *AuthMiddlewareTestSuite) TestAuthorize_Forbidden() {
	token := &jwt.Token{
		Valid: true,
		Claims: jwt.MapClaims{
			"username": "testuser",
			"role":     "user",
		},
	}
	suite.jwtService.On("ValidateToken", "valid_token").Return(token, nil)

	suite.router.Use(suite.authMiddleware.Authenticate())
	suite.router.Use(suite.authMiddleware.Authorize("admin"))

	suite.router.GET("/admin", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Authorized"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "You are not authorized for this action")
	suite.jwtService.AssertExpectations(suite.T())
}
