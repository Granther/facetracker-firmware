package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"captive-portal/internal/network"
)

func IndexHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
	return gin.HandlerFunc(fn)
}

func SubmitHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ssid := c.PostForm("ssid")
		pass := c.PostForm("password")

		err := network.ConnectNetwork(ssid, pass) // If we connect, set state to setup
		// We dont loose the captive portal because it is a 'separate' adapter
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.HTML(http.StatusOK, "index.html", gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "done.html", gin.H{})
	}
	return gin.HandlerFunc(fn)
}
