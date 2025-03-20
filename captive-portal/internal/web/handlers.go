package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"captive-portal/internal/network"
)

func IndexHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		network.DiscoverNets()
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
	return gin.HandlerFunc(fn)
}

