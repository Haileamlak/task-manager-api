package main

import (
	"log"
	"os"
	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	infrastructure "task-manager/Infrastructure"
	repositories "task-manager/Repositories"
	usecases "task-manager/Usecases"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	mongoURI := os.Getenv("MONGO_URI")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Initialize services
	jwtService := infrastructure.NewJWTService(jwtSecret)
	passwordService := infrastructure.NewPasswordService()

	// Initialize database
	databaseService := infrastructure.NewDatabase()
	db := databaseService.Connect(mongoURI)

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
	if r.Run(":" + port) != nil {
		panic("Failed to start server")
	}
}
