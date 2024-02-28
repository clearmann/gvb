package model

import (
	"gorm.io/gorm"
	"time"
)

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
	Model
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
	ArticleId int
	TagId     int
}
type BlogArticleVO struct {
	Article

	CommentCount int64 `json:"comment_count"` // 评论数量
	LikeCount    int64 `json:"like_count"`    // 点赞数量
	ViewCount    int64 `json:"view_count"`    // 访问数量

	LastArticle       ArticlePaginationVO  `gorm:"-" json:"last_article"`       // 上一篇
	NextArticle       ArticlePaginationVO  `gorm:"-" json:"next_article"`       // 下一篇
	RecommendArticles []RecommendArticleVO `gorm:"-" json:"recommend_articles"` // 推荐文章
	NewestArticles    []RecommendArticleVO `gorm:"-" json:"newest_articles"`    // 最新文章
}
type ArticlePaginationVO struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}
type RecommendArticleVO struct {
	ID        int       `json:"id"`
	Img       string    `json:"img"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func GetArticleList(db *gorm.DB, page, size int, categoryId, tagId int) (data []Article, total int64, err error) {
	db = db.Model(Article{})
	// 文章未删除 且 帖子状态为公开
	db = db.Where("is_delete = 0 and status = 1")
	if categoryId != 0 {
		db = db.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		db.Where("id in (select article_id from article_tag where tag_id = ?)", tagId)
	}
	db = db.Count(&total)
	result := db.Preload("Tags").
		Preload("Category").
		Order("is_top desc, id desc").
		Scopes(Paginate(page, size)).
		Find(&data)
	return data, total, result.Error
}

// 前台文章详情（不在回收站并且状态为公开）
func GetBlogArticle(db *gorm.DB, id int) (data *Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(Article{Model: Model{ID: id}}).
		Where("is_delete = 0 AND status = 1"). // *
		First(&data)
	return data, result.Error
}

// 查询 n 篇推荐文章 (根据标签)
func GetRecommendList(db *gorm.DB, id, n int) (list []RecommendArticleVO, err error) {
	// sub1: 查出标签id列表
	// SELECT tag_id FROM `article_tag` WHERE `article_id` = ?
	sub1 := db.Table("article_tag").
		Select("tag_id").
		Where("article_id", id)
	// sub2: 查出这些标签对应的文章id列表 (去重, 且不包含当前文章)
	// SELECT DISTINCT article_id FROM (sub1) t
	// JOIN article_tag t1 ON t.tag_id = t1.tag_id
	// WHERE `article_id` != ?
	sub2 := db.Table("(?) t1", sub1).
		Select("DISTINCT article_id").
		Joins("JOIN article_tag t ON t.tag_id = t1.tag_id").
		Where("article_id != ?", id)
	// 根据 文章id列表 查出文章信息 (前 n 个)
	result := db.Table("(?) t2", sub2).
		Select("id, title, img, created_at").
		Joins("JOIN article a ON t2.article_id = a.id").
		Where("a.is_delete = 0").
		Order("is_top, id DESC").
		Limit(n).
		Find(&list)
	return list, result.Error
}

// 查询最新的 n 篇文章
func GetNewestList(db *gorm.DB, n int) (data []RecommendArticleVO, err error) {
	result := db.Model(&Article{}).
		Select("id, title, img, created_at").
		Where("is_delete = 0 AND status = 1").
		Order("created_at DESC, id ASC").
		Limit(n).
		Find(&data)
	return data, result.Error
}

// 查询上一篇文章 (id < 当前文章 id)
func GetLastArticle(db *gorm.DB, id int) (val ArticlePaginationVO, err error) {
	sub := db.Table("article").Select("max(id)").Where("id < ?", id)
	result := db.Table("article").
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id = (?)", sub).
		Limit(1).
		Find(&val)
	return val, result.Error
}

// 查询下一篇文章 (id > 当前文章 id)
func GetNextArticle(db *gorm.DB, id int) (data ArticlePaginationVO, err error) {
	result := db.Model(&Article{}).
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id > ?", id).
		Limit(1).
		Find(&data)
	return data, result.Error
}

// 前台文章列表（不在回收站并且状态为公开）
func GetBlogArticleList(db *gorm.DB, page, size, categoryId, tagId int) (data []Article, total int64, err error) {
	db = db.Model(&Article{})
	// 未删除 且 状态公开
	db = db.Where("is_delete = 0 and status = 1")
	if categoryId != 0 {
		db.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		db.Where("id in (select article_id from article_tag where tag_id = ?)", tagId)
	}
	db = db.Count(&total)
	result := db.Preload("Tags").Preload("Category").Order("is_top desc, id desc").
		Scopes(Paginate(page, size)).Find("&date")
	return data, total, result.Error
}
