package baidunetdisk

import (
	"showta.cc/app/storage"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

type Extra struct {
}

func init() {
	logic.RegisterEngine(func() storage.Storage {
		return &Baidunetdisk{}
	})
}

type Baidunetdisk struct {
	model.Storage
	Extra
}

var config = storage.Config{
	Name: "baidunetdisk",
}

func (self *Baidunetdisk) GetConfig() storage.Config {
	return config
}

func (self *Baidunetdisk) AllowCache() bool {
	return !config.NoCache
}

func (self *Baidunetdisk) IsDirect() bool {
	return config.Direct
}

func (self *Baidunetdisk) GetExtra() storage.ExtraItem {
	return &self.Extra
}

func (self *Baidunetdisk) Mount() error {
	return nil
}

func (self *Baidunetdisk) List(info msg.Finfo) (list []msg.Finfo, err error) {
	return
}

func (self *Baidunetdisk) Link(info msg.Finfo) (*msg.LinkInfo, error) {
	return nil, nil
}
