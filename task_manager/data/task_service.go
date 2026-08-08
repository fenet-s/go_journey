package data

import (
	"context"
	"fmt"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskService struct {
	Collection *mongo.Collection
}

func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{Collection: collection}
}

func (s *TaskService) EnsureIndexes(ctx context.Context) error {
	_, err := s.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("task_id_unique"),
	})
	if err != nil {
		return fmt.Errorf("create task indexes: %w", err)
	}

	return nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	cursor, err := s.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	tasks := make([]models.Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}

	return tasks, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int) (models.Task, error) {
	var task models.Task
	if err := s.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&task); err != nil {
		return models.Task{}, fmt.Errorf("find task by id %d: %w", id, err)
	}

	return task, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	nextID, err := s.nextTaskID(ctx)
	if err != nil {
		return models.Task{}, err
	}

	task.ID = nextID
	if _, err := s.Collection.InsertOne(ctx, task); err != nil {
		return models.Task{}, fmt.Errorf("insert task: %w", err)
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, updatedTask models.Task) (models.Task, error) {
	updatedTask.ID = id

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"due_date":    updatedTask.DueDate,
			"status":      updatedTask.Status,
		},
	}

	var task models.Task
	err := s.Collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("update task %d: %w", id, err)
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	result, err := s.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("delete task %d: %w", id, err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *TaskService) nextTaskID(ctx context.Context) (int, error) {
	var lastTask models.Task
	err := s.Collection.FindOne(
		ctx,
		bson.D{},
		options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}}),
	).Decode(&lastTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}

		return 0, fmt.Errorf("find last task id: %w", err)
	}

	return lastTask.ID + 1, nil
}
