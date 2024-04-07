package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/msg"
)

func AddRouterFile(g *gin.RouterGroup) {
	group := g.Group("/file")
	group.POST("/list", listFile)
	group.POST("/get", getFile)
	group.POST("/subdir", subdir)
}

func listFile(c *gin.Context) {
	var req msg.ListFileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	res, err := logic.IsFolderForbidden(c, req.Rpath, *req.Password)
	if err != nil {
		if res {
			msg.RespError(c, http.StatusForbidden, err)
		} else {
			msg.RespError(c, http.StatusInternalServerError, err)
		}

		return
	}

	data, err := logic.ViewListFile(c, req.Rpath)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func getFile(c *gin.Context) {
	var req msg.GetFileReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	res, err := logic.IsFolderForbidden(c, req.Rpath, *req.Password)
	if err != nil {
		if res {
			msg.RespError(c, http.StatusForbidden, err)
		} else {
			msg.RespError(c, http.StatusInternalServerError, err)
		}

		return
	}

	data, err := logic.GetStorageFile(c, req.Rpath)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func ProxyFile(c *gin.Context) {
	rpath := c.Param("path")
	logic.ProxyFile(c.Request, c.Writer, rpath)
}

func subdir(c *gin.Context) {
	var req msg.SubdirReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	data, err := logic.Subdir(c, req.Rpath)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}
