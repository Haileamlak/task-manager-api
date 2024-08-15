package routers

import (
	"task-manager/Delivery/controllers"
	infrastructure "task-manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(apiController controllers.ApiController, jwtService infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", apiController.Register)
	r.POST("/login", apiController.Login)

	// Protected routes
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)
	r.Use(authMiddleware.Authenticate())

	// All users routes
	r.GET("/tasks", apiController.GetTasks)
	r.GET("/tasks/:id", apiController.GetTask)

	adminAuthoriser := authMiddleware.Authorize("admin")

	// Admin only routes
	r.POST("/promote", adminAuthoriser, apiController.PromoteUser)
	r.POST("/tasks", adminAuthoriser, apiController.CreateTask)
	r.PUT("/tasks/:id", adminAuthoriser, apiController.UpdateTask)
	r.DELETE("/tasks/:id", adminAuthoriser, apiController.DeleteTask)

	return r
}
