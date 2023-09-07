package middleware

import (
	"net/http"
	"strings"

	"github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	DB, err = model.NewConnection(user, password, host, port, database)
	if err != nil {
		panic(err)
	}
}

type health struct {
	Status string   `json:"status"`
	Checks []string `json:"checks"`
}

type Headers struct {
	Authorization string
}

type Authentication struct {
	Status    AuthenticationStatus
	Role      constants.Role
	AccountId uint
	Email     string
}

type AuthenticationStatus string

const (
	Authorized   AuthenticationStatus = "Authorized"
	Unauthorized AuthenticationStatus = "Unauthorized"
)

func IsRoot(c *gin.Context, next func(c *gin.Context, account core.Account)) {
	GetAccount(c,
		func(c *gin.Context, account core.Account) {
			if account.Role != constants.ROOT {
				code.GinHandler(c, code.ErrPermissionDenied)
				return
			}
			next(c, account)
		})
}

func GetAccount(c *gin.Context, next func(c *gin.Context, account core.Account)) {
	auth := Authorize(c)
	if auth.Status != Authorized {
		code.GinHandler(c, code.ErrUnauthorized)
		return
	}
	var account core.Account
	DB.Find(&account, auth.AccountId)
	next(c, account)
}

func Authorize(c *gin.Context) Authentication {
	var headers Headers
	c.ShouldBindHeader(&headers)
	authorization := headers.Authorization
	splited := strings.Split(authorization, " ")
	if authorization == "" || len(splited) != 2 {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return AuthorizeByJWT(splited[1])
}

func AuthorizeByJWT(token string) Authentication {
	claims, err := utils.VerifyJwt(token)
	if err != nil {
		return Authentication{
			Status: Unauthorized,
		}
	}
	return Authentication{
		Status:    Authorized,
		Email:     claims["email"].(string),
		AccountId: utils.StringToUint(claims["id"].(string)),
		Role:      constants.Role(claims["role"].(string)),
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func Health(router *gin.Engine) {
	router.GET("/health/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, health{
			Status: "UP",
			Checks: make([]string, 0),
		})
	})
	router.GET("/health/live", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, health{
			Status: "UP",
			Checks: make([]string, 0),
		})
	})
}
