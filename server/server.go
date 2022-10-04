package server

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	router.Run("localhost:8080")
}
