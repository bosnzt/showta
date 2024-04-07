package native

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"showta.cc/app/lib/util"
	"showta.cc/app/storage"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"strings"
)

type Extra struct {
	RootPath string `json:"root_path" required:"true" tip:"true"`
}

func (self *Extra) GetRootPath() string {
	return self.RootPath
}

func init() {
	logic.RegisterEngine(func() storage.Storage {
		return &Native{}
	})
}

type Native struct {
	model.Storage
	Extra
}

var config = storage.Config{
	Name:    "native",
	Direct:  true,
	NoCache: true,
}

func (self *Native) GetConfig() storage.Config {
	return config
}

func (self *Native) AllowCache() bool {
	return !config.NoCache
}

func (self *Native) IsDirect() bool {
	return config.Direct
}

func (self *Native) GetExtra() storage.ExtraItem {
	return &self.Extra
}

func (self *Native) Mount() error {
	exist, err := util.IsDirExist(self.RootPath)
	if !exist {
		return fmt.Errorf("native mount err: %+v", err)
	}
	return nil
}

func (self *Native) Get(rpath string) (info msg.Finfo, err error) {
	apath := self.getApath(rpath)
	fileinfo, err := os.Stat(apath)
	if err != nil {
		err = errors.New("dir err")
		return
	}

	info = &msg.FileInfo{
		Name:     fileinfo.Name(),
		IsFolder: fileinfo.IsDir(),
		Modified: fileinfo.ModTime(),
		Size:     fileinfo.Size(),
	}
	return
}

func (self *Native) List(info msg.Finfo) (list []msg.Finfo, err error) {
	rpath := info.GetPath()
	apath := self.getApath(rpath)
	dir, err := ioutil.ReadDir(apath)
	if err != nil {
		return
	}

	for _, file := range dir {
		list = append(list, &msg.FileInfo{
			Path:     path.Join("/", rpath, file.Name()),
			Name:     file.Name(),
			Size:     file.Size(),
			Modified: file.ModTime(),
			IsFolder: file.IsDir(),
		})
	}

	return
}

func (self *Native) Link(info msg.Finfo) (*msg.LinkInfo, error) {
	rpath := info.GetPath()
	apath := self.getApath(rpath)
	return &msg.LinkInfo{Url: apath}, nil
}

func (self *Native) getApath(rpath string) string {
	mountPath := self.GetData().MountPath
	subpath := strings.TrimPrefix(rpath, mountPath)
	apath := filepath.Join(self.GetRootPath(), subpath)
	apath = filepath.ToSlash(apath)
	return apath
}
