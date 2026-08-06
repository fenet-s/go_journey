package main

import (
	"Library_management/controllers"
	"Library_management/models"
	"Library_management/services"
)

func main() {
	library := services.NewLibrary()

	library.AddBook(models.Book{ID: 1, Title: "Go Programming Language", Author: "Alan A. A. Donovan"})
	library.AddBook(models.Book{ID: 2, Title: "Clean Code", Author: "Robert C. Martin"})
	library.AddMember(models.Member{ID: 1, Name: "John Doe"})

	controllers.RunLibraryConsole(library)

}
