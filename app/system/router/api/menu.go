package api

import (
	"github.com/gin-gonic/gin"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

func GetMenu(c *gin.Context) {
	user := c.MustGet("identity").(*model.User)
	if user.IsSuper() {
		msg.Response(c, conf.AdminMenuList)
	} else {
		msg.Response(c, conf.CommonMenuList)
	}
}
