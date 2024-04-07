package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "showta.cc/app/system/logic"
    "showta.cc/app/system/model"
    "showta.cc/app/system/msg"
)

func AddRouterFolder(g *gin.RouterGroup) {
    group := g.Group("/folder")
    group.GET("/list_setting", listFolderSetting)
    group.POST("/get_setting", getFolderSetting)
    group.POST("/add_setting", addFolderSetting)
    group.POST("/update_setting", updateFolderSetting)
    group.POST("/delete_setting", deleteFolderSetting)
}

func listFolderSetting(c *gin.Context) {
    list, err := model.GetAllFolderSetting()
    if err != nil {
        msg.RespError(c, http.StatusInternalServerError, err)
        return
    }

    msg.Response(c, list)
}

func getFolderSetting(c *gin.Context) {
    var req model.FolderSetting
    err := c.ShouldBindJSON(&req)
    if err != nil {
        msg.RespError(c, http.StatusBadRequest, err)
        return
    }

    data, err := model.GetFolderSetting(req.ID)
    if err != nil {
        msg.RespError(c, http.StatusInternalServerError, err)
        return
    }

    msg.Response(c, data)
}

func addFolderSetting(c *gin.Context) {
    var req model.FolderSetting
    err := c.ShouldBindJSON(&req)
    if err != nil {
        msg.RespError(c, http.StatusBadRequest, err)
        return
    }

    err = logic.AddFolderSetting(c, req)
    if err != nil {
        msg.RespError(c, http.StatusInternalServerError, err)
        return
    }

    msg.Response(c, nil)
}

func updateFolderSetting(c *gin.Context) {
    var req model.FolderSetting
    err := c.ShouldBindJSON(&req)
    if err != nil {
        msg.RespError(c, http.StatusBadRequest, err)
        return
    }

    err = logic.UpdateFolderSetting(c, req)
    if err != nil {
        msg.RespError(c, http.StatusInternalServerError, err)
        return
    }

    msg.Response(c, nil)
}

func deleteFolderSetting(c *gin.Context) {
    var req model.FolderSetting
    err := c.ShouldBindJSON(&req)
    if err != nil {
        msg.RespError(c, http.StatusBadRequest, err)
        return
    }

    err = logic.DeleteFolderSetting(c, req.ID)
    if err != nil {
        msg.RespError(c, http.StatusInternalServerError, err)
        return
    }

    msg.Response(c, nil)
}
