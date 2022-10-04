package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"github.com/xtravanilla/go-gin-test/db"
	"github.com/xtravanilla/go-gin-test/models"
)

type LibrarianController struct{}

func (h LibrarianController) GetAllBooks(c *gin.Context) {
	dyna := db.GetDB()

	out, err := dyna.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("books-table"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)

	c.IndentedJSON(http.StatusNotFound, gin.H{"books": out.Items})
}

func (h LibrarianController) RemoveABookByID(c *gin.Context) {
	bookId := c.Param("id")

	dyna := db.GetDB()

	params := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bookId},
		},
		TableName: aws.String("book-table"),
	}

	out, err := dyna.DeleteItem(context.TODO(), params)

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)

	c.IndentedJSON(http.StatusNotFound, gin.H{"book": out.Attributes})
}

func (h LibrarianController) AddABook(c *gin.Context) {
	var newBook models.Book

	marshalledBook, err := dynamodbattribute.MarshalMap(newBook)
	if err != nil {
		log.Fatalf("Got error marshalling book map: %s", err)
	}

	dyna := db.GetDB()

	params := &dynamodb.PutItemInput{
		TableName: aws.String("book-table"),
		Item:      marshalledBook,
	}

	out, err := dyna.PutItem(context.TODO(), params)

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)

	c.IndentedJSON(http.StatusNotFound, gin.H{"book": out.Attributes})

}
