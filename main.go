package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	initDB()

	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)

	log.Fatal(router.Run("localhost:8080"))
}
