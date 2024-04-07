package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "showta.cc/app/storage/engine"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

func AddRouterStorage(g *gin.RouterGroup) {
	group := g.Group("/storage")
	group.GET("/list", listStorage)
	group.POST("/get", getStorage)
	group.POST("/mount", mountStorage)
	group.POST("/update", updateStorage)
	group.POST("/switch", switchStorage)
	group.POST("/delete", deleteStorage)

	group.GET("/listname", listEngineName)
	group.GET("/listform", listEngineForm)
}

func listStorage(c *gin.Context) {
	list, err := model.GetAllStorage()
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, list)
}

func getStorage(c *gin.Context) {
	var req model.Storage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	data, err := model.GetStorage(req.ID)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func mountStorage(c *gin.Context) {
	var req model.Storage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.MountStorage(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func updateStorage(c *gin.Context) {
	var req model.Storage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.UpdateStorage(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func switchStorage(c *gin.Context) {
	var req model.Storage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.SwitchStorage(c, req.ID)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func deleteStorage(c *gin.Context) {
	var req model.Storage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.DeleteStorage(c, req.ID)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func listEngineName(c *gin.Context) {
	resp := logic.GetAllEngineName()
	msg.Response(c, resp)
}

func listEngineForm(c *gin.Context) {
	resp := logic.GetAllEngineForm()
	msg.Response(c, resp)
}
