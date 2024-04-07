package showta

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"path/filepath"
	"showta.cc/app/storage"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"strings"
)

type SimpleResp struct {
	Code int `json:"code"`
}

type Extra struct {
	Url       string `json:"url" required:"true" tip:"true"`
	RootPath  string `json:"root_path" required:"true" tip:"true"`
	FolderPwd string `json:"folder_pwd" tip:"true"`
	Username  string `json:"username" required:"true" tip:"true"`
	Password  string `json:"password" required:"true" tip:"true"`
}

func (self *Extra) GetRootPath() string {
	return self.RootPath
}

func init() {
	logic.RegisterEngine(func() storage.Storage {
		return &Showta{}
	})
}

type Showta struct {
	model.Storage
	Extra
}

var config = storage.Config{
	Name: "showta",
}

func (self *Showta) GetConfig() storage.Config {
	return config
}

func (self *Showta) AllowCache() bool {
	return !config.NoCache
}

func (self *Showta) IsDirect() bool {
	return config.Direct
}

func (self *Showta) GetExtra() storage.ExtraItem {
	return &self.Extra
}

func (self *Showta) Mount() error {
	err := self.auth()
	if err != nil {
		return err
	}

	return nil
}

func (self *Showta) Get(rpath string) (info msg.Finfo, err error) {
	apath := self.getApath(rpath)
	var result msg.GenericResp[msg.GetFileResp]
	err = self.remote("/file/get", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]string{
			"rpath":    apath,
			"password": self.FolderPwd,
		})
	}, true)

	if err != nil {
		return
	}

	if result.Code != 0 {
		err = errors.New(result.Msg)
		return
	}

	result.Data.FileInfo.RawUrl = result.Data.RawUrl
	info = &result.Data.FileInfo
	return
}

func (self *Showta) List(info msg.Finfo) (list []msg.Finfo, err error) {
	rpath := info.GetPath()
	mountPath := self.GetData().MountPath
	apath := self.getApath(rpath)

	var result msg.GenericResp[msg.ListFileResp]
	err = self.remote("/file/list", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]string{
			"rpath":    apath,
			"password": self.FolderPwd,
		})
	}, true)

	if err != nil {
		return
	}

	if result.Code != 0 {
		err = errors.New(result.Msg)
		return
	}

	for _, v := range result.Data.List {
		subpath := strings.TrimPrefix(v.Path, self.GetRootPath())
		apath := filepath.Join("/", mountPath, subpath)
		apath = filepath.ToSlash(apath)
		list = append(list, &msg.FileInfo{
			Path:     apath,
			Name:     v.Name,
			Size:     v.Size,
			Modified: v.Modified,
			IsFolder: v.IsFolder,
		})
	}

	return
}

func (self *Showta) Link(info msg.Finfo) (*msg.LinkInfo, error) {
	rpath := info.GetPath()
	apath := self.getApath(rpath)
	var result msg.GenericResp[msg.GetFileResp]
	err := self.remote("/file/get", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]string{
			"rpath":    apath,
			"password": self.FolderPwd,
		})
	}, true)

	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		err = errors.New(result.Msg)
		return nil, err
	}

	return &msg.LinkInfo{Url: result.Data.RawUrl}, nil
}

func (self *Showta) remote(api string, callback func(req *resty.Request), refresh bool) error {
	req := resty.New().R()
	req.SetHeader("Authorization", self.Token)
	callback(req)
	resp, err := req.Execute(http.MethodPost, self.Url+api)
	if err != nil {
		return err
	}

	var simpleResp SimpleResp
	err = json.Unmarshal(resp.Body(), &simpleResp)
	if err != nil {
		return errors.New("remote decode error")
	}

	if refresh && simpleResp.Code == http.StatusUnauthorized {
		err = self.auth()
		if err != nil {
			return err
		}

		return self.remote(api, callback, false)
	}

	return nil
}

func (self *Showta) auth() error {
	var result msg.GenericResp[msg.LoginResp]
	req := resty.New().R()
	req.SetResult(&result).SetBody(map[string]string{
		"username": self.Username,
		"password": self.Password,
	})
	_, err := req.Execute(http.MethodPost, self.Url+"/admin/login")
	if err != nil {
		return err
	}

	if result.Code != 0 {
		return errors.New(result.Msg)
	}

	if result.Data.Token == "" {
		return errors.New("auth url error")
	}

	logic.SyncUpdateStorage(self, result.Data.Token)

	return nil
}

func (self *Showta) getApath(rpath string) string {
	mountPath := self.GetData().MountPath
	subpath := strings.TrimPrefix(rpath, mountPath)
	apath := filepath.Join(self.GetRootPath(), subpath)
	apath = filepath.ToSlash(apath)
	return apath
}
