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

type UserController struct{}

func (h UserController) CheckOutBook(c *gin.Context) {
	bookId := c.Param("book_id")
	userId := c.Param("user_id")

	dyna := db.GetDB()

	book_out, err := dyna.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("book-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bookId},
		},
		UpdateExpression: aws.String("set available = :available"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":available": &types.AttributeValueMemberS{Value: "false"},
		},
	})

	user_out, err := dyna.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("user-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
		UpdateExpression: aws.String("set checkedoutBooks = :checkedoutBooks"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":checkedoutBooks": &types.AttributeValueMemberS{Value: bookId},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(user_out.Attributes)
	fmt.Println(book_out.Attributes)

	c.IndentedJSON(http.StatusNotFound, gin.H{"book": book_out.Attributes, "user": user_out.Attributes})
}

func (h UserController) ReturnBook(c *gin.Context) {
	bookId := c.Param("book_id")
	userId := c.Param("user_id")

	dyna := db.GetDB()

	book_out, err := dyna.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("book-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bookId},
		},
		UpdateExpression: aws.String("set available = :available"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":available": &types.AttributeValueMemberS{Value: "true"},
		},
	})

	user_out, err := dyna.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("user-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
		UpdateExpression: aws.String("DELETE checkedoutBooks = :checkedoutBooks"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":checkedoutBooks": &types.AttributeValueMemberS{Value: bookId},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(user_out.Attributes)
	fmt.Println(book_out.Attributes)

	c.IndentedJSON(http.StatusNotFound, gin.H{"book": book_out.Attributes, "user": user_out.Attributes})
}

func (h UserController) ListAllBooksByUser(c *gin.Context) {
	userId := c.Param("user_id")

	dyna := db.GetDB()

	out, err := dyna.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("user-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userId},
		},
	})

	user := models.User{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &user)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	marshalledBookIds, err := dynamodbattribute.MarshalMap(user.CheckedoutBooks)
	if err != nil {
		log.Fatalf("Got error marshalling book map: %s", err)
	}

	batch, err := dyna.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			"user-table": {
				Keys: []map[string]types.AttributeValue{marshalledBookIds},
			},
		},
	})

	c.IndentedJSON(http.StatusNotFound, gin.H{"books": batch.Responses})
}
