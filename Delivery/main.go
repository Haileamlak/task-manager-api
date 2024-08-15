package main

import (
	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	infrastructure "task-manager/Infrastructure"
	usecases "task-manager/Usecases"
	repositories "task-manager/Repositories"
)

func main() {
	
	// Initialize services
	jwtService := infrastructure.NewJWTService()
	passwordService := infrastructure.NewPasswordService()
	
	// Initialize database
	databaseService := infrastructure.NewDatabase()
	db := databaseService.Connect("mongodb://localhost:27017")
	
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db, "users")
	taskRepo := repositories.NewTaskRepository(db, "tasks")
	// Initialize use cases
	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService)
	taskUsecase := usecases.NewTaskUsecase(taskRepo)

	// Initialize controllers
	apiController := controllers.NewApiController(taskUsecase, userUsecase)

	// Setup router
	r := routers.SetupRouter(apiController, jwtService)

	// Start the server
	if r.Run(":8080") != nil {
		panic("Failed to start server")
	}
}
