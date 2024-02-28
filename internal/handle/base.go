package handle

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	g "gvb/internal/global"
	"gvb/internal/model"
	"net/http"
)

/*
响应设计方案：不使用 HTTP 码来表示业务状态，采用业务状态码来表示
· 只要能到达后端的业务请求，HTTP 状态码都为 200
· 当后端发生 panic 错误，并且被 gin 中间件捕获时，才会返回  HTTP 500 状态码
-- 业务状态码为 0 表示成功，其他都表示失败
*/
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// Http 码 + 业务码 + 消息 + 数据
func ReturnHttpResponse(c *gin.Context, httpCode int, code int, message string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: message,
		Data:    data,
	})
	return
}

// 业务码 + 数据
func ReturnResponse(c *gin.Context, r g.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Message(), data)
}

// 成功业务码 + 数据
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.OkResult, data)
}

// 所有可预料的错误 = 业务错误 + 系统错误, 在业务层面处理, 返回 HTTP 200 状态码
// 对于不可预料的错误, 会触发 panic, 由 gin 中间件捕获, 并返回 HTTP 500 状态码
// err 是业务错误, data 是错误数据 (可以是 error 或 string)
func ReturnError(c *gin.Context, r g.Result, data any) {
	zap.L().Info(r.Message())
	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Message(),
			Data:    data,
		},
	)
	return
}
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}

// 获取当前登录用户的信息
/*
1. 能从 gin Context 中获取到 user 对象，说明本次请求链路中获取过了
2. 从 session 中获取到 uid
3. 根据 uid 获取用户信息，并设置到 gin Context上
*/
func GetCurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := g.CTX_USER_AUTH
	// 1.
	if cache, exist := c.Get(key); exist && cache != nil {
		zap.L().Debug("[Func-CurrentUserAuth] get from cache: " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth), nil
	}
	// 2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	//3
	db := GetDB(c)
	user, err := model.GetUserAuthById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}

// 分页获取数据
type PageQuery struct {
	Page    int    `form:"page_num"`  // 当前页数（从1开始）
	Size    int    `form:"page_size"` // 每页条数
	Keyword string `form:"keyword"`   // 搜索关键字
}

// 分页响应数据
type PageResult[T any] struct {
	Page  int   `json:"page_num"`  // 每页条数
	Size  int   `json:"page_size"` // 上次页数
	Total int64 `json:"total"`     // 总条数
	List  []T   `json:"page_data"` // 分页数据
}
