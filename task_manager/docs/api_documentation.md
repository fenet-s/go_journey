Task Management REST API

A simple REST API built with Go + Gin for managing tasks. Data is persisted in MongoDB.

## Configuration

Create a `.env` file in the project root or set environment variables directly.

Required:

- `MONGODB_URI` — MongoDB connection string, for example `mongodb://localhost:27017`

Optional:

- `MONGODB_DATABASE` — defaults to `task_manager`
- `MONGODB_COLLECTION` — defaults to `tasks`
- `SERVER_PORT` — defaults to `8080`

Example `.env`:

```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=task_manager
MONGODB_COLLECTION=tasks
SERVER_PORT=8080
```

## Run

```bash
go mod tidy
go run .
```

Server:

`http://localhost:8080`

## Task format

```json
{
  "id": 1,
  "title": "Learn Go",
  "description": "Practice Gin",
  "due_date": "2026-08-20",
  "status": "Pending"
}
```

Allowed statuses:

- `Pending`
- `In Progress`
- `Completed`

## Endpoints

### GET /tasks

Returns all tasks.

Response: `200 OK`

### GET /tasks/:id

Returns one task.

Example:

`GET /tasks/1`

Responses: `200 OK`, `400 Bad Request`, `404 Not Found`

### POST /tasks

Creates a task.

Body:

```json
{
  "title": "Learn Gin",
  "description": "Build a REST API",
  "due_date": "2026-08-25",
  "status": "Pending"
}
```

Response: `201 Created`

### PUT /tasks/:id

Updates a task.

Body:

```json
{
  "title": "Learn Advanced Gin",
  "description": "Improve the API",
  "due_date": "2026-08-30",
  "status": "In Progress"
}
```

Responses: `200 OK`, `400 Bad Request`, `404 Not Found`

### DELETE /tasks/:id

Deletes a task.

Response: `200 OK`

```json
{
  "message": "task deleted successfully"
}
```

## Postman verification

Create a collection named `Task Management API` with:

- `GET /tasks`
- `GET /tasks/:id`
- `POST /tasks`
- `PUT /tasks/:id`
- `DELETE /tasks/:id`

Use:

- `{{baseUrl}} = http://localhost:8080`

Test successful requests and error cases such as invalid IDs and missing tasks.

## MongoDB verification

After creating or updating a task, verify the document directly in MongoDB Compass or with a query such as:

```js
db.tasks.find().sort({ id: 1 });
```

Because tasks are stored in MongoDB, data remains available after API restarts.
