package api

import (
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func ItemList(c *gin.Context) {
	m := struct {
		Rid   int64
		Table string
		Page  int64
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

	// 数量
	total, err := dbc.Count(m.Table, dbs.H{})
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 列表
	rows, err := dbc.Find(m.Table, "*", dbs.H{}, "", m.Page, 200)
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}
	var columns []string
	var list []map[string]interface{}
	for rows.Next() {
		m := map[string]interface{}{}
		columns, err = dbs.MapScan(rows, m)
		if err != nil {
			c.Message("-1", "MapScan 失败: "+err.Error())
			return
		}
		list = append(list, m)
	}

	c.Message("0", "success", gin.H{"total": total, "columns": columns, "list": list})
}
