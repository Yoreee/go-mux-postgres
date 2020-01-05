package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "Riaz"
	dbname = "testdb"
)

// Book Struct (Model)
type Book struct {
	ID    int    `json:"id"`
	ISBN  string `json:"isbn"`
	Title string `json:"title"`
	// Author *Author `json:"author"`
}

// Init books var as a slice Book struct
var books []Book

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	fmt.Println("Successfully connected!")
	rows, err := db.Query(`SELECT id, isbn, title FROM book`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.ID, &book.ISBN, &book.Title)
		if err != nil {
			panic(err)
		}
		fmt.Println(book)
		books = append(books, book)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&books)
	books = nil
}

// Get one book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Successfully connected!")

	id := params["id"]
	book := Book{}
	sqlStatement := `SELECT id, isbn, title FROM book WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	errr := row.Scan(&book.ID, &book.ISBN, &book.Title)
	if errr != nil {
		if errr == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
	}
	fmt.Println(book)
	json.NewEncoder(w).Encode(&book)

}

// Add new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Successfully connected!")

	isbn := &book.ISBN
	title := &book.Title
	id := 0
	sqlStatement := `insert into book (isbn, title) values ($1, $2) returning id`
	err = db.QueryRow(sqlStatement, isbn, title).Scan(&id)

	if err != nil {
		panic(err)
	}
	book.ID = id
	json.NewEncoder(w).Encode(&book)

	fmt.Println(&book)
	fmt.Println("New ID is:", id)
}

// Update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// params := mux.Vars(r)
	// for index, item := range books {
	// 	if item.ID == params["id"] {
	// 		books = append(books[:index], books[index+1:]...)
	// 		var book Book
	// 		_ = json.NewDecoder(r.Body).Decode(&book)
	// 		book.ID = params["id"]
	// 		books = append(books, book)
	// 		json.NewEncoder(w).Encode(book)
	// 		return
	// 	}
	// }

}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]
	sqlStatement := `delete from book where id = $1`
	db.Exec(sqlStatement, id)
}

func main() {
	// Init mux
	r := mux.NewRouter()

	// Route handlers / endpoints

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}
