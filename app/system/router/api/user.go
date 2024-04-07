package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

func AddRouterUser(g *gin.RouterGroup) {
	group := g.Group("/user")
	group.POST("/list", listUser)
	group.POST("/enable", enableUser)
	group.POST("/add", addUser)
	group.POST("/update", updateUser)
	group.POST("/delete", deleteUser)
}

func AboutUser(c *gin.Context) {
	data := c.MustGet("identity").(*model.User)
	msg.Response(c, data)
}

func ResetPwd(c *gin.Context) {
	var req msg.ResetPwdReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.ResetPwd(c, req.Password)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func UserLogin(c *gin.Context) {
	var req msg.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	data, err := logic.Auth(c, req.Username, req.Password)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	msg.Response(c, data)
}

func listUser(c *gin.Context) {
	var req msg.ListUserReq
	err := c.ShouldBind(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	data, err := logic.ListUser(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func enableUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBind(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.EnableUser(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func addUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.AddUser(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func updateUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.UpdateUser(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func deleteUser(c *gin.Context) {
	var req model.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.DeleteUser(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}
