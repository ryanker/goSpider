package api

import (
	"github.com/ryanker/gin_v1"

	"../../lib/dbs"
	"../../model"
)

func TableList(c *gin.Context) {
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

	// 验证表名称
	if m.Table != "List" && m.Table != "Content" && m.Table != "ListDownload" && m.Table != "ContentDownload" {
		c.Message("-1", "Table value is error")
		return
	}

	// 读取规则
	row, err := model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	// 打开数据库
	dbFile := "./db/" + row.Database + ".db"
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

	// 排序
	order := ""
	if m.Table == "List" || m.Table == "Content" {
		order = "Lid DESC"
	} else if m.Table == "ListDownload" || m.Table == "ContentDownload" {
		order = "Id DESC"
	}

	// 列表
	list, columns, err := dbc.FindMap(m.Table, "*", h, order, m.Page, 20)
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "columns": columns, "list": list})
}
