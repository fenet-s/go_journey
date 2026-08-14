package repositories

import (
	"context"
	"refactored_task_management_api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
	DueDate     time.Time          `bson:"due_date"`
	UserID      primitive.ObjectID `bson:"user_id"`
}

func (d *taskDocument) toDomain() *domain.Task {
	return &domain.Task{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
		DueDate:     d.DueDate,
		UserID:      d.UserID,
	}
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) domain.TaskRepository {
	return &taskRepository{
		collection: collection,
	}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	doc := &taskDocument{
		ID:          primitive.NewObjectID(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return err
	}

	task.ID = doc.ID
	return nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []taskDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0, len(docs))
	for _, doc := range docs {
		tasks = append(tasks, *doc.toDomain())
	}

	return tasks, nil
}

func (r *taskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var doc taskDocument
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	return doc.toDomain(), nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"due_date":    task.DueDate,
			"user_id":     task.UserID,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": task.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
