package api

import (
	"strconv"

	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func ShowList(c *gin.Context) {
	m := struct {
		Rid      int64
		Status   int64
		Keyword1 string
		Keyword2 string
		Keyword3 string
		Keyword4 string
		Keyword5 string
		Keyword6 string
		Keyword7 string
		Keyword8 string
		Keyword9 string
		Page     int64
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
	dbFile := "./db/" + row.Database + ".db"
	dbc, err := dbs.Open(dbFile)
	if err != nil {
		c.Message("-1", "打开数据库失败: "+err.Error())
		return
	}

	h := dbs.H{}
	if m.Status > 0 {
		h["Status"] = m.Status
	}

	// 判断头图是否勾选下载
	isDown, err := model.RuleParamCount(dbs.H{"Rid": row.Rid, "Type": "List", "Field": "Image", "DownType": 1})
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}

	// 需要搜索的字段
	ParamList, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List", "IsSearch": 1}, 0, 0)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	keywords := dbs.H{}
	for k, v := range ParamList {
		keywords[v.Field] = "Keyword" + strconv.Itoa(k+1)
		f := v.Field + " LIKE"
		if k == 0 {
			h[f] = "%" + m.Keyword1 + "%"
		} else if k == 1 {
			h[f] = "%" + m.Keyword2 + "%"
		} else if k == 2 {
			h[f] = "%" + m.Keyword3 + "%"
		} else if k == 3 {
			h[f] = "%" + m.Keyword4 + "%"
		} else if k == 4 {
			h[f] = "%" + m.Keyword5 + "%"
		} else if k == 5 {
			h[f] = "%" + m.Keyword6 + "%"
		} else if k == 6 {
			h[f] = "%" + m.Keyword7 + "%"
		} else if k == 7 {
			h[f] = "%" + m.Keyword8 + "%"
		} else if k == 8 {
			h[f] = "%" + m.Keyword9 + "%"
		}
	}

	// 数量
	total, err := dbc.Count("List", h)
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 列表
	list, columns, err := dbc.FindMap("List", "*", h, "Lid DESC", m.Page, 20)
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "columns": columns, "list": list, "isDown": isDown, "keywords": keywords})
}
