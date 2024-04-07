package alipan

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"path/filepath"
	"showta.cc/app/lib/apilimit"
	"showta.cc/app/lib/util"
	"showta.cc/app/storage"
	"showta.cc/app/system/log"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"strings"
	"time"
)

type Extra struct {
	SpaceType    string `json:"space_type" required:"true" etype:"select" options:"resource,backup" tip:"true"`
	RootId       string `json:"root_id" dvalue:"root" required:"true" tip:"true"`
	RefreshToken string `json:"refresh_token" required:"true" etype:"textarea" tip:"true"`
	ClientId     string `json:"client_id" tip:"true"`
	ClientSecret string `json:"client_secret" tip:"true"`
}

var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

func init() {
	logic.RegisterEngine(func() storage.Storage {
		return &Alipan{
			Domain:        "https://openapi.alipan.com",
			ProxyOauthUrl: "https://api.showta.cc/alipan/access_token",
		}
	})
}

type Alipan struct {
	model.Storage
	Extra
	Domain        string
	ProxyOauthUrl string
	DriveId       string
	rateLimiter   *apilimit.ApiRateLimiter
}

var config = storage.Config{
	Name: "alipan",
}

func (self *Alipan) GetConfig() storage.Config {
	return config
}

func (self *Alipan) AllowCache() bool {
	return !config.NoCache
}

func (self *Alipan) IsDirect() bool {
	return config.Direct
}

func (self *Alipan) GetExtra() storage.ExtraItem {
	return &self.Extra
}

func (self *Alipan) Mount() error {
	self.rateLimiter = apilimit.NewApiRateLimiter(map[string]apilimit.ApiLimit{
		"list":           {MaxCount: 40, Interval: 10 * time.Second},
		"getDownloadUrl": {MaxCount: 10, Interval: 10 * time.Second},
		"others":         {MaxCount: 150, Interval: 10 * time.Second},
	})

	var result GetDriveInfoResp
	err := self.remote("/adrive/v1.0/user/getDriveInfo", func(req *resty.Request) {
		req.SetResult(&result)
	}, true)
	if err != nil {
		return err
	}

	if self.SpaceType == "resource" {
		self.DriveId = result.ResourceDriveId
	} else if self.SpaceType == "backup" {
		self.DriveId = result.BackupDriveId
	} else {
		self.DriveId = result.DefaultDriveId
	}

	return nil
}

func (self *Alipan) List(info msg.Finfo) (list []msg.Finfo, err error) {
	if !self.rateLimiter.Allow("list") {
		err = errors.New("too many requests")
		return
	}

	rpath := info.GetPath()
	fileId := info.GetFileId()
	if fileId == "" {
		fileItem, err := self.getByPath(rpath)
		if err != nil {
			return list, err
		}

		fileId = fileItem.FileId
	}

	var result ListResp
	err = self.remote("/adrive/v1.0/openFile/list", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]string{
			"drive_id":       self.DriveId,
			"parent_file_id": fileId,
		})
	}, true)
	if err != nil {
		return
	}

	for _, v := range result.ItemList {
		apath := filepath.Join(rpath, v.Name)
		apath = filepath.ToSlash(apath)
		list = append(list, &msg.FileInfo{
			FileId:   v.FileId,
			Path:     apath,
			Name:     v.Name,
			Size:     v.Size,
			Modified: v.UpdatedAt,
			IsFolder: v.Type == "folder",
		})
	}

	return
}

func (self *Alipan) Link(info msg.Finfo) (*msg.LinkInfo, error) {
	if !self.rateLimiter.Allow("getDownloadUrl") {
		return nil, errors.New("too many requests")
	}

	var result getDownloadUrlResp
	err := self.remote("/adrive/v1.0/openFile/getDownloadUrl", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]interface{}{
			"drive_id":   self.DriveId,
			"file_id":    info.GetFileId(),
			"expire_sec": 900,
		})
	}, true)
	if err != nil {
		return nil, err
	}

	return &msg.LinkInfo{Url: result.Url, Expire: 900 * time.Second}, nil
}

func (self *Alipan) getByPath(rpath string) (result FileItem, err error) {
	mountPath := self.GetData().MountPath
	subpath := strings.TrimPrefix(rpath, mountPath)
	if subpath == "" {
		result.FileId = self.RootId
		return
	}

	if !self.rateLimiter.Allow("others") {
		err = errors.New("too many requests")
		return
	}

	err = self.remote("/adrive/v1.0/openFile/get_by_path", func(req *resty.Request) {
		req.SetResult(&result).SetBody(map[string]string{
			"drive_id":  self.DriveId,
			"file_path": subpath,
		})
	}, true)

	return
}

func (self *Alipan) remote(api string, callback func(req *resty.Request), refresh bool) error {
	var simpleResp SimpleResp
	req := util.HttpClient().R()
	req.SetHeader("Authorization", "Bearer "+self.GetData().Token)
	req.SetHeader("Content-Type", "application/json")
	if callback != nil {
		callback(req)
	}

	req.SetError(&simpleResp)
	resp, err := req.Execute(http.MethodPost, self.Domain+api)
	if err != nil {
		log.Errorf("alipan remote execute err:%+v", err)
		return err
	}

	if simpleResp.Code != "" {
		if !refresh {
			log.Errorf("alipan remote resp err:%+v", simpleResp)
		}

		if resp.StatusCode() > 399 && refresh &&
			(self.GetData().Token == "" || simpleResp.Code == "AccessTokenExpired") {
			err = self.auth()
			if err != nil {
				return err
			}

			return self.remote(api, callback, false)
		}

		return errors.New(simpleResp.Code)
	}

	return nil
}

func (self *Alipan) auth() error {
	var result AccessTokenResp
	url := self.Domain + "/oauth/access_token"
	if self.ClientId == "" {
		url = self.ProxyOauthUrl
	}

	req := util.HttpClient().SetHeader("user-agent", UserAgent).R()
	req.SetResult(&result).SetBody(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": self.RefreshToken,
		"client_id":     self.ClientId,
		"client_secret": self.ClientSecret,
	})

	resp, err := req.Execute(http.MethodPost, url)
	if err != nil {
		log.Errorf("auth err:%+v", err)
		return err
	}

	if result.AccessToken == "" {
		var errResp OauthErrResp
		json.Unmarshal(resp.Body(), &errResp)
		if errResp.Code != "" {
			return errors.New(errResp.Message)
		}

		return errors.New("auth url error")
	}

	logic.SyncUpdateStorage(self, result.AccessToken)

	return nil
}
