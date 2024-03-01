package model

import (
	"gorm.io/gorm"
)

type UserInfo struct {
	Model
	Nickname string `json:"nickname" gorm:"unique;type:varchar(30);not null"`
	Avatar   string `json:"avatar" gorm:"type:varchar(1024);not null"`
	Intro    string `json:"intro" gorm:"type:varchar(255)"`
	Website  string `json:"website" gorm:"type:varchar(255)"`
}
type UserInfoV0 struct {
	UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}

func GetUserInfoByUserInfoId(db *gorm.DB, userInfoId int) (*UserInfo, error) {
	var userInfo UserInfo
	result := db.Model(&userInfo).Where("id", userInfoId).First(&userInfo)
	return &userInfo, result.Error
}
func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	userInfo := UserInfo{
		Model:    Model{ID: id},
		Nickname: nickname,
		Avatar:   avatar,
		Intro:    intro,
		Website:  website,
	}

	result := db.Select("nickname", "avatar", "intro", "website").Updates(userInfo)
	return result.Error
}
