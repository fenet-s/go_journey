package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"task_manager/data"
	"task_manager/router"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using system environment")
	}

	uri := os.Getenv("MONGODB_URI")
	uri = strings.TrimSpace(uri)
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	databaseName = strings.TrimSpace(databaseName)
	if databaseName == "" {
		databaseName = "task_manager"
	}

	collectionName := os.Getenv("MONGODB_COLLECTION")
	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		collectionName = "tasks"
	}

	serverPort := os.Getenv("SERVER_PORT")
	serverPort = strings.TrimSpace(serverPort)
	if _, err := strconv.Atoi(serverPort); err != nil {
		serverPort = "8080"
	}
	if serverPort == "" {
		serverPort = "8080"
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to MongoDB")

	db := client.Database(databaseName)

	taskCollection := db.Collection(collectionName)
	userCollection := db.Collection("users")

	taskService := data.TaskService{
		Collection: taskCollection,
	}
	userService := data.UserService{
		Collection: userCollection,
	}
	if err := taskService.EnsureIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := userService.EnsureIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}

	r := router.SetupRouter(&taskService, &userService)

	if err := r.Run(":" + serverPort); err != nil {
		log.Fatal(err)
	}
}
