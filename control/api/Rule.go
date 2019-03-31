package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func RuleCreate(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	Rid, err := model.RuleCreate(dbs.H{
		"CateId":        m.CateId,
		"Name":          m.Name,
		"Brief":         m.Brief,
		"ListTable":     m.ListTable,
		"ListUrl":       m.ListUrl,
		"ListPageStart": m.ListPageStart,
		"ListPageEnd":   m.ListPageEnd,
		"ListPageSize":  m.ListPageSize,
		"ListRule":      m.ListRule,
		"ContentUrl":    m.ContentUrl,
	})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "添加成功", gin.H{"Rid": Rid})
}

func RuleUpdate(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	_, err = model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	err = model.RuleUpdate(dbs.H{
		"CateId":        m.CateId,
		"Name":          m.Name,
		"Brief":         m.Brief,
		"ListTable":     m.ListTable,
		"ListUrl":       m.ListUrl,
		"ListPageStart": m.ListPageStart,
		"ListPageEnd":   m.ListPageEnd,
		"ListPageSize":  m.ListPageSize,
		"ListRule":      m.ListRule,
		"ContentUrl":    m.ContentUrl,
	}, m.Rid)
	if err != nil {
		c.Message("-1", "更新数据库失败："+err.Error())
		return
	}

	c.Message("0", "修改成功")
}

func RuleRead(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	row, err := model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", row)
}

func RuleDelete(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	err = model.RuleDelete(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "删除成功")
}

func RuleList(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}
	if m.CateId > 0 {
		h["CateId"] = m.CateId
	}
	if m.Name != "" {
		h["Name LIKE"] = m.Name
	}

	total, err := model.RuleCount(h)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	list, err := model.RuleList(h, 0, 0)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}

func RuleTest(c *gin.Context) {
	m := struct {
		Url  string
		Rule string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	res, err := http.Get(m.Url)
	if err != nil {
		c.Message("-1", "抓取页面失败: "+err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		c.Message("-1", fmt.Sprintf("抓取页面页面状态出错: %d %s", res.StatusCode, res.Status))
		return
	}

	// 如果没有匹配规则，则直接返回HTML
	if m.Rule == "" {
		b, _ := ioutil.ReadAll(res.Body)
		c.Message("0", "success", gin.H{"html": string(b)})
		return
	}

	// 加载HTML代码
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		c.Message("-1", "不是正确的HTML页面: "+err.Error())
		return
	}

	// 匹配代码
	html, err := doc.Find(m.Rule).Html()
	if err != nil {
		c.Message("-1", "匹配代码失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"html": html})
}
