package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"path"
	"showta.cc/app/lib/util"
	"showta.cc/app/system/log"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"sync"
)

var (
	pwdSettingMap sync.Map
)

func loadAllFolderPwd() {
	dataList, err := model.GetAllFolderSetting()
	if err != nil {
		log.Error(err)
		return
	}

	for _, data := range dataList {
		if data.Password != "" {
			pwdSettingMap.Store(data.Folder, data)
		}
	}
}

func AddFolderSetting(ctx context.Context, data model.FolderSetting) error {
	err := model.CreateFolderSetting(&data)
	if err != nil {
		return err
	}

	if data.Password != "" {
		pwdSettingMap.Store(data.Folder, data)
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

	if data.Password != "" {
		pwdSettingMap.Store(data.Folder, data)
	}

	return nil
}

func DeleteFolderSetting(ctx context.Context, id uint) error {
	data, err := model.GetFolderSetting(id)
	if err != nil {
		return err
	}

	err = model.DeleteFolderSetting(id)
	if err != nil {
		return err
	}

	pwdSettingMap.Delete(data.Folder)

	return nil
}

func IsFolderForbidden(c *gin.Context, rpath string, password string) (res bool, err error) {
	user := c.MustGet("identity").(*model.User)
	if user.IsSuper() {
		return
	}

	rpath = util.StandardPath(rpath)
	setting := findMatchSetting(rpath)
	if setting.Folder == "" || (setting.Folder != rpath && !setting.ApplySub) {
		return
	}

	if setting.Password != password {
		res = true
		err = msg.ErrAccessPwd
		return
	}

	return
}

func findMatchSetting(rpath string) (setting model.FolderSetting) {
	data, ok := pwdSettingMap.Load(rpath)
	if ok {
		return data.(model.FolderSetting)
	}

	prevPath := path.Dir(rpath)
	if prevPath == "/" {
		return
	}

	return findMatchSetting(prevPath)
}
