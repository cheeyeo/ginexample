package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cheeyeo/ginexample/books"

	"github.com/gin-gonic/gin"
)

func main() {
	// If running locally without docker, get the mongodb container IP
	// docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <mongo container name>
	// export MONGO_URL=mongodb://172.0.0.2:27017
	// If running via compose the host is the service name set in docker compose file i.e. mongodb

	uri := os.Getenv("MONGO_URL")
	log.Printf("URI IS %s", uri)
	database := "test"
	collection := "books"

	err := books.InitDB(uri, database, collection)
	if err != nil {
		log.Fatal(err)
	}

	defer books.CloseDB()

	router := gin.Default()
	router.POST("/books", books.CreateBook)
	router.GET("/books", books.GetBooks)
	router.GET("/books/:bookId", books.GetBook)
	router.PATCH("/books/:bookId", books.UpdateBook)
	router.DELETE("/books/:bookId", books.DeleteBook)
	fmt.Println("Service up and running")
	router.Run(":8080")
}
