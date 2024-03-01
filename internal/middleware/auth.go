package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	g "gvb/internal/global"
	"gvb/internal/handle"
	"gvb/internal/model"
	"gvb/internal/utils/jwt"
	"strings"
	"time"
)

// 基于 JWT 的授权
// 如果存在 session 则直接从 session 中获取用户信息
// 如果不存在 session ，则从 Authorization 中获取 token ，并解析 token 获取用户信息，并设置到 session 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zap.L().Info("[middleware-JWTAuth] user auth not exist, do jwt auth")
		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		//有的业务需要做身份验证，没有加进来的则不需要
		url, method := c.FullPath()[4:], c.Request.Method
		resource, err := model.GetResource(db, url, method)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//没有找到资源，可直接跳过验证信息
				zap.L().Info("[middleware-JWTAuth] resource not exist, skip jwt auth")
				c.Set("skip_check", true)
				c.Next()
				c.Set("skip_check", false)
				return
			}
			handle.ReturnError(c, g.ErrDbOp, nil)
		}
		// 匿名资源, 直接跳过后续验证
		if resource.Anonymous {
			zap.L().Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s is anonymous, skip jwt auth!", url, method))
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, g.ErrRequest, nil)
		}
		// token 的正确格式为 "Bearer [tokenString]"
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ErrTokenType, nil)
		}
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			handle.ReturnError(c, g.ErrTokenWrong, nil)
		}
		// 判断token 是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ErrTokenRuntime, nil)
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		// session
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserId)
		session.Save()

		// gin context
		c.Set(g.CTX_USER_AUTH, user)
	}
}
