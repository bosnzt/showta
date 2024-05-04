package msg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/system/model"
	"time"
)

type GenericResp[V any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data V      `json:"data"`
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RespError(c *gin.Context, errCode int, errMsg error) {
	c.JSON(http.StatusOK, Resp{
		Code: errCode,
		Msg:  errMsg.Error(),
	})
}

func Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type AboutUserResp struct {
	*model.User
	ClientIp string `json:"client_ip"`
}

type ResetPwdReq struct {
	Password string `json:"password"`
}

type ListUserReq struct {
	Query    string `form:"query"`
	Pagenum  int    `form:"pagenum"`
	Pagesize int    `form:"pagesize"`
}

type ListUserResp struct {
	Total    int          `json:"total"`
	Pagenum  int          `json:"pagenum"`
	UserList []model.User `json:"users"`
}

type Finfo interface {
	GetFileId() string
	GetPath() string
	IsDir() bool
	GetName() string
	GetSize() int64
	ModTime() time.Time
	GetRaw() string
}

type GetFileReq struct {
	Rpath    string  `json:"rpath" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

type GetFileResp struct {
	FileInfo
	RawUrl string `json:"raw_url"`
}

type ListFileReq struct {
	Rpath    string  `json:"rpath" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

type ListFileResp struct {
	List   []FileInfo `json:"list"`
	Topmd  string     `json:"topmd"`
	Readme string     `json:"readme"`
}

type FileInfo struct {
	FileId   string    `json:"file_id"`
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	IsFolder bool      `json:"is_folder"`
	Ptype    int       `json:"ptype"`
	RawUrl   string    `json:"raw_url"`
}

func (self *FileInfo) GetFileId() string {
	return self.FileId
}

func (self *FileInfo) GetPath() string {
	return self.Path
}

func (self *FileInfo) IsDir() bool {
	return self.IsFolder
}

func (self *FileInfo) GetName() string {
	return self.Name
}

func (self *FileInfo) GetSize() int64 {
	return self.Size
}

func (self *FileInfo) ModTime() time.Time {
	return self.Modified
}

func (self *FileInfo) GetRaw() string {
	return self.RawUrl
}

type LinkInfo struct {
	Url    string
	Expire time.Duration
}

type SubdirReq struct {
	Rpath string `json:"rpath" binding:"required"`
}

type DisplayTemplate struct {
	Video   []string `json:"video"`
	Picture []string `json:"picture"`
	Text    []string `json:"text"`
	Audio   []string `json:"audio"`
}

type GetDisplayResp struct {
	Template DisplayTemplate    `json:"template"`
	Data     []model.Preference `json:"data"`
}

type UpdateDisplayReq struct {
	Video   []string `json:"video"`
	Picture []string `json:"picture"`
	Text    []string `json:"text"`
	Audio   []string `json:"audio"`
	Office  string   `json:"office"`
}

type GetSiteResp struct {
	Data []model.Preference `json:"data"`
}

type UpdateSiteReq struct {
	Title          string `json:"title"`
	Logo           string `json:"logo"`
	Favicon        string `json:"favicon"`
	Domain         string `json:"domain"`
	Notice         string `json:"notice"`
	GlobalSign     bool   `json:"global_sign"`
	SignExpiration int64  `json:"sign_expiration"`
}
