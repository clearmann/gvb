package handle

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	g "gvb/internal/global"
	"gvb/internal/model"
)

// redis context
var rctx = context.Background()

// 从 Redis 中获取页面列表缓存
// rdb.Get 如果不存在 key, 会返回 redis.Nil 错误
func getPageCache(rdb *redis.Client) ([]model.Page, error) {
	var cache []model.Page
	s, err := rdb.Get(rctx, g.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}
	return cache, nil
}

// 将页面列表缓存到 redis
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, g.PAGE, string(data), 0).Err()
}
