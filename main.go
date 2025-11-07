package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	_"path"
	"strconv"
	"strings"
)

// Моделі даних
type Book struct {
	BookID int    `json:"book_id"`
	Title  string `json:"title"`
	Pages  int    `json:"pages"`
}

type Library struct {
	LibraryID  int    `json:"library_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Books      []Book `json:"books"`
}

// Сховище даних
type Storage struct {
	Libraries []Library `json:"libraries"`
	NextID    int       `json:"next_id"`
}

const storageFile = "libraries.json"

var storage Storage

// Ініціалізація сховища
func initStorage() {
	data, err := os.ReadFile(storageFile)
	if err != nil {
		// Якщо файл не існує, створюємо нове сховище
		storage = Storage{
			Libraries: []Library{},
			NextID:    1,
		}
		saveStorage()
		return
	}

	err = json.Unmarshal(data, &storage)
	if err != nil {
		log.Printf("Помилка десеріалізації: %v", err)
		storage = Storage{
			Libraries: []Library{},
			NextID:    1,
		}
	}
}

// Збереження даних у файл
func saveStorage() error {
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storageFile, data, 0644)
}

// Відповідь з помилкою
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Відповідь з JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Помилка серіалізації"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Обробники для бібліотек

// GET /libraries - отримати всі бібліотеки
// POST /libraries - створити нову бібліотеку
func librariesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		respondWithJSON(w, http.StatusOK, storage.Libraries)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неможливо прочитати тіло запиту")
			return
		}
		defer r.Body.Close()

		var library Library
		err = json.Unmarshal(body, &library)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Невірний формат JSON")
			return
		}

		library.LibraryID = storage.NextID
		storage.NextID++
		storage.Libraries = append(storage.Libraries, library)

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusCreated, library)

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не підтримується")
	}
}

// GET /libraries/{library_id} - отримати бібліотеку за ID
// PUT /libraries/{library_id} - оновити бібліотеку
// DELETE /libraries/{library_id} - видалити бібліотеку
func libraryHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID з URL
	path := strings.TrimPrefix(r.URL.Path, "/libraries/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Невірний ID")
		return
	}

	// Знаходимо індекс бібліотеки
	index := -1
	for i, lib := range storage.Libraries {
		if lib.LibraryID == id {
			index = i
			break
		}
	}

	switch r.Method {
	case http.MethodGet:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Бібліотека не знайдена")
			return
		}
		respondWithJSON(w, http.StatusOK, storage.Libraries[index])

	case http.MethodPut:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Бібліотека не знайдена")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неможливо прочитати тіло запиту")
			return
		}
		defer r.Body.Close()

		var library Library
		err = json.Unmarshal(body, &library)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Невірний формат JSON")
			return
		}

		library.LibraryID = id
		storage.Libraries[index] = library

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusOK, library)

	case http.MethodDelete:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Бібліотека не знайдена")
			return
		}

		storage.Libraries = append(storage.Libraries[:index], storage.Libraries[index+1:]...)

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Бібліотека видалена"})

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не підтримується")
	}
}

// GET /libraries/{library_id}/books - отримати книги бібліотеки
// POST /libraries/{library_id}/books - додати книгу до бібліотеки
func libraryBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID з URL
	path := strings.TrimPrefix(r.URL.Path, "/libraries/")
	path = strings.TrimSuffix(path, "/books")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Невірний ID")
		return
	}

	// Знаходимо індекс бібліотеки
	index := -1
	for i, lib := range storage.Libraries {
		if lib.LibraryID == id {
			index = i
			break
		}
	}

	if index == -1 {
		respondWithError(w, http.StatusNotFound, "Бібліотека не знайдена")
		return
	}

	switch r.Method {
	case http.MethodGet:
		respondWithJSON(w, http.StatusOK, storage.Libraries[index].Books)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неможливо прочитати тіло запиту")
			return
		}
		defer r.Body.Close()

		var book Book
		err = json.Unmarshal(body, &book)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Невірний формат JSON")
			return
		}

		storage.Libraries[index].Books = append(storage.Libraries[index].Books, book)

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusCreated, book)

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не підтримується")
	}
}

// GET /libraries/{library_id}/books/{book_id} - отримати книгу за ID
// PUT /libraries/{library_id}/books/{book_id} - оновити книгу
// DELETE /libraries/{library_id}/books/{book_id} - видалити книгу
func libraryBookHandler(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID з URL
	path := strings.TrimPrefix(r.URL.Path, "/libraries/")
	path = strings.TrimSuffix(path, "/books/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Невірний ID")
		return
	}

	// Знаходимо індекс бібліотеки
	index := -1
	for i, lib := range storage.Libraries {
		if lib.LibraryID == id {
			index = i
			break
		}
	}

	if index == -1 {
		respondWithError(w, http.StatusNotFound, "Бібліотека не знайдена")
		return
	}

	// Отримуємо ID книги
	path = strings.TrimPrefix(r.URL.Path, "/libraries/"+strconv.Itoa(id)+"/books/")
	id, err = strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Невірний ID")
		return
	}

	// Знаходимо індекс книги
	index = -1
	for i, book := range storage.Libraries[index].Books {
		if book.BookID == id {
			index = i
			break
		}
	}

	switch r.Method {
	case http.MethodGet:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Книга не знайдена")
			return
		}
		respondWithJSON(w, http.StatusOK, storage.Libraries[index].Books[index])

	case http.MethodPut:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Книга не знайдена")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неможливо прочитати тіло запиту")
			return
		}
		defer r.Body.Close()

		var book Book
		err = json.Unmarshal(body, &book)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Невірний формат JSON")
			return
		}

		book.BookID = id
		storage.Libraries[index].Books[index] = book

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusOK, book)

	case http.MethodDelete:
		if index == -1 {
			respondWithError(w, http.StatusNotFound, "Книга не знайдена")
			return
		}

		storage.Libraries[index].Books = append(storage.Libraries[index].Books[:index], storage.Libraries[index].Books[index+1:]...)

		err = saveStorage()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Помилка збереження")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Книга видалена"})

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не підтримується")
	}
}

func main() {
	// Ініціалізуємо сховище
	initStorage()

	// Реєструємо обробники
	http.HandleFunc("/libraries", librariesHandler)
	http.HandleFunc("/libraries/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/books") {
			libraryBooksHandler(w, r)
		} else {
			libraryHandler(w, r)
		}
	})

	// Запускаємо сервер
	port := ":8080"
	fmt.Printf("🚀 Сервер запущено на http://localhost%s\n", port)
	fmt.Println("\nДоступні ендпоінти:")
	fmt.Println("  GET    /libraries           - отримати всі бібліотеки")
	fmt.Println("  POST   /libraries           - створити бібліотеку")
	fmt.Println("  GET    /libraries/{library_id}      - отримати бібліотеку")
	fmt.Println("  PUT    /libraries/{library_id}      - оновити бібліотеку")
	fmt.Println("  DELETE /libraries/{library_id}      - видалити бібліотеку")
	fmt.Println("  GET    /libraries/{library_id}/books - отримати книги")
	fmt.Println("  POST   /libraries/{library_id}/books - додати книгу")
	fmt.Println("  GET    /libraries/{library_id}/books/{book_id} - отримати книгу")
	fmt.Println("  PUT    /libraries/{library_id}/books/{book_id} - оновити книгу")
	fmt.Println("  DELETE /libraries/{library_id}/books/{book_id} - видалити книгу")

	log.Fatal(http.ListenAndServe(port, nil))
}