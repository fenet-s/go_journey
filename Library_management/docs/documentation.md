## Library Management System

This project is a simple console-based library management system written in Go.

### Features

- Add books and members
- Remove books
- Borrow and return books
- List all books
- List available books
- List borrowed books for a specific member
- List registered members

### Project structure

- `main.go` starts the application and seeds a few sample records.
- `controllers/library_controller.go` handles terminal input/output.
- `models/book.go` defines the `Book` struct and book status helpers.
- `models/member.go` defines the `Member` struct and borrowed-book helpers.
- `services/library_service.go` contains the library business logic.

### How it works

The service layer stores books in a map keyed by book ID and members in a map keyed by member ID. Borrowing a book updates both the book status and the member's borrowed-book slice. Returning a book reverses that change.

### Running the app

Run the program from the module root:

- `go run .`

### Notes

- Book status values are `Available` and `Borrowed`.
- The controller validates input and prints user-friendly errors for missing books, missing members, and invalid operations.
