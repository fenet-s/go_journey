package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"task_manager/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	Collection *mongo.Collection
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUsername   = errors.New("username is required")
	ErrInvalidPassword   = errors.New("password is required")
	ErrUserNotFound      = errors.New("user not found")
)

func (us *UserService) EnsureIndexes(ctx context.Context) error {
	_, err := us.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("user_username_unique"),
	})
	if err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}

	return nil
}

func (us *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)
	user.Role = models.NormalizeRole(user.Role)

	if user.Username == "" {
		return models.User{}, ErrInvalidUsername
	}
	if user.Password == "" {
		return models.User{}, ErrInvalidPassword
	}

	var existingUser models.User
	err := us.Collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return models.User{}, ErrUserAlreadyExists
	}
	if err != mongo.ErrNoDocuments {
		return models.User{}, fmt.Errorf("error checking for existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("error hashing password: %w", err)
	}

	user.ID = uuid.NewString()
	user.Password = string(hashedPassword)

	if _, err := us.Collection.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.User{}, ErrUserAlreadyExists
		}
		return models.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (us *UserService) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrUserNotFound
	}

	var user models.User
	err := us.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	user.Role = models.NormalizeRole(user.Role)
	return &user, nil
}
