package handlers

import (
	"encoding/json"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/models"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/storage"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/utils"
	"net/http"
	"strings"
)

func GetLibrariansHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	utils.RespondJSON(w, http.StatusOK, storage.Store.Librarians)
}

func GetLibrarianHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/librarians/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, librarian := range storage.Store.Librarians {
		if librarian.ID == id {
			utils.RespondJSON(w, http.StatusOK, librarian)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Librarian not found")
}

func CreateLibrarianHandler(w http.ResponseWriter, r *http.Request) {
	var librarian models.Librarian
	if err := json.NewDecoder(r.Body).Decode(&librarian); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	librarian.ID = len(storage.Store.Librarians) + 1
	storage.Store.Librarians = append(storage.Store.Librarians, librarian)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, librarian)
}

func UpdateLibrarianHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/librarians/")

	var updatedLibrarian models.Librarian
	if err := json.NewDecoder(r.Body).Decode(&updatedLibrarian); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, librarian := range storage.Store.Librarians {
		if librarian.ID == id {
			updatedLibrarian.ID = id
			storage.Store.Librarians[i] = updatedLibrarian
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedLibrarian)
			return
		}
	}

	utils.RespondError(w, http.StatusNotFound, "Librarian not found")
}

func DeleteLibrarianHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/librarians/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, librarian := range storage.Store.Librarians {
		if librarian.ID == id {
			storage.Store.Librarians = append(storage.Store.Librarians[:i], storage.Store.Librarians[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Librarian deleted"})
			return
		}
	}
	
	utils.RespondError(w, http.StatusNotFound, "Librarian not found")
}
func LibrariansRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/librarians/") && len(r.URL.Path) > len("/librarians/") {
			GetLibrarianHandler(w, r)
		} else {
			GetLibrariansHandler(w, r)
		}
	case http.MethodPost:
		CreateLibrarianHandler(w, r)
	case http.MethodPut:
		UpdateLibrarianHandler(w, r)
	case http.MethodDelete:
		DeleteLibrarianHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
