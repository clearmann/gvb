package model

import "gorm.io/gorm"

const (
	STATUS_PUBLIC = iota + 1 // 公开
	STATUS_SECRET            // 私密
	STATUS_DRAFT             //草稿
)
const (
	TYPE_ORIGINAL  = iota + 1 // 原创
	TYPE_REPRINT              //转载
	TYPE_TRANSLATE            //翻译
)

// 一个文章属于一个分类
// 一个文章属于一个用户
// 一个文章可以有多个标签 一个标签也可以有多个文章
type Article struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null" json:"title"` //标题
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"` // 1-原创 2-转载 3-翻译
	Status      int    `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`    // 1-公开 2-私密
	IsTop       bool   `json:"is_top"`                                              //是否置顶
	IsDelete    bool   `json:"is_delete"`                                           //是否删除
	OriginalUrl string `json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id

	Tags     []Tag     `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
	Category *Category `gorm:"foreignkey:CategoryId" json:"category"`
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}
type ArticleTag struct {
	ArticleId uint
	Tag       uint
}
