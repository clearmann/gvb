package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	g "gvb/internal/global"
	"gvb/internal/model"
)

type Page struct{}

func (*Page) GetList(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)
	cache, err := getPageCache(rdb)
	if err == nil && cache != nil {
		zap.L().Debug("[handle-page-GetList] get page list from cache")
		ReturnSuccess(c, cache)
	}
	switch err {
	case redis.Nil:
		break
	default:
		ReturnError(c, g.ErrRedisOp, err)
	}
	// 从数据库中获取 PageList
	data, _, err := model.GetPageList(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
		return
	}
	// 添加到 cache
	if err = addPageCache(GetRDB(c), data); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}
	ReturnSuccess(c, data)
}
