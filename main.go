package main

import (
	"log"
	"net/http"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/handlers"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/storage"
)

func main() {
	if err := storage.Store.LoadFromFile(); err != nil {
		log.Println("Creating new storage file...")
	}

	// Routers
	http.HandleFunc("/authors", handlers.AuthorsRouter)
	http.HandleFunc("/authors/", handlers.AuthorsRouter)
	http.HandleFunc("/books", handlers.BooksRouter)
	http.HandleFunc("/books/", handlers.BooksRouter)
	http.HandleFunc("/librarians", handlers.LibrariansRouter)
	http.HandleFunc("/librarians/", handlers.LibrariansRouter)

	// Root page with api description
	http.HandleFunc("/", handlers.RootHandler)

	port := ":8080"

	log.Fatal(http.ListenAndServe(port, nil))
}
