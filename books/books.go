package books

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	qmgo "github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookCreateUpdateRequest struct {
	Title  string `form:"title" binding:"required"`
	Author string `form:"author"`
}

type BookResponse struct {
	Id        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	Author    string             `json:"author"`
	CreatedAt time.Time          `json:"createdAt" binding:"required"`
	UpdatedAt time.Time          `json:"updatedAt" binding:"required"`
}

type BookListResponse struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
}

type Book struct {
	field.DefaultField `bson: "inline"`
	Title              string `bson:"title" validate:"required"`
	Author             string `bson:"author"`
}

var db *qmgo.QmgoClient
var ctx context.Context

func InitDB(databaseUri string, databaseName string, collection string) error {
	var err error

	ctx = context.Background()
	db, err = qmgo.Open(ctx, &qmgo.Config{Uri: databaseUri, Database: databaseName, Coll: collection})
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	if err := db.Close(ctx); err != nil {
		log.Fatal(err)
	}
}

func CreateBook(ctx *gin.Context) {
	var newBook BookCreateUpdateRequest

	if err := ctx.BindJSON(&newBook); err != nil {
		// fmt.Printf("ERR IN CREATE: %+v\n", err)
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	book := Book{
		Title:  newBook.Title,
		Author: newBook.Author,
	}

	_, err := db.InsertOne(ctx, &book)
	if err != nil {
		fmt.Printf("ERR: %+v\n", err)
		ctx.JSON(http.StatusInternalServerError, "Something went wrong")
		return
	}

	ctx.JSON(http.StatusCreated, GetBooksResponse(book))
}

func GetBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	var book Book
	err = db.Find(ctx, bson.M{"_id": bookId}).One(&book)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book not found")
		return
	}

	ctx.JSON(http.StatusOK, GetBooksResponse(book))
}

func GetBooks(ctx *gin.Context) {
	var books []BookListResponse

	err := db.Find(ctx, bson.M{}).All(&books)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}

	// to send success response on completion
	ctx.JSON(http.StatusOK, books)
}

func UpdateBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	var newBook BookCreateUpdateRequest
	if err := ctx.BindJSON(&newBook); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Update Request")
		return
	}

	var book Book

	err = db.Find(ctx, bson.M{"_id": bookId}).One(&book)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book not found")
		return
	}

	// Set updated values in book
	book.Title = newBook.Title
	book.Author = newBook.Author
	err = db.ReplaceOne(ctx, bson.M{"_id": bookId}, &book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Something went wrong with updating book")
		return
	}

	ctx.JSON(http.StatusOK, GetBooksResponse(book))
}

func DeleteBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	var book Book

	err = db.Find(ctx, bson.M{"_id": bookId}).One(&book)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "Book not found")
		return
	}

	// deleting the book
	err = db.RemoveId(ctx, bookId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Something went wrong with deleting book")
		return
	}

	ctx.JSON(http.StatusOK, true)
}

func GetBooksResponse(book Book) (bookResponse BookResponse) {
	return BookResponse{
		Id:        book.DefaultField.Id,
		Title:     book.Title,
		Author:    book.Author,
		CreatedAt: book.CreateAt,
		UpdatedAt: book.UpdateAt,
	}
}
