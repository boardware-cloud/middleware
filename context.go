package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Authorize(c)
		if auth.Status == Authorized {
			// account, _ := core.FindAccount(auth.AccountId)
			// c.Set("account", core.Account{Email: "good"})
		}
		c.Next()
	}
}
