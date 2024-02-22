package model

import "gorm.io/gorm"

//一个分类下可以有多篇文章，一篇文章只能属于一个分类 为一对多的关系
type Category struct {
	gorm.Model
	Name    string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Article []Article `gorm:"foreignKey:CategoryId"`
}