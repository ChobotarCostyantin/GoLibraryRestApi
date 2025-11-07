package handlers

import (
	"encoding/json"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/models"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/storage"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/utils"
	"net/http"
	"strings"
)

func GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	utils.RespondJSON(w, http.StatusOK, storage.Store.Authors)
}

func GetAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/authors/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, author := range storage.Store.Authors {
		if author.ID == id {
			utils.RespondJSON(w, http.StatusOK, author)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Author not found")
}

func CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	author.ID = len(storage.Store.Authors) + 1
	storage.Store.Authors = append(storage.Store.Authors, author)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, author)
}

func UpdateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/authors/")

	var updatedAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&updatedAuthor); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, author := range storage.Store.Authors {
		if author.ID == id {
			updatedAuthor.ID = id
			storage.Store.Authors[i] = updatedAuthor
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedAuthor)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Author not found")
}

func DeleteAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/authors/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, author := range storage.Store.Authors {
		if author.ID == id {
			storage.Store.Authors = append(storage.Store.Authors[:i], storage.Store.Authors[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Author deleted"})
			return
		}
	}
	
	utils.RespondError(w, http.StatusNotFound, "Author not found")
}
func AuthorsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/authors/") && len(r.URL.Path) > len("/authors/") {
			GetAuthorHandler(w, r)
		} else {
			GetAuthorsHandler(w, r)
		}
	case http.MethodPost:
		CreateAuthorHandler(w, r)
	case http.MethodPut:
		UpdateAuthorHandler(w, r)
	case http.MethodDelete:
		DeleteAuthorHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
