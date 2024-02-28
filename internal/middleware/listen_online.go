package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	g "gvb/internal/global"
	"gvb/internal/handle"
	"strconv"
	"time"
)

// 监听在线状态
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户
		// 判断当前用户是否在线
		ctx := context.Background()
		rdb := c.MustGet(g.CTX_RDB).(*redis.Client)
		auth, err := handle.GetCurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ErrUserAuth, err)
			return
		}
		onlineKey := g.ONLINE_USER + strconv.Itoa(int(auth.ID))
		offlineKey := g.OFFLINE_USER + strconv.Itoa(int(auth.ID))
		// 判断当前用户是否被强制下线
		if rdb.Exists(ctx, offlineKey).Val() == 1 {
			fmt.Println("用户被强制下线")
			handle.ReturnError(c, g.ErrForceOffline, nil)
			c.Abort()
			return
		}
		// 每次发送请求会更新 Redis 中的在线状态: 重新计算 10 分钟
		rdb.Set(ctx, onlineKey, auth, 10*time.Minute)
		c.Next()
	}
}
