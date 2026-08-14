package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" binding:"required"`
	Description string             `json:"description" binding:"required"`
	DueDate     time.Time          `json:"due_date" binding:"required"`
	Status      string             `json:"status" binding:"required"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrEmptyTitle   = errors.New("title cannot be empty")
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Role     string             `json:"role" bson:"role"`
}

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type JWTService interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetAll(ctx context.Context) ([]Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FetchByEmail(ctx context.Context, email string) (*User, error)
	Count(ctx context.Context) (int64, error)
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context, id primitive.ObjectID, requesterID, requesterRole string) (*Task, error)
	GetAllTasks(ctx context.Context, requesterID, requesterRole string) ([]Task, error)
	UpdateTask(ctx context.Context, task *Task, requesterID, requesterRole string) error
	DeleteTask(ctx context.Context, id primitive.ObjectID, requesterID, requesterRole string) error
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (string, error)
}
