package model

import "gorm.io/gorm"

func GetResource(db *gorm.DB, uri, method string) (resource Resource, err error) {
	result := db.Where(&Resource{Url: uri, Method: method}).First(&resource)
	return resource, result.Error
}
