package main

import (
	"context"
	"log"
	"os"
	"time"

	"refactored_task_management_api/delivery/controllers"
	"refactored_task_management_api/delivery/routers"
	"refactored_task_management_api/infrastructure"
	"refactored_task_management_api/repositories"
	"refactored_task_management_api/usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	mongoURI := envOr("MONGODB_URI", "mongodb://localhost:27017")
	dbName := envOr("DB_NAME", "task_manager")
	port := envOr("PORT", "8080")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping mongodb: %v", err)
	}
	db := client.Database(dbName)

	// Infrastructure
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService(jwtSecret, 24*time.Hour)

	// Repositories (depend on Domain interfaces, implemented against Mongo collections)
	taskRepo := repositories.NewTaskRepository(db.Collection("tasks"))
	userRepo := repositories.NewUserRepository(db.Collection("users"))

	// Usecases (depend only on Domain interfaces)
	taskUsecase := usecases.NewTaskUsecase(taskRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService)

	// Delivery
	taskController := controllers.NewTaskController(taskUsecase)
	userController := controllers.NewUserController(userUsecase)
	router := routers.SetupRouter(taskController, userController, jwtService)

	log.Printf("server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
