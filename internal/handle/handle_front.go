package handle

import (
	"github.com/gin-gonic/gin"
	g "gvb/internal/global"
	"gvb/internal/model"
	"strconv"
	"strings"
	"time"
)

type Front struct{}

func (*Front) GetHomeInfo(c *gin.Context) {
	//从数据库中获取相关信息
	db := GetDB(c)
	rdb := GetRDB(c)
	data, err := model.GetFrontStatics(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
		return
	}
	data.ViewCount, _ = rdb.Get(rctx, g.VIEW_COUNT).Int64()
	ReturnSuccess(c, data)
}

type ArticleQuery struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}

/*
文章相关接口
*/
// 获取文章列表
func (*Front) GetArticleList(c *gin.Context) {
	db := GetDB(c)
	var query ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
	}
	list, _, err := model.GetArticleList(db, query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
	}
	ReturnSuccess(c, list)
}

// 根据文章 [:id] 获取文章详情
func (*Front) GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, nil)
	}
	// 获取 db 和 rdb
	db := GetDB(c)
	rdb := GetRDB(c)
	//文章详情
	val, err := model.GetBlogArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	article := model.BlogArticleVO{Article: *val}
	//推荐文章 5 篇
	article.RecommendArticles, err = model.GetRecommendList(db, id, 5)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	// 最新文章（5篇）
	article.NewestArticles, err = model.GetNewestList(db, 5)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}

	// 目前请求一次就会增加访问量, 即刷新可以刷访问量
	rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 1, strconv.Itoa(id))
	// 上一篇文章
	article.LastArticle, err = model.GetLastArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	// 下一篇文章
	article.NextArticle, err = model.GetNextArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	// 点赞量, 浏览量
	article.ViewCount = int64(rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, strconv.Itoa(id)).Val())
	article.LikeCount = int64(rdb.ZScore(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())

	// 评论数量
	article.CommentCount, err = model.GetArticleCommentCount(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	ReturnSuccess(c, article)
}

type FArticleQuery struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}
type ArchiveVO struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Created_at time.Time `json:"created_at"`
}

// GetArchiveList 获取文章归档
func (*Front) GetArchiveList(c *gin.Context) {
	var query FArticleQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		ReturnError(c, g.ErrRequest, nil)
	}
	db := GetDB(c)
	list, total, err := model.GetBlogArticleList(db, query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, nil)
	}
	archives := make([]ArchiveVO, 0)
	for _, article := range list {
		archives = append(archives, ArchiveVO{
			ID:         article.ID,
			Title:      article.Title,
			Created_at: article.CreatedAt,
		})
	}
	ReturnSuccess(c, PageResult[ArchiveVO]{
		Page:  query.Page,
		Size:  query.Size,
		Total: total,
		List:  archives,
	})
}

// 文章搜索
type ArticleSearchVO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (*Front) SearchArticle(c *gin.Context) {
	db := GetDB(c)
	result := make([]ArticleSearchVO, 0)
	keyword := c.Query("keyword")
	if keyword == "" {
		ReturnSuccess(c, result)
	}
	articleList, err := model.List(db, []model.Article{}, "*", "",
		"is_delete = 0 and status = 1 and (title like ? or content like ?)",
		"%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}
	for _, article := range articleList {
		// 高亮标题中的关键字
		title := strings.ReplaceAll(article.Title, keyword,
			"<span style='color:#f47466'>"+keyword+"</span>")

		content := article.Content
		// 关键字在内容中的起始位置
		keywordStartIndex := unicodeIndex(content, keyword)
		if keywordStartIndex != -1 { // 关键字在内容中
			preIndex, afterIndex := 0, 0
			if keywordStartIndex > 25 {
				preIndex = keywordStartIndex - 25
			}
			// 防止中文截取出乱码 (中文在 golang 是 3 个字符, 使用 rune 中文占一个数组下标)
			preText := substring(content, preIndex, keywordStartIndex)
			// string([]rune(content[preIndex:keywordStartIndex]))

			// 关键字在内容中的结束位置
			keywordEndIndex := keywordStartIndex + unicodeLen(keyword)
			afterLength := len(content) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			// afterText := string([]rune(content)[keywordStartIndex:afterIndex])
			afterText := substring(content, keywordStartIndex, afterIndex)
			// 高亮内容中的关键字
			content = strings.ReplaceAll(preText+afterText, keyword,
				"<span style='color:#f47466'>"+keyword+"</span>")
		}

		result = append(result, ArticleSearchVO{
			ID:      article.ID,
			Title:   title,
			Content: content,
		})
	}

	ReturnSuccess(c, result)
}

// 获取带中文的字符串中子字符串的实际位置，非字节位置
func unicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

// 解决中文获取位置不正确问题
func substring(source string, start int, end int) string {
	var unicodeStr = []rune(source)
	length := len(unicodeStr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(unicodeStr[i])
	}
	return substring
}

// 获取带中文的字符串实际长度，非字节长度
func unicodeLen(str string) int {
	var r = []rune(str)
	return len(r)
}
