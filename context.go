package middleware

import (
	"fmt"

	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Authorize(c)
		fmt.Println(auth)
		if auth.Status == Authorized {
			account, _ := core.FindAccount(auth.AccountId)
			c.Set("account", account)
		}
		c.Next()
	}
}
