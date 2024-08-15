package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type PasswordServiceTestSuite struct {
	suite.Suite
	passwordService PasswordService
}

func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.passwordService = NewPasswordService()
}

func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}

func (suite *PasswordServiceTestSuite) TestHashPassword_Success() {
	password := "securePassword123"
	hashedPassword, err := suite.passwordService.HashPassword(password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)

	// Ensure that the hashed password can be verified against the original password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(suite.T(), err)
}

func (suite *PasswordServiceTestSuite) TestComparePasswords_Success() {
	password := "securePassword123"
	hashedPassword, _ := suite.passwordService.HashPassword(password)

	err := suite.passwordService.ComparePasswords(hashedPassword, password)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordServiceTestSuite) TestComparePasswords_Failure() {
	password := "securePassword123"
	hashedPassword, _ := suite.passwordService.HashPassword(password)

	// Attempt to compare with an incorrect password
	err := suite.passwordService.ComparePasswords(hashedPassword, "wrongPassword")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), bcrypt.ErrMismatchedHashAndPassword, err)
}
