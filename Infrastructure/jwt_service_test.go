package infrastructure

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	jwtService JWTService
}

func (suite *JWTServiceTestSuite) SetupTest() {
	suite.jwtService = NewJWTService()
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}

func (suite *JWTServiceTestSuite) TestGenerateToken_Success() {
	tokenString, err := suite.jwtService.GenerateToken("testuser", "admin")

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)

	token, err := suite.jwtService.ValidateToken(tokenString)
	assert.NoError(suite.T(), err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", claims["user"])
	assert.Equal(suite.T(), "admin", claims["role"])
}

func (suite *JWTServiceTestSuite) TestGenerateToken_Expiration() {
	tokenString, err := suite.jwtService.GenerateToken("testuser", "admin")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenString)

	token, err := suite.jwtService.ValidateToken(tokenString)
	assert.NoError(suite.T(), err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(suite.T(), ok)

	expiration := claims["exp"].(float64)
	assert.True(suite.T(), expiration > float64(time.Now().Unix()))
}

func (suite *JWTServiceTestSuite) TestValidateToken_Success() {
	tokenString, _ := suite.jwtService.GenerateToken("testuser", "admin")

	token, err := suite.jwtService.ValidateToken(tokenString)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), token)
	assert.True(suite.T(), token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", claims["user"])
	assert.Equal(suite.T(), "admin", claims["role"])
}

func (suite *JWTServiceTestSuite) TestValidateToken_InvalidTokenFormat() {
	token, err := suite.jwtService.ValidateToken("invalid.token.string")

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token.Valid)
}