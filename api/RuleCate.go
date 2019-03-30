package api

import (
	"github.com/xiuno/gin"

	"../lib/dbs"
	"../model"
)

func RuleCateCreate(c *gin.Context) {
	m := model.RuleCate{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	_, err = model.RuleCateCreate(dbs.H{
		"CateId":   m.CateId,
		"Name":     m.Name,
		"Brief":    m.Brief,
		"Url":      m.Url,
		"DateBase": m.DateBase,
	})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "添加成功")
}

func RuleCateUpdate(c *gin.Context) {
	m := model.RuleCate{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	_, err = model.RuleCateRead(m.CateId)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	err = model.RuleCateUpdate(dbs.H{
		"CateId":   m.CateId,
		"Name":     m.Name,
		"Brief":    m.Brief,
		"Url":      m.Url,
		"DateBase": m.DateBase,
	}, m.CateId)
	if err != nil {
		c.Message("-1", "更新数据库失败："+err.Error())
		return
	}

	c.Message("0", "修改成功")
}

func RuleCateRead(c *gin.Context) {
	m := model.RuleCate{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	row, err := model.RuleCateRead(m.CateId)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", row)
}

func RuleCateDelete(c *gin.Context) {
	m := model.RuleCate{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	err = model.RuleCateDelete(m.CateId)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "删除成功")
}

func RuleCateList(c *gin.Context) {
	m := model.RuleCate{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}
	if m.Name != "" {
		h["Name LIKE"] = m.Name
	}

	total, err := model.RuleCateCount(h)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	list, err := model.RuleCateList(h, m.Page, 20)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}
