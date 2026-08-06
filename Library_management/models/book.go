package models

type Book struct {
	ID     int
	Title  string
	Author string
	Status string
}

const (
	BookStatusAvailable = "Available"
	BookStatusBorrowed  = "Borrowed"
)

func (b Book) IsAvailable() bool {
	return b.Status == BookStatusAvailable
}

func (b *Book) MarkBorrowed() {
	b.Status = BookStatusBorrowed
}

func (b *Book) MarkAvailable() {
	b.Status = BookStatusAvailable
}
