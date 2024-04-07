package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/orcaman/concurrent-map/v2"
	"reflect"
	"showta.cc/app/lib/util"
	"showta.cc/app/storage"
	"showta.cc/app/system/log"
	"showta.cc/app/system/model"
	"sync"
)

const (
	WORK     = "work"
	DISABLED = "disabled"
)

type StorageFunc func() storage.Storage

var (
	engineMap     = cmap.New[StorageFunc]()
	engineFormMap = cmap.New[storage.Form]()
	storageMap    sync.Map
)

func RegisterEngine(engine StorageFunc) {
	inst := engine()
	instCfg := inst.GetConfig()
	engineMap.Set(instCfg.Name, engine)
	formExtra := getEngineFormExtra(inst.GetExtra())
	engineFormMap.Set(instCfg.Name, storage.Form{
		Extra: formExtra,
	})
}

func GetAllEngineForm() map[string]storage.Form {
	return engineFormMap.Items()
}

func getEngineFormExtra(extra storage.ExtraItem) []storage.FormItem {
	elem := reflect.TypeOf(extra).Elem()
	itemList := make([]storage.FormItem, 0)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag
		name, ok := tag.Lookup("json")
		if !ok {
			continue
		}

		elementType := tag.Get("etype")
		if elementType == "" {
			elementType = "input"
		}

		item := storage.FormItem{
			Name:     name,
			Etype:    elementType,
			Dvalue:   tag.Get("dvalue"),
			Options:  tag.Get("options"),
			Required: tag.Get("required") == "true",
			Tip:      tag.Get("tip") == "true",
		}
		itemList = append(itemList, item)
	}

	return itemList
}

func GetEngine(name string) (StorageFunc, error) {
	n, ok := engineMap.Get(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf("no storage named: %s", name))
	}
	return n, nil
}

func GetAllEngineName() []string {
	var nameList []string
	engineMap.IterCb(func(k string, v StorageFunc) {
		nameList = append(nameList, k)
	})

	return nameList
}

func loadAllStorage() {
	dataList, err := model.GetAllEnableStorage()
	if err != nil {
		log.StdErrorf("failed get enabled storages: %+v", err)
		return
	}

	var success int
	for _, data := range dataList {
		err = LoadStorage(context.Background(), data)
		if err != nil {
			log.StdErrorf("load storage: [%s], engine: [%s], %+v", data.MountPath, data.Engine, err)
		} else {
			success++
		}
	}

	log.StdInfof("load storage success: [%d/%d]", success, len(dataList))

}

func MountStorage(ctx context.Context, data model.Storage) error {
	data.MountPath = util.StandardPath(data.MountPath)
	engine, err := GetEngine(data.Engine)
	if err != nil {
		return err
	}

	data.SetStatus(WORK)

	inst := engine()
	if err != nil {
		return err
	}

	err = initStorage(ctx, data, inst)
	if err != nil {
		return err
	}

	jsonData, _ := json.Marshal(inst.GetExtra())
	data.Extra = string(jsonData)
	err = model.CreateStorage(&data)
	if err != nil {
		return err
	}

	return nil
}

func UpdateStorage(ctx context.Context, data model.Storage) error {
	oldData, err := model.GetStorage(data.ID)
	if err != nil {
		return err
	}

	data.MountPath = util.StandardPath(data.MountPath)
	err = model.UpdateStorage(&data)
	if err != nil {
		return err
	}

	if data.Disabled {
		return nil
	}

	inst, err := GetStorageByMountPath(oldData.MountPath)
	if err != nil {
		return err
	}

	if data.MountPath != oldData.MountPath {
		storageMap.Delete(oldData.MountPath)
	}

	err = initStorage(ctx, data, inst)
	return err
}

func SwitchStorage(ctx context.Context, id uint) error {
	data, err := model.GetStorage(id)
	if err != nil {
		return err
	}

	if data.Disabled {
		//Enable
		data.Disabled = false
		data.SetStatus(WORK)
		err = model.UpdateStorage(data)
		if err != nil {
			return err
		}

		err = LoadStorage(ctx, *data)
		if err != nil {
			return err
		}
	} else {
		//Disable
		data.Disabled = true
		data.SetStatus(DISABLED)
		err = model.UpdateStorage(data)
		if err != nil {
			return err
		}

		storageMap.Delete(data.MountPath)
	}

	return nil
}

func DeleteStorage(ctx context.Context, id uint) error {
	data, err := model.GetStorage(id)
	if err != nil {
		return err
	}

	if !data.Disabled {
		storageMap.Delete(data.MountPath)
	}

	err = model.DeleteStorage(id)
	if err != nil {
		return err
	}

	return nil
}

func initStorage(ctx context.Context, data model.Storage, inst storage.Storage) (err error) {
	inst.SetData(data)
	err = json.Unmarshal([]byte(data.Extra), inst.GetExtra())
	if err != nil {
		return
	}

	err = inst.Mount()
	if err != nil && data.ID == 0 {
		return
	}

	storageMap.Store(data.MountPath, inst)
	status := WORK
	if err != nil {
		status = err.Error()
	}

	if status != data.Status {
		data.SetStatus(status)
		model.UpdateStorage(&data)
	}

	return
}

func GetStorageByMountPath(mountPath string) (storage.Storage, error) {
	data, ok := storageMap.Load(mountPath)
	if !ok {
		return nil, errors.New("no mount path for an storage is: " + mountPath)
	}

	return data.(storage.Storage), nil
}

func LoadStorage(ctx context.Context, data model.Storage) error {
	name := data.Engine
	engine, ok := engineMap.Get(name)
	if !ok {
		return fmt.Errorf("failed get engine:%s", name)
	}

	inst := engine()
	err := initStorage(ctx, data, inst)
	return err
}

func SyncUpdateStorage(inst storage.Storage, token string) {
	data := inst.GetData()
	if data.ID > 0 && data.Token != token {
		data.Token = token
		model.UpdateStorage(data)
	}
}
