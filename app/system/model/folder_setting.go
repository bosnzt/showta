package model

type FolderSetting struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Folder   string `json:"folder" gorm:"unique"`
	Write    bool   `json:"write"`
	Password string `json:"password"`
	ApplySub bool   `json:"apply_sub"`
	Topmd    string `json:"topmd"`
	Readme   string `json:"readme"`
}

func GetAllFolderSetting() ([]FolderSetting, error) {
	var dataList []FolderSetting
	err := db.Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func GetFolderSetting(id uint) (*FolderSetting, error) {
	var data FolderSetting
	if err := db.First(&data, id).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func GetFolderSettingByFolder(folder string) (*FolderSetting, error) {
	var data FolderSetting
	if err := db.Where("folder = ?", folder).Limit(1).Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateFolderSetting(data *FolderSetting) error {
	return db.Create(data).Error
}

func UpdateFolderSetting(data *FolderSetting) error {
	return db.Save(&data).Error
}

func DeleteFolderSetting(id uint) error {
	return db.Delete(&FolderSetting{}, id).Error
}
