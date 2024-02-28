package model

type Tag struct {
	Model
	Name    string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Article []Article `gorm:"many2many:article_tag" json:"articles,omitempty"`
}
