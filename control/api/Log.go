package api

import (
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func LogList(c *gin.Context) {
	m := struct {
		model.Log
		Order string
		Page  int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}
	if m.Status > 0 {
		h["Status"] = m.Status
	}
	if m.Message != "" {
		h["Message LIKE"] = "%" + m.Message + "%"
	}

	// 排序
	order := ""
	if m.Order == "Runtime" {
		order = "Runtime DESC"
	} else {
		order = "LogId DESC"
	}

	total, err := model.LogCount(h)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	list, err := model.LogList(h, order, m.Page, 20)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}

func LogDeleteDB(c *gin.Context) {
	err := model.LogDeleteDB()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "删除完成")
}
