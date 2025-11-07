package models

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type Book struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Title    string `json:"title"`
	Pages    int    `json:"pages"`
}

type Librarian struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int `json:"age"`
}