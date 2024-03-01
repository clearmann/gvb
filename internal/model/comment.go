package model

import "gorm.io/gorm"

const (
	TYPE_ARTICLE = iota + 1 //文章
	TYPE_LINK               //友链
	TYPE_TALK               //说说
)

/*
如果评论类型是文章，那么 topic_id 就是文章的 id
如果评论类型是友链，不需要 topic_id
*/
type Comment struct {
	Model
	UserId      int    `json:"user_id"`       // 评论者
	ReplyUserId int    `json:"reply_user_id"` // 被回复者
	TopicId     int    `json:"topic_id"`      // 评论的文章
	ParentId    int    `json:"parent_id"`     // 父评论
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
	IsReview    bool   `json:"is_review"`

	// Belongs To
	User      *UserAuth `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article   *Article  `gorm:"foreignKey:TopicId" json:"article"`
}

// 获取某篇文章的评论数
func GetArticleCommentCount(db *gorm.DB, articleId int) (count int64, err error) {
	result := db.Model(&Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return count, result.Error
}

// 新增评论
func AddComment(db *gorm.DB, userId, typ, topicId int, content string, isReview bool) (*Comment, error) {
	comment := Comment{
		UserId:   userId,
		TopicId:  topicId,
		Content:  content,
		Type:     typ,
		IsReview: isReview,
	}
	result := db.Create(&comment)
	return &comment, result.Error
}

// 回复评论
func ReplyComment(db *gorm.DB, userId, replyUserId, parentId int, content string, isReview bool) (*Comment, error) {
	var parent Comment
	result := db.First(&parent, parentId)
	if result.Error != nil {
		return nil, result.Error
	}

	comment := Comment{
		UserId:      userId,
		Content:     content,
		ReplyUserId: replyUserId,
		ParentId:    parentId,
		IsReview:    isReview,
		TopicId:     parent.TopicId, // 主题和父评论一样
		Type:        parent.Type,    // 类型和父评论一样
	}
	result = db.Create(&comment)
	return &comment, result.Error
}
