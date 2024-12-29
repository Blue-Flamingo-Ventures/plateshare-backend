package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func main() {
	r := gin.Default()

	r.POST("/test", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	r.Run(":8080") // Run on port 8080
}
