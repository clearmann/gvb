package internal

import (
	"github.com/gin-gonic/gin"
	"gvb/internal/handle"
	"gvb/internal/middleware"
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
	//不需要登录
	base := r.Group("/api/front")
	base.GET("/about", blogInfoAPI.GetAbout) //关于我
	base.GET("/home", frontAPI.GetHomeInfo)  // 前台首页
	base.GET("/page", pageAPI.GetList)       //前台页面

	article := base.Group("/article")
	{
		article.GET("/list", frontAPI.GetArticleList)    // 前台文章列表
		article.GET("/:id", frontAPI.GetArticleInfo)     // 前台文章详情
		article.GET("/archive", frontAPI.GetArchiveList) // 前台文章归档
		article.GET("/search", frontAPI.SearchArticle)   // 前台文章搜索
	}
	base.GET("/category/list", frontAPI.GetCategoryList) //前台分类列表
	base.GET("/tag/list", frontAPI.GetTagList)           // 前台标签列表
	//base.GET("/link/list", frontAPI.GetLinkList)//前台友链列表
	//base.GET("/message/list", frontAPI.GetMessageList) // 前台留言列表
	base.GET("/comment/list")                //前台评论列表
	base.GET("/comment/replies/:comment_id") // 根据评论 id 查询回复
	// 需要登录
	base.Use(middleware.JWTAuth())
	{
		//base.POST("/upload", uploadAPI.UploadFile)    // 文件上传
		base.GET("/user/info", userAPI.GetInfo)       // 根据 Token 获取用户信息
		base.PUT("/user/info", userAPI.UpdateCurrent) // 根据 Token 更新当前用户信息

		base.POST("/message", frontAPI.SaveMessage)                 // 前台新增留言
		base.POST("/comment", frontAPI.SaveComment)                 // 前台新增评论
		base.GET("/comment/like/:comment_id", frontAPI.LikeComment) // 前台点赞评论
		base.GET("/article/like/:article_id", frontAPI.LikeArticle) // 前台点赞文章
	}
}
