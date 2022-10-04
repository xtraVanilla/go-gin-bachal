package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xtravanilla/go-gin-test/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(middlewares.AuthMiddleware())

	librarian := new(controllers.LibrarianController)
	user := new(controllers.UserController)

	router.GET("/getBooks", librarian.GetAllBooks)
	router.GET("/addBook", librarian.AddABook)
	router.DELETE("/removeBook/:id", librarian.RemoveABookByID)

	router.POST("/checkoutbook/:id", user.CheckOutBook)
	router.POST("/return/:id", user.ReturnBook)
	router.GET("/checkedoutbooks/:id", user.ListAllBooksByUser)

	return router
}
