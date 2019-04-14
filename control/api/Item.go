package api

import (
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func ItemList(c *gin.Context) {
	m := struct {
		Rid    int64
		Table  string
		Status int64
		Page   int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	// 读取规则
	row, err := model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	// 打开数据库
	dbFile := "./db/" + row.DateBase + ".db"
	dbc, err := dbs.Open(dbFile)
	if err != nil {
		c.Message("-1", "打开数据库失败: "+err.Error())
		return
	}

	h := dbs.H{}
	if m.Status > 0 && m.Table != "Content" {
		h["Status"] = m.Status
	}

	// 数量
	total, err := dbc.Count(m.Table, h)
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 列表
	list, columns, err := dbc.FindMap(m.Table, h, "", m.Page, 20)
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "columns": columns, "list": list})
}
