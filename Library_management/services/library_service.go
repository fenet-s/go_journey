package services

import (
	"Library_management/models"
	"fmt"
)

type LibraryManager interface {
	AddBook(book models.Book)
	AddMember(member models.Member)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAllBooks() []models.Book
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ListMembers() map[int]models.Member
}

type Library struct {
	Books   map[int]models.Book
	members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	if book.Status == "" {
		book.Status = models.BookStatusAvailable
	}
	l.Books[book.ID] = book
}

func (l *Library) AddMember(member models.Member) {
	l.members[member.ID] = member
}

func (l *Library) ListMembers() map[int]models.Member {
	members := make(map[int]models.Member, len(l.members))
	for id, member := range l.members {
		members[id] = member
	}
	return members
}

func (l *Library) RemoveBook(bookID int) error {
	_, exists := l.Books[bookID]

	if !exists {
		return fmt.Errorf("book with ID %d not found", bookID)

	}
	delete(l.Books, bookID)
	return nil
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	member, exists := l.members[memberID]
	if !exists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	if !book.IsAvailable() {
		return fmt.Errorf("book %q is already borrowed", book.Title)
	}

	if member.HasBorrowedBook(bookID) {
		return fmt.Errorf("member %q already borrowed this book", member.Name)
	}

	book.MarkBorrowed()
	if err := member.BorrowBook(book); err != nil {
		return err
	}

	l.Books[bookID] = book
	l.members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]
	if !exists {
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	member, exists := l.members[memberID]
	if !exists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	if book.IsAvailable() {
		return fmt.Errorf("book %q is already available", book.Title)
	}

	if !member.HasBorrowedBook(bookID) {
		return fmt.Errorf("member %q did not borrow this book", member.Name)
	}

	if _, found := member.ReturnBook(bookID); !found {
		return fmt.Errorf("member %q did not borrow this book", member.Name)
	}

	book.MarkAvailable()

	l.Books[bookID] = book
	l.members[memberID] = member

	return nil
}

func (l *Library) ListAllBooks() []models.Book {
	books := make([]models.Book, 0, len(l.Books))
	for _, book := range l.Books {
		books = append(books, book)
	}
	return books
}

// available books

func (l *Library) ListAvailableBooks() []models.Book {

	var availableBooks []models.Book

	for _, book := range l.Books {

		if book.IsAvailable() {
			availableBooks = append(availableBooks, book)
		}

	}

	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.members[memberID]
	if !exists {
		return nil
	}

	books := make([]models.Book, len(member.BorrowedBooks))
	copy(books, member.BorrowedBooks)
	return books
}
