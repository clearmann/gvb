package internal

import (
	"github.com/gin-gonic/gin"
	"gvb/internal/handle"
)

var (
	userAPI     handle.User     //用户
	userAuthAPI handle.UserAuth //用户账号
	blogInfoAPI handle.BlogInfo //博客
	pageAPI     handle.Page     //页面接口
	frontAPI    handle.Front    //前台接口汇总
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
	base.GET("/code", userAuthAPI.SendCode) // 验证码
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	//auth := r.Group("/api")
	//auth.Use(middleware.JWTAuth())
	//auth.Use(middleware.PermissionCheck())
	// 用户模块
	//user := auth.Group("/user")
	//{
	//	user.GET("/list")
	//}
}

// 博客前台的接口：大部分不需要登录，部分需要登录
func registerBlogHandler(r *gin.Engine) {
	base := r.Group("/api/front")
	base.GET("/about", blogInfoAPI.GetAbout) //关于我
	base.GET("/home", frontAPI.GetHomeInfo)  // 博客首页
	base.GET("/page", pageAPI.GetList)       //前台页面

	article := base.Group("/article")
	{
		article.GET("/list", frontAPI.GetArticleList)    // 前台文章列表
		article.GET("/:id", frontAPI.GetArticleInfo)     // 前台文章详情
		article.GET("/archive", frontAPI.GetArchiveList) // 前台文章归档
		article.GET("/search", frontAPI.SearchArticle)   // 前台文章搜索
	}
}
