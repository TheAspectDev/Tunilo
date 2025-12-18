package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NOT A PART OF THE SOURCE CODE
// IT IS FOR TESTING LOCALLY

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	// GET /
	r.GET("/", func(c *gin.Context) {
		html := `
		<button onclick="fetch('/echo', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({})
		})">
			heck
		</button>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// POST /echo
	r.POST("/echo", func(c *gin.Context) {
		var jsonData map[string]interface{}

		// Try to bind JSON
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			// Not JSON → read raw body
			body, _ := c.GetRawData()
			c.JSON(http.StatusOK, gin.H{
				"received": string(body),
			})
			return
		}

		// Valid JSON
		c.JSON(http.StatusOK, gin.H{
			"json": jsonData,
		})
	})

	// Run server
	r.Run("localhost:8999")
}
