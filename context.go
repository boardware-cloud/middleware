package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// auth := Authorize(c)
		// fmt.Println(auth)
		// if auth.Status == Authorized {
		// 	var account core.Account
		// 	core.FindAccount(auth.AccountId)
		// 	c.Set("account", account)
		// 	return
		// }
		c.Next()
	}
}
