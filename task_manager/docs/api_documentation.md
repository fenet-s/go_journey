Task Management REST API

A simple REST API built with Go + Gin for managing tasks with JWT authentication and role-based authorization. Data is persisted in MongoDB.

## Configuration

Create a `.env` file in the project root or set environment variables directly.

Required:

- `MONGODB_URI` — MongoDB connection string, for example `mongodb://localhost:27017`
- `JWT_SECRET` — secret used to sign and verify JWT tokens

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
JWT_SECRET=replace-with-a-long-random-secret
```

## Run

```bash
go mod tidy
go run .
```

Server:

`http://localhost:8080`

## Authentication flow

1. Register a user with `POST /register`.
2. Log in with `POST /login` to receive a JWT access token.
3. Send the token on protected requests as:

```http
Authorization: Bearer <token>
```

The token includes the user ID, username, role, and expiration time.

Role rules:

- `user` — can access protected task routes
- `admin` — can access protected task routes and delete tasks

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

### POST /register

Creates a new user account.

Body:

```json
{
  "username": "alice",
  "password": "P@ssw0rd123"
}
```

Response: `201 Created`

```json
{
  "message": "user registered successfully",
  "id": "...",
  "username": "alice",
  "role": "user"
}
```

### POST /login

Authenticates a user and returns a JWT token.

Body:

```json
{
  "username": "alice",
  "password": "P@ssw0rd123"
}
```

Response: `200 OK`

```json
{
  "message": "login successful",
  "token": "<jwt-token>",
  "role": "user"
}
```

### GET /tasks

Returns all tasks. Requires a valid JWT.

Response: `200 OK`

### GET /tasks/:id

Returns one task. Requires a valid JWT.

Example:

`GET /tasks/1`

Responses: `200 OK`, `400 Bad Request`, `404 Not Found`

### POST /tasks

Creates a task. Requires a valid JWT.

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

Updates a task. Requires a valid JWT.

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

Deletes a task. Requires a valid JWT and the `admin` role.

Response: `200 OK`

```json
{
  "message": "task deleted successfully"
}
```

## Postman verification

Create a collection named `Task Management API` with:

- `POST /register`
- `POST /login`
- `GET /tasks`
- `GET /tasks/:id`
- `POST /tasks`
- `PUT /tasks/:id`
- `DELETE /tasks/:id`

Use:

- `{{baseUrl}} = http://localhost:8080`

For protected requests, add the `Authorization` header and paste the JWT token from the login response.

Test successful requests and error cases such as invalid IDs, missing tasks, missing tokens, expired tokens, and non-admin delete attempts.

## MongoDB verification

After creating or updating a task, verify the document directly in MongoDB Compass or with a query such as:

```js
db.tasks.find().sort({ id: 1 });
```

Because tasks are stored in MongoDB, data remains available after API restarts.

## Security notes

- Passwords are hashed with bcrypt before storage.
- JWT tokens are signed with `HS256` using `JWT_SECRET`.
- Protected routes reject missing, malformed, expired, or invalid tokens.
