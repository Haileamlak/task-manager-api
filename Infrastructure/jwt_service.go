package infrastructure

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTService interface
type JWTService interface {
	GenerateToken(username string, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string) JWTService {
	return &jwtService{secretKey: secret, issuer: "task-manager"}
}

// GenerateToken generates a new JWT token
func (s *jwtService) GenerateToken(username string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates a JWT token
func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check if the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(s.secretKey), nil
	})
}