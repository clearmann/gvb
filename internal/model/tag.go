package model

import "gorm.io/gorm"

type Tag struct {
	Model
	Name    string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Article []Article `gorm:"many2many:article_tag" json:"articles,omitempty"`
}
type TagVO struct {
	Model
	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

func GetTagList(db *gorm.DB, page, size int, keyword string) (list []TagVO, total int64, err error) {
	db = db.Table("tag t").Joins("left join article_tag at on t.id = at.tag_id").
		Select("t.id", "t.name", "count(at.article_id) sc article_count", "t.created_at", "t.updated_at")

	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	result := db.
		Group("t.id").Order("t.updated_at desc").Count(&total).
		Scopes(Paginate(page, size)).Find(&list)

	return list, total, result.Error
}
