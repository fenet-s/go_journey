package controllers

import (
	"Library_management/models"
	"Library_management/services"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RemoveBookController(library *services.Library, bookID int) {
	err := library.RemoveBook(bookID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book removed successfully")
}

func RunLibraryConsole(library *services.Library) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("==== Library Management System ====")
		fmt.Println("1. Add book")
		fmt.Println("2. Add member")
		fmt.Println("3. Remove book")
		fmt.Println("4. Borrow book")
		fmt.Println("5. Return book")
		fmt.Println("6. List all books")
		fmt.Println("7. List available books")
		fmt.Println("8. List borrowed books for a member")
		fmt.Println("9. List members")
		fmt.Println("0. Exit")

		choice, err := readInt(reader, "Choose an option: ")
		if err != nil {
			fmt.Println("Invalid input:", err)
			continue
		}

		switch choice {
		case 1:
			addBook(library, reader)
		case 2:
			addMember(library, reader)
		case 3:
			bookID, err := readInt(reader, "Enter book ID to remove: ")
			if err != nil {
				fmt.Println("Invalid input:", err)
				continue
			}
			RemoveBookController(library, bookID)
		case 4:
			borrowBook(library, reader)
		case 5:
			returnBook(library, reader)
		case 6:
			printBooks("All books", library.ListAllBooks())
		case 7:
			printBooks("Available books", library.ListAvailableBooks())
		case 8:
			memberID, err := readInt(reader, "Enter member ID: ")
			if err != nil {
				fmt.Println("Invalid input:", err)
				continue
			}
			members := library.ListMembers()
			member, exists := members[memberID]
			if !exists {
				fmt.Printf("Member with ID %d not found\n", memberID)
				continue
			}
			printBooks(fmt.Sprintf("Borrowed books for %s", member.Name), library.ListBorrowedBooks(memberID))
		case 9:
			printMembers(library.ListMembers())
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Please choose a valid option.")
		}
	}
}

func addBook(library *services.Library, reader *bufio.Reader) {
	id, err := readInt(reader, "Enter book ID: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	title, err := readString(reader, "Enter book title: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	author, err := readString(reader, "Enter book author: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	library.AddBook(models.Book{ID: id, Title: title, Author: author})

	fmt.Println("Book added successfully")
}

func addMember(library *services.Library, reader *bufio.Reader) {
	id, err := readInt(reader, "Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	name, err := readString(reader, "Enter member name: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	library.AddMember(models.Member{ID: id, Name: name})
	fmt.Println("Member added successfully")
}

func borrowBook(library *services.Library, reader *bufio.Reader) {
	bookID, err := readInt(reader, "Enter book ID to borrow: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	memberID, err := readInt(reader, "Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	if err := library.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book borrowed successfully")
}

func returnBook(library *services.Library, reader *bufio.Reader) {
	bookID, err := readInt(reader, "Enter book ID to return: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	memberID, err := readInt(reader, "Enter member ID: ")
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	if err := library.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book returned successfully")
}

func readString(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func readInt(reader *bufio.Reader, prompt string) (int, error) {
	input, err := readString(reader, prompt)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("%q is not a valid number", input)
	}

	return value, nil
}

func printBooks(title string, books []models.Book) {
	fmt.Println()
	fmt.Println(title)
	if len(books) == 0 {
		fmt.Println("  No books found.")
		return
	}

	for _, book := range books {
		fmt.Printf("  ID:%d | Title:%s | Author:%s | Status:%s\n", book.ID, book.Title, book.Author, book.Status)
	}
}

func printMembers(members map[int]models.Member) {
	fmt.Println()
	fmt.Println("Members")
	if len(members) == 0 {
		fmt.Println("  No members found.")
		return
	}

	for _, member := range members {
		fmt.Printf("  ID:%d | Name:%s | Borrowed books:%d\n", member.ID, member.Name, len(member.BorrowedBooks))
	}
}
