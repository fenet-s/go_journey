package models

type Member struct {
	ID            int
	Name          string
	BorrowedBooks []Book
}

func (m *Member) BorrowBook(book Book) error {
	for _, borrowed := range m.BorrowedBooks {
		if borrowed.ID == book.ID {
			return nil
		}
	}

	m.BorrowedBooks = append(m.BorrowedBooks, book)
	return nil
}

func (m *Member) ReturnBook(bookID int) (Book, bool) {
	for i, borrowed := range m.BorrowedBooks {
		if borrowed.ID == bookID {
			removed := borrowed
			m.BorrowedBooks = append(m.BorrowedBooks[:i], m.BorrowedBooks[i+1:]...)
			return removed, true
		}
	}

	return Book{}, false
}

func (m Member) HasBorrowedBook(bookID int) bool {
	for _, borrowed := range m.BorrowedBooks {
		if borrowed.ID == bookID {
			return true
		}
	}

	return false
}
