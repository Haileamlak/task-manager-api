package infrastructure

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware interface
type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	Authorize(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService JWTService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService JWTService) AuthMiddleware {
	return &authMiddleware{jwtService}
}

// Authenticate middleware
func (m *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			ctx.Abort()
			return
		}

		tokenString := authParts[1]
		token, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		
		ctx.Set("username", claims["username"])
		ctx.Set("role", claims["role"])

		ctx.Next()
	}
}

// Authorize middleware
func (m *authMiddleware) Authorize(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, _ := m.jwtService.ValidateToken(tokenString)
		claims, _ := token.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		if !contains(roles, role) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized for this action"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}