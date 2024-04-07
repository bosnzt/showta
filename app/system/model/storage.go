package model

import (
	"time"
)

type Storage struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	MountPath string `json:"mount_path" gorm:"unique"`
	Engine    string `json:"engine"`
	Status    string `json:"status"`
	Extra     string `json:"extra"`
	Disabled  bool   `json:"disabled"`
	Remark    string `json:"remark"`
	Token     string `json:"-"`
	UpdatedAt time.Time
}

func (self *Storage) SetData(data Storage) {
	*self = data
}

func (self *Storage) GetData() *Storage {
	return self
}

func (self *Storage) SetStatus(status string) {
	self.Status = status
}

func GetStorage(id uint) (*Storage, error) {
	var data Storage
	if err := db.First(&data, id).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func GetAllStorage() ([]Storage, error) {
	var dataList []Storage
	err := db.Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func CreateStorage(data *Storage) error {
	return db.Create(data).Error
}

func UpdateStorage(data *Storage) error {
	return db.Save(&data).Error
}

func DeleteStorage(id uint) error {
	return db.Delete(&Storage{}, id).Error
}

func GetAllEnableStorage() ([]Storage, error) {
	var dataList []Storage
	err := db.Where("disabled = ?", false).Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}
