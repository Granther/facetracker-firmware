package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{"captcha": numCaptcha, "csrf": token})
	}
	return gin.HandlerFunc(fn)
}

