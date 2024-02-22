package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name    string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Article []Article `gorm:"many2many:article_tag" json:"articles,omitempty"`
}
