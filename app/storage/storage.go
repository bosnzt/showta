package storage

import (
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
)

type Config struct {
	Name    string `json:"name"`
	Direct  bool   `json:"direct"`
	NoCache bool
}

type ExtraItem interface{}

type Storage interface {
	GetConfig() Config
	SetData(model.Storage)
	GetData() *model.Storage
	GetExtra() ExtraItem
	Mount() error
	List(info msg.Finfo) (list []msg.Finfo, err error)
	Link(info msg.Finfo) (*msg.LinkInfo, error)
	AllowCache() bool
	IsDirect() bool
}

type FormItem struct {
	Name     string `json:"name"`
	Etype    string `json:"etype"`
	Dvalue   string `json:"dvalue"`
	Options  string `json:"options"`
	Required bool   `json:"required"`
	Tip      bool   `json:"tip"`
}

type Form struct {
	Extra []FormItem `json:"extra"`
}

type Getter interface {
	Get(rpath string) (info msg.Finfo, err error)
}
