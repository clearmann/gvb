package model

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Name  string `gorm:"unique;type:varchar(64)" json:"name"`
	Label string `gorm:"unique;type:varchar(64)" json:"label"`
	Cover string `gorm:"unique;type:varchar(64)" json:"cover"`
}
