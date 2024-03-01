package model

import "gorm.io/gorm"

type Config struct {
	Model
	Key   string `gorm:"unique;type:varchar(256)" json:"key"`
	Value string `gorm:"type:varchar(256)" json:"value"`
	Desc  string `gorm:"type:varchar(256)" json:"desc"`
}

func GetConfigByKey(db *gorm.DB, key string) string {
	var config Config
	result := db.Where("key", key).First(&config)
	if result.Error != nil {
		return ""
	}
	return config.Value
}
func GetConfigMap(db *gorm.DB) (map[string]string, error) {
	var configs []Config
	result := db.Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	m := make(map[string]string)
	for _, config := range configs {
		m[config.Key] = config.Value
	}
	return m, nil
}
func GetConfigBool(db *gorm.DB, key string) bool {
	var config Config
	result := db.Where("key", key).First(&config)
	if result.Error != nil {
		return false
	}
	return config.Value == "true"
}
