package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"showta.cc/app/lib/util"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

func AddFolderSetting(ctx context.Context, data model.FolderSetting) error {
	err := model.CreateFolderSetting(&data)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFolderSetting(ctx context.Context, data model.FolderSetting) error {
	_, err := model.GetFolderSetting(data.ID)
	if err != nil {
		return err
	}

	err = model.UpdateFolderSetting(&data)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFolderSetting(ctx context.Context, id uint) error {
	_, err := model.GetFolderSetting(id)
	if err != nil {
		return err
	}

	err = model.DeleteFolderSetting(id)
	if err != nil {
		return err
	}

	return nil
}

func IsFolderForbidden(c *gin.Context, rpath string, password string) (res bool, err error) {
	user := c.MustGet("identity").(*model.User)
	if user.IsSuper() {
		return
	}

	rpath = util.StandardPath(rpath)
	setting, err := model.GetFolderSettingByFolder(rpath)
	if err != nil {
		return
	}

	if setting.Password != "" && setting.Password != password {
		res = true
		err = msg.ErrAccessPwd
		return
	}

	return
}
