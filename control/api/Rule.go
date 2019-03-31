package api

import (
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

	_, err = model.RuleCreate(dbs.H{
		"CateId":        m.CateId,
		"Name":          m.Name,
		"Brief":         m.Brief,
		"ListTable":     m.ListTable,
		"ListUrl":       m.ListUrl,
		"ListPageStart": m.ListPageStart,
		"ListPageEnd":   m.ListPageEnd,
		"ListPageSize":  m.ListPageSize,
		"ListRange":     m.ListRange,
		"ListRule":      m.ListRule,
		"ContentTable":  m.ContentTable,
		"ContentUrl":    m.ContentUrl,
	})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "添加成功")
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
		"ListRange":     m.ListRange,
		"ListRule":      m.ListRule,
		"ContentTable":  m.ContentTable,
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

	list, err := model.RuleList(h, m.Page, 20)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}
