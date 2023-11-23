package middleware

import (
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
)

var db = model.GetDB()
var accountRepository = core.GetAccountRepository()

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Authorize(c)
		if auth.Status == Authorized {
			account := accountRepository.GetById(auth.AccountId)
			if account != nil {
				c.Set("account", account)
			}
		}
		c.Next()
	}
}
