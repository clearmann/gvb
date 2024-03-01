package handle

import (
	"github.com/gin-gonic/gin"
	g "gvb/internal/global"
	"gvb/internal/model"
	"strconv"
)

type User struct {
}

func (*User) GetInfo(c *gin.Context) {
	rdb := GetRDB(c)
	user, err := GetCurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
	}
	userInfoV0 := model.UserInfoV0{UserInfo: *user.UserInfo}
	userInfoV0.ArticleLikeSet, err = rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(user.UserInfoId)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
	}
	userInfoV0.CommentLikeSet, err = rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(user.UserInfoId)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
	}
	ReturnSuccess(c, userInfoV0)
}

// 更新当前用户信息, 不需要传 id, 从 Token 中解析出来
type UpdateCurrentUserReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

func (user *User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
	}

	auth, _ := GetCurrentUserAuth(c)
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}

	ReturnSuccess(c, nil)
}
