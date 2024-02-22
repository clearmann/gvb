package internal

import (
	"github.com/gin-gonic/gin"
	"gvb/internal/handle"
)

var (
	userApi     handle.User
	userAuthAPI handle.UserAuth
)

func RegisterHandlers(r *gin.Engine) {
	registerBaseHandlers(r)
	registerAdminHandler(r)
	registerBlogHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandlers(r *gin.Engine) {
	base := r.Group("/api")

	base.POST("/login", userAuthAPI.Login)       // 登录
	base.POST("/register", userAuthAPI.Register) // 注册
	base.GET("/logout", userAuthAPI.Logout)      // 退出登录
	//base.POST("/report", blogInfoAPI.Report)        // 上报信息
	//base.GET("/config", blogInfoAPI.GetConfigMap)   // 获取配置
	//base.PATCH("/config", blogInfoAPI.UpdateConfig) // 更新配置
	//base.GET("/code", userAuthAPI.SendCode)         // 验证码
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	//auth := r.Group("/api")
	//auth.GET("/home", blogInfo.GetHomeInfo)
}
func registerBlogHandler(r *gin.Engine) {

}
