package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Authorize(c)
		fmt.Println(auth)
		c.Next()
	}
}
