package model

import (
	"time"
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username" gorm:"unique" binding:"required"`
	Password   string    `json:"password"`
	EncryptPwd string    `json:"-"`
	PwdStamp   int64     `json:"-"`
	Salt       string    `json:"-"`
	Role       int       `json:"role"`
	Enable     bool      `json:"enable"`
	Perm       int       `json:"perm"`
	LoginIp    string    `json:"login_ip"`
	LoginTime  time.Time `json:"login_time"`
	UpdatedAt  time.Time `json:"modified"`
}

func (self *User) IsSuper() bool {
	return self.Username == "admin"
}

func GetAllUser() ([]User, error) {
	var dataList []User
	err := db.Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func GeUserListBySearch(username string, limit int, offset int) ([]User, error) {
	var dataList []User
	tx := db.Limit(limit).Offset(offset)
	if username != "" {
		tx = tx.Where("username LIKE ?", "%"+username+"%")
	}

	err := tx.Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func CountUser() (int, error) {
	var count int64
	err := db.Model(&User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func CountUserBySearch(username string) (int, error) {
	var count int64
	tx := db.Model(&User{})
	if username != "" {
		tx = tx.Where("username LIKE ?", "%"+username+"%")
	}

	err := tx.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func CreateUser(data *User) error {
	return db.Create(data).Error
}

func BatchCreateUser(data []*User) error {
	return db.Create(data).Error
}

func GetUser(id uint) (*User, error) {
	var data User
	if err := db.First(&data, id).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func GetUserByName(username string) (*User, error) {
	var data User
	if err := db.Where("username = ?", username).Limit(1).Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func UpdateUser(data *User) error {
	return db.Save(&data).Error
}

func DeleteUser(id uint) error {
	return db.Delete(&User{}, id).Error
}
