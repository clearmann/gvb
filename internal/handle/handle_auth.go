package handle

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	g "gvb/internal/global"
	"gvb/internal/model"
	"gvb/internal/utils/jwt"
	"regexp"
	"strconv"
)

type UserAuth struct{}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}
type LoginVO struct {
	model.UserInfo

	// 点赞 Set: 用于记录用户点赞过的文章, 评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

func isValidEmail(email string) bool {
	// 正则表达式检查邮箱格式
	// 这里使用一个简单的正则表达式，实际应用中可能需要更复杂的规则
	// 可以根据需求自行调整正则表达式
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}
func (*UserAuth) Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		ReturnError(c, g.ErrRequest, err)
		fmt.Println(err)
		return
	}
	db := GetDB(c)
	if len(registerRequest.Password) < 6 || len(registerRequest.Password) > 15 {
		ReturnError(c, g.ErrRequest, "密码需大于六位且小于16位~~")
		return
	}
	if !isValidEmail(registerRequest.Email) {
		ReturnError(c, g.ErrRequest, "邮箱格式错误~~")
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Debug(err.Error())
	}
	err = model.CreateUser(db, registerRequest.Email, registerRequest.Username, string(passwordHash))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, registerRequest)
	return
}
func (*UserAuth) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)
	rdb := GetRDB(c)
	userAuth, err := model.GetUserAuthInfoByName(db, loginRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrRequest, nil)
			return
		}
	}
	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(loginRequest.Password))
	if err != nil {
		ReturnError(c, g.ErrRequest, "密码错误")
		return
	}
	userInfo, err := model.GetUserInfoByUserInfoId(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, "不能通过UserInfoId在UserInfo表中找到")
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	// 删除 Redis 中的离线状态
	offlineKey := g.OFFLINE_USER + strconv.Itoa(int(userAuth.ID))
	rdb.Del(rctx, offlineKey).Result()

	token, err := jwt.GenToken(userInfo.ID, roleIds)
	fmt.Println(roleIds)
	ReturnSuccess(c, LoginVO{
		UserInfo: *userInfo,
		Token:    token,
	})
}
func (*UserAuth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)
}
