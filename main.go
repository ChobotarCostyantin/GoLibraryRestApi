package main

import (
	"github.com/ChobotarCostyantin/GoLibraryRestApi/handlers"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/middleware"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/storage"
	"log"
	"net/http"
)

func main() {
	if err := storage.Store.LoadFromFile(); err != nil {
		log.Println("Creating new storage file...")
	}

	// Routers
	http.HandleFunc("/authors", middleware.Chain(
		handlers.AuthorsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/authors/", middleware.Chain(
		handlers.AuthorsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/books", middleware.Chain(
		handlers.BooksRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/books/", middleware.Chain(
		handlers.BooksRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/librarians", middleware.Chain(
		handlers.LibrariansRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/librarians/", middleware.Chain(
		handlers.LibrariansRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/", middleware.LoggingMiddleware(handlers.RootHandler))

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("API Key for authorization: %s", middleware.ValidAPIKey)

	log.Fatal(http.ListenAndServe(port, nil))
}
