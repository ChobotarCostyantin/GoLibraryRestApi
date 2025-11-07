package handlers

import (
	"github.com/ChobotarCostyantin/GoLibraryRestApi/utils"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.RespondError(w, http.StatusNotFound, "Not found")
		return
	}
	info := map[string]interface{}{
		"message": "Library REST API",
		"endpoints": map[string]string{
			"GET /authors":            "Get all authors",
			"GET /authors/{id}":       "Get author by ID",
			"POST /authors":           "Create a new author",
			"PUT /authors/{id}":       "Update author",
			"DELETE /authors/{id}":    "Delete author",
			"GET /books":              "Get all books",
			"GET /books/{id}":         "Get book by ID",
			"POST /books":             "Create a new book",
			"PUT /books/{id}":         "Update book",
			"DELETE /books/{id}":      "Delete book",
			"GET /librarians":         "Get all librarians",
			"GET /librarians/{id}":    "Get librarian by ID",
			"POST /librarians":        "Create a new librarian",
			"PUT /librarians/{id}":    "Update librarian",
			"DELETE /librarians/{id}": "Delete librarian",
		},
	}
	utils.RespondJSON(w, http.StatusOK, info)
}
