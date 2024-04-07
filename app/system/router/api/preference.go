package api

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/msg"
)

func AddRouterPreference(g *gin.RouterGroup) {
	group := g.Group("/preference")
	group.GET("/display_get", getDisplay)
	group.POST("/display_update", updateDisplay)
	group.GET("/site_get", getSite)
	group.POST("/site_update", updateSite)
}

func GetPreference(c *gin.Context) {
	data, err := logic.GetPreference()
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func getDisplay(c *gin.Context) {
	data, err := logic.GetDisplay()
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func updateDisplay(c *gin.Context) {
	var req msg.UpdateDisplayReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.UpdateDisplay(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, nil)
}

func getSite(c *gin.Context) {
	data, err := logic.GetSite()
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	msg.Response(c, data)
}

func updateSite(c *gin.Context) {
	var req msg.UpdateSiteReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		msg.RespError(c, http.StatusBadRequest, err)
		return
	}

	err = logic.UpdateSite(c, req)
	if err != nil {
		msg.RespError(c, http.StatusInternalServerError, err)
		return
	}

	UpdateHtml()
	msg.Response(c, nil)
}
