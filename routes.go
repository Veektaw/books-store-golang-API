package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./books.db")
	if err != nil {
		panic(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT,
		quantity INTEGER
	);`
	if _, err = db.Exec(createTableSQL); err != nil {
		panic(err)
	}
}

func getBooks(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, author, quantity FROM books")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving books"})
		return
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error scanning book"})
			return
		}
		books = append(books, b)
	}

	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	var b book
	err := db.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = ?", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving book"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, b)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	_, err := db.Exec("UPDATE books SET title = ?, author = ?, quantity = ? WHERE id = ?",
		updatedBook.Title, updatedBook.Author, updatedBook.Quantity, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error updating book"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedBook)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting book"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	newBook.ID = uuid.New().String()

	_, err := db.Exec("INSERT INTO books (id, title, author, quantity) VALUES (?, ?, ?, ?)",
		newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error adding book"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newBook)
}
