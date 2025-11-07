package handlers

import (
	"encoding/json"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/models"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/storage"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/utils"
	"net/http"
	"strings"
)

func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	utils.RespondJSON(w, http.StatusOK, storage.Store.Books)
}

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/books/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, book := range storage.Store.Books {
		if book.ID == id {
			utils.RespondJSON(w, http.StatusOK, book)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Book not found")
}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	book.ID = len(storage.Store.Books) + 1
	storage.Store.Books = append(storage.Store.Books, book)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, book)
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/books/")

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, book := range storage.Store.Books {
		if book.ID == id {
			updatedBook.ID = id
			storage.Store.Books[i] = updatedBook
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedBook)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Book not found")
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/books/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, book := range storage.Store.Books {
		if book.ID == id {
			storage.Store.Books = append(storage.Store.Books[:i], storage.Store.Books[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Book deleted"})
			return
		}
	}
	
	utils.RespondError(w, http.StatusNotFound, "Book not found")
}
func BooksRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/books/") && len(r.URL.Path) > len("/books/") {
			GetBookHandler(w, r)
		} else {
			GetBooksHandler(w, r)
		}
	case http.MethodPost:
		CreateBookHandler(w, r)
	case http.MethodPut:
		UpdateBookHandler(w, r)
	case http.MethodDelete:
		DeleteBookHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
