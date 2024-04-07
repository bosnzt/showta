package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/internal/jwt"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

const (
	Permissive = iota
	Enforcing
	Strict
)

func PermissiveAuth(c *gin.Context) {
	auth(c, Permissive)
}

func EnforcingAuth(c *gin.Context) {
	auth(c, Enforcing)
}

func StrictAuth(c *gin.Context) {
	auth(c, Strict)
}

func auth(c *gin.Context, permission int) {
	token := c.GetHeader("Authorization")
	if token == "" {
		if permission == Permissive {
			user, err := model.GetUserByName("guest")
			if err != nil {
				msg.RespError(c, http.StatusUnauthorized, err)
				c.Abort()
				return
			}

			if !user.Enable {
				msg.RespError(c, http.StatusUnauthorized, msg.ErrMustLogin)
				c.Abort()
				return
			}

			c.Set("identity", user)
			c.Next()
			return
		}

		msg.RespError(c, http.StatusUnauthorized, msg.ErrMustLogin)
		c.Abort()
		return
	}

	appClaims, err := jwt.ParseToken(token)
	if err != nil {
		msg.RespError(c, http.StatusUnauthorized, err)
		c.Abort()
		return
	}

	user, err := model.GetUserByName(appClaims.Username)
	if err != nil {
		msg.RespError(c, http.StatusUnauthorized, err)
		c.Abort()
		return
	}

	if permission == Strict && !user.IsSuper() {
		msg.RespError(c, http.StatusForbidden, fmt.Errorf("no permission"))
		c.Abort()
		return
	}

	if !user.Enable {
		msg.RespError(c, http.StatusUnauthorized, fmt.Errorf("user disabled"))
		c.Abort()
		return
	}

	if appClaims.PwdStamp != user.PwdStamp {
		msg.RespError(c, http.StatusUnauthorized, fmt.Errorf("user status changed"))
		c.Abort()
		return
	}

	c.Set("identity", user)
	c.Next()
}
