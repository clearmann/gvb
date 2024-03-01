package model

import "gorm.io/gorm"

// 一个分类下可以有多篇文章，一篇文章只能属于一个分类 为一对多的关系
type Category struct {
	Model
	Name    string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Article []Article `gorm:"foreignKey:CategoryId"`
}
type CategoryVO struct {
	Category
	ArticleCount int `json:"article_count"`
}

// 获取分类列表
func GetCategoryList(db *gorm.DB, number, size int, keyword string) ([]CategoryVO, int64, error) {
	var list = make([]CategoryVO, 0)
	var total int64
	db = db.Table("category c").Select("c.id", "c.name", "count(c.id) as article_count", "c.updated_at").
		Joins("left join article a on c.id = a.category_id and a.is_delete = 0 and a.status = 1")
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	result := db.Group("c.id").Order("c.updated_at DESC").Count(&total).
		Scopes(Paginate(number, size)).Find(&list)

	return list, total, result.Error
}
