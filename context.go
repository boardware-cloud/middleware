package middleware

import (
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB
var accountRepository core.AccountRepository

func init() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	db, err = model.NewConnection(user, password, host, port, database)
	accountRepository = core.NewAccountRepository(db)
	if err != nil {
		panic(err)
	}
}

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
