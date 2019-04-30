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

	// 数量
	size := doc.Find(m.ListRule).Size()

	// 列表
	var htmlList []interface{}
	doc.Find(m.ListRule).Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		htmlList = append(htmlList, html)
	})

	c.Message("0", "success", gin.H{"size": size, "htmlList": htmlList})
}

// 3.匹配列表页规则的单个字段
func HttpGetListRule(c *gin.Context) {
	m := struct {
		Url          string
		ListRule     string
		ParamRule    string
		ValueType    string
		ValueAttr    string
		FilterType   string
		FilterRegexp string
		FilterRepl   string
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

	var dom *goquery.Selection
	if m.ListRule == "" {
		dom = doc.Find(m.ParamRule).Eq(0)
	} else {
		dom = doc.Find(m.ListRule).Find(m.ParamRule).Eq(0)
	}

	value := ""
	if m.ValueType == "Html" {
		value, _ = dom.Html()
	} else if m.ValueType == "Text" {
		value = dom.Text()
	} else if m.ValueType == "Attr" {
		value, _ = dom.Attr(m.ValueAttr)
	} else {
		value, _ = dom.Html()
	}
	value = misc.StrClear(value, m.FilterType, m.FilterRegexp, m.FilterRepl)

	c.Message("0", "success", gin.H{"html": value})
}

// 4.匹配内容页规则的单个字段
func HttpGetContentRule(c *gin.Context) {
	m := struct {
		Url          string
		ParamRule    string
		ValueType    string
		ValueAttr    string
		FilterType   string
		FilterRegexp string
		FilterRepl   string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	bodyByte, _, err := misc.HttpGetRetry(m.Url)
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

	dom := doc.Find(m.ParamRule).Eq(0)
	value := ""
	if m.ValueType == "Html" {
		value, _ = dom.Html()
	} else if m.ValueType == "Text" {
		value = dom.Text()
	} else if m.ValueType == "Attr" {
		value, _ = dom.Attr(m.ValueAttr)
	} else {
		value, _ = dom.Html()
	}
	value = misc.StrClear(value, m.FilterType, m.FilterRegexp, m.FilterRepl)

	c.Message("0", "success", gin.H{"html": value})
}

// 5.匹配列表页规则下载
func HttpGetListRuleDown(c *gin.Context) {
	m := struct {
		Url           string
		ListRule      string
		ParamRule     string
		ValueType     string
		ValueAttr     string
		FilterType    string
		FilterRegexp  string
		FilterRepl    string
		DownRule      string
		DownValueType string
		DownValueAttr string
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

	doc2 := doc.Find(m.ListRule).Find(m.ParamRule).Eq(0)
	value := ""
	if m.ValueType == "Html" {
		value, _ = doc2.Html()
	} else if m.ValueType == "Text" {
		value = doc2.Text()
	} else if m.ValueType == "Attr" {
		value, _ = doc2.Attr(m.ValueAttr)
	} else {
		value, _ = doc2.Html()
	}
	value = misc.StrClear(value, m.FilterType, m.FilterRegexp, m.FilterRepl)

	// 解析HTML代码
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(value))
	if err != nil {
		c.Message("-2", "不是正确的HTML页面: "+err.Error())
		return
	}

	// 数量
	size := doc.Find(m.DownRule).Size()

	// 列表
	var htmlList []interface{}
	doc.Find(m.DownRule).Each(func(i int, s *goquery.Selection) {
		html := ""
		if m.DownValueType == "Text" {
			html = s.Text()
		} else if m.DownValueType == "Attr" {
			html, _ = s.Attr(m.DownValueAttr)
		}
		htmlList = append(htmlList, html)
	})

	c.Message("0", "success", gin.H{"size": size, "htmlList": htmlList})
}

// 6.匹配内容页规则下载
func HttpGetContentRuleDown(c *gin.Context) {
	m := struct {
		Url           string
		ParamRule     string
		ValueType     string
		ValueAttr     string
		FilterType    string
		FilterRegexp  string
		FilterRepl    string
		DownRule      string
		DownValueType string
		DownValueAttr string
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

	doc2 := doc.Find(m.ParamRule).Eq(0)
	value := ""
	if m.ValueType == "Html" {
		value, _ = doc2.Html()
	} else if m.ValueType == "Text" {
		value = doc2.Text()
	} else if m.ValueType == "Attr" {
		value, _ = doc2.Attr(m.ValueAttr)
	} else {
		value, _ = doc2.Html()
	}
	value = misc.StrClear(value, m.FilterType, m.FilterRegexp, m.FilterRepl)

	// 解析HTML代码
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(value))
	if err != nil {
		c.Message("-2", "不是正确的HTML页面: "+err.Error())
		return
	}

	// 数量
	size := doc.Find(m.DownRule).Size()

	// 列表
	var htmlList []interface{}
	doc.Find(m.DownRule).Each(func(i int, s *goquery.Selection) {
		html := ""
		if m.DownValueType == "Text" {
			html = s.Text()
		} else if m.DownValueType == "Attr" {
			html, _ = s.Attr(m.DownValueAttr)
		}
		htmlList = append(htmlList, html)
	})

	c.Message("0", "success", gin.H{"size": size, "htmlList": htmlList})
}
