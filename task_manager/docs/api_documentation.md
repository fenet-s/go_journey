Task Management REST API

A simple REST API built with Go + Gin for managing tasks. Data is stored in memory.

Run

go mod tidy
go run .

Server:

http://localhost:8080

Task Format

{
"id": 1,
"title": "Learn Go",
"description": "Practice Gin",
"due_date": "2026-08-20",
"status": "Pending"
}

Allowed statuses:

Pending

In Progress

Completed

Endpoints

GET /tasks

Returns all tasks.

Response: 200 OK

GET /tasks/:id

Returns one task.

Example:

GET /tasks/1

Responses: 200 OK, 400 Bad Request, 404 Not Found

POST /tasks

Creates a task.

Body:

{
"title": "Learn Gin",
"description": "Build a REST API",
"due_date": "2026-08-25",
"status": "Pending"
}

Response: 201 Created

PUT /tasks/:id

Updates a task.

Body:

{
"title": "Learn Advanced Gin",
"description": "Improve the API",
"due_date": "2026-08-30",
"status": "In Progress"
}

Responses: 200 OK, 400 Bad Request, 404 Not Found

DELETE /tasks/:id

Deletes a task.

Response: 200 OK

{
"message": "task deleted successfully"
}

Postman

Create a collection named Task Management API with:

GET /tasks
GET /tasks/:id
POST /tasks
PUT /tasks/:id
DELETE /tasks/:id

Use:

{{baseUrl}} = http://localhost:8080

Test successful requests and error cases such as invalid IDs and missing tasks.

Data is lost when the server restarts because this project uses in-memory storage.
