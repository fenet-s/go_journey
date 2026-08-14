package repositories

import (
	"context"
	"refactored_task_management_api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userDocument struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
}

func (d *userDocument) toDomain() *domain.User {
	return &domain.User{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Role:     d.Role,
	}
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	doc := &userDocument{
		ID:       primitive.NewObjectID(),
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return err
	}

	user.ID = doc.ID
	return nil
}

func (r *userRepository) FetchByEmail(ctx context.Context, email string) (*domain.User, error) {
	var doc userDocument
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return doc.toDomain(), nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
