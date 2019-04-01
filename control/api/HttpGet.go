package api

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xiuno/gin"

	"../../lib/misc"
)

// 1.抓取整个网页
func HttpGetPage(c *gin.Context) {
	m := struct {
		Url string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	// 抓取页面
	bodyByte, err := misc.HttpGet(m.Url)
	if err != nil {
		c.Message("-1", "抓取页面失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"html": string(bodyByte)})
}

// 2.匹配列表页规则
func HttpGetList(c *gin.Context) {
	m := struct {
		Url      string
		ListRule string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	bodyByte, err := misc.HttpGet(m.Url)
	if err != nil {
		c.Message("-1", "抓取页面失败: "+err.Error())
		return
	}

	// 解析HTML代码
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		c.Message("-1", "不是正确的HTML页面: "+err.Error())
		return
	}

	var htmlList []interface{}
	doc.Find(m.ListRule).Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		htmlList = append(htmlList, html)
	})

	c.Message("0", "success", gin.H{"htmlList": htmlList})
}

// 3.匹配列表页规则的单个字段
func HttpGetListRule(c *gin.Context) {
	m := struct {
		Url       string
		ListRule  string
		ParamRule string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	bodyByte, err := misc.HttpGet(m.Url)
	if err != nil {
		c.Message("-1", "抓取页面失败: "+err.Error())
		return
	}

	// 解析HTML代码
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		c.Message("-1", "不是正确的HTML页面: "+err.Error())
		return
	}

	html, err := doc.Find(m.ListRule).Find(m.ParamRule).Html()
	if err != nil {
		c.Message("-1", "匹配代码失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"html": html})
}

// 4.匹配内容页规则的单个字段
func HttpGetContentRule(c *gin.Context) {
	m := struct {
		Url       string
		ParamRule string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	bodyByte, err := misc.HttpGet(m.Url)
	if err != nil {
		c.Message("-1", "抓取页面失败: "+err.Error())
		return
	}

	// 解析HTML代码
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		c.Message("-1", "不是正确的HTML页面: "+err.Error())
		return
	}

	html, err := doc.Find(m.ParamRule).Html()
	if err != nil {
		c.Message("-1", "匹配代码失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"html": html})
}
