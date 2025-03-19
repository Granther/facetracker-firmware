package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", IndexHandler())

	r.POST("/submit", SubmitHandler())

	return r
}
