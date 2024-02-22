package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	g "gvb/internal/global"
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
}
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}
