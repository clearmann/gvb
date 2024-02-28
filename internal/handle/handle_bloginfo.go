package handle

import (
	"github.com/gin-gonic/gin"
	g "gvb/internal/global"
	"gvb/internal/model"
)

type BlogInfo struct{}

type BlogHomeV0 struct {
	ArticleCount int `json:"article_count"` // 文章数量
	UserCount    int `json:"user_count"`    // 用户数量
	MessageCount int `json:"message_count"` // 留言数量
	ViewCount    int `json:"view_count"`    // 访问量
	// CategoryCount int64 `json:"category_count"` // 分类数量
	// TagCount      int64 `json:"tag_count"`      // 标签数量
	// BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // 博客信息
	// PageList      []Page                 `json:"pageList"`
}
type AboutReq struct {
	Content string `json:"content"`
}

func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetConfigByKey(GetDB(c), g.CONFIG_ABOUT))
}
