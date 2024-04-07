package model

type Preference struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
	Term  int    `json:"-"`
}

func GetAllPreference() ([]Preference, error) {
	var dataList []Preference
	err := db.Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func GetPreferenceByTerm(term int) ([]Preference, error) {
	var dataList []Preference
	err := db.Where("term = ?", term).Find(&dataList).Error
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func CountPreference() (int, error) {
	var count int64
	err := db.Model(&Preference{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func BatchCreatePreference(data []*Preference) error {
	return db.Create(data).Error
}

func UpdatePreference(key, value string) error {
	return db.Model(&Preference{}).Where("key = ?", key).Update("value", value).Error
}
