package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(inject context.Context) {
	db = inject.Value("db").(*gorm.DB)
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
	db.Find(&account, auth.AccountId)
	next(c, account)
}

func Authorize(c *gin.Context) Authentication {
	var headers Headers
	c.ShouldBindHeader(&headers)
	authorization := headers.Authorization
	splited := strings.Split(authorization, " ")
	fmt.Println(splited)
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
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
