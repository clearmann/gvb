package middleware

import "github.com/gin-gonic/gin"

// 基于 JWT 的授权
// 如果存在 session 则直接从 session 中获取用户信息
// 如果不存在 session ，则从 Authorization 中获取 token ，并解析 token 获取用户信息，并设置到 session 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
