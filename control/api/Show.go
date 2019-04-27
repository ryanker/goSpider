package api

import (
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func ShowRead(c *gin.Context) {
	m := struct {
		Rid int64
		Lid int64
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

	data, columns, err := dbc.ReadMap("Content", "*", dbs.H{"Lid": m.Lid})
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}

	c.Message("0", "success", gin.H{"columns": columns, "data": data})
}

func ShowList(c *gin.Context) {
	m := struct {
		Rid         int64
		OrderField  string
		SearchField string
		SearchWord  string
		Page        int64
		PageSize    int64
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

	// ========== 优先读取优先级： Content -> List ==========
	// 内容页表总数量
	totalContent, err := dbc.Count("Content", dbs.H{})
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 列表页表总数量
	totalList, err := dbc.Count("List", dbs.H{})
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}

	// 每页数量
	if m.PageSize == 0 {
		m.PageSize = 20
	}

	// 优先读取内容页表
	if totalContent > 0 {
		ParamList, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}

		var isName bool           // 是否有标题字段
		var isImage bool          // 是否有图片字段
		var isImageDown bool      // 图片字段是否需要下载
		var searchFields []string // 参与搜索的字段
		var orderFields []string  // 参与排序的字段
		order := "Lid"
		for _, v := range ParamList {
			if v.Field == "Name" {
				isName = true
			} else if v.Field == "Image" {
				isImage = true
				if v.DownType == 1 {
					isImageDown = true
				}
			}
			if v.IsSearch == 1 {
				searchFields = append(searchFields, v.Field)
				if m.SearchField == v.Field && m.SearchWord != "" {
					h[m.SearchField+" LIKE"] = "%" + m.SearchWord + "%"
				}
			}
			if v.IsOrder == 1 {
				orderFields = append(orderFields, v.Field)
				if m.OrderField == v.Field {
					order = "`" + v.Field + "`"
				}
			}
		}

		// 总数
		total, err := dbc.Count("Content", h)
		if err != nil {
			c.Message("-1", "获取数量失败: "+err.Error())
			return
		}

		fields := "`Lid`"
		if isName {
			fields += ",`Name`"
		}
		if isImage {
			fields += ",`Image`"
		}
		if order != "Lid" {
			fields += "," + order
		}
		list, _, err := dbc.FindMap("Content", fields, h, order+" DESC", m.Page, m.PageSize)
		if err != nil {
			c.Message("-1", "读取表失败: "+err.Error())
			return
		}

		// "内容页表"没有"图片字段"，尝试去"列表页表"读取是否有"图片字段"
		if !isImage {
			ParamList, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List", "Field": "Image"}, 0, 0)
			if err != nil {
				c.Message("-1", err.Error())
				return
			}
			for _, v := range ParamList {
				if v.Field == "Image" {
					isImage = true
					if v.DownType == 1 {
						isImageDown = true
					}
					continue
				}
			}

			if isImage {
				for k, v := range list {
					Lid, ok := v["Lid"]
					if !ok {
						continue
					}

					// 优先去"下载附件表"读取
					if isImageDown {
						v2, _, err := dbc.ReadMap("ListDownload", "`NewUrl`", dbs.H{"Lid": Lid, "Status": 2, "Field": "Image"})
						if err != nil {
							continue
						}
						NewUrl, ok := v2["NewUrl"]
						if ok {
							NewUrl, ok := NewUrl.(string)
							if ok && NewUrl != "" {
								list[k]["Image"] = NewUrl
								continue
							}
						}
					}

					// 如果没有下载完成，再到"列表页表"读取原始图片路径
					v2, _, err := dbc.ReadMap("List", "`Image`", dbs.H{"Lid": Lid})
					if err != nil {
						continue
					}
					url, ok := v2["Image"]
					if ok {
						url, ok := url.(string)
						if ok && url != "" {
							list[k]["Image"] = url
						}
					}
				}
			}
		} else if isImageDown {
			// 用下载完成的图片替换远程图片
			for k, v := range list {
				Lid, ok := v["Lid"]
				if !ok {
					continue
				}
				v2, _, err := dbc.ReadMap("ContentDownload", "`NewUrl`", dbs.H{"Lid": Lid, "Status": 2, "Field": "Image"})
				if err != nil {
					continue
				}
				NewUrl, ok := v2["NewUrl"]
				if ok {
					NewUrl, ok := NewUrl.(string)
					if ok && NewUrl != "" {
						list[k]["Image"] = NewUrl
					}
				}
			}
		}

		c.Message("0", "success", gin.H{
			"ruleName":     row.Name,
			"totalContent": totalContent,
			"totalList":    totalList,
			"total":        total,
			"list":         list,
			"searchFields": searchFields,
			"orderFields":  orderFields,
			"isName":       isName,
			"isImage":      isImage,
		})
	} else {
		// 如果内容页表没有数据，说明待采集，临时读取列表页数据代替
		ParamList, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}

		var isName bool           // 是否有标题字段
		var isImage bool          // 是否有图片字段
		var isImageDown bool      // 图片字段是否需要下载
		var searchFields []string // 参与搜索的字段
		var orderFields []string  // 参与排序的字段
		order := "Lid"
		for _, v := range ParamList {
			if v.Field == "Name" {
				isName = true
			} else if v.Field == "Image" {
				isImage = true
				if v.DownType == 1 {
					isImageDown = true
				}
			}
			if v.IsSearch == 1 {
				searchFields = append(searchFields, v.Field)
				if m.SearchField == v.Field && m.SearchWord != "" {
					h[m.SearchField+" LIKE"] = "%" + m.SearchWord + "%"
				}
			}
			if v.IsOrder == 1 {
				orderFields = append(orderFields, v.Field)
				if m.OrderField == v.Field {
					order = "`" + v.Field + "`"
				}
			}
		}

		// 总数
		total, err := dbc.Count("List", h)
		if err != nil {
			c.Message("-1", "获取数量失败: "+err.Error())
			return
		}

		fields := "`Lid`"
		if isName {
			fields += ",`Name`"
		}
		if isImage {
			fields += ",`Image`"
		}
		if order != "Lid" {
			fields += "," + order
		}
		list, _, err := dbc.FindMap("List", fields, h, order+" DESC", m.Page, m.PageSize)
		if err != nil {
			c.Message("-1", "读取表失败: "+err.Error())
			return
		}
		if isImageDown {
			// 用下载完成的图片替换远程图片
			for k, v := range list {
				Lid, ok := v["Lid"]
				if !ok {
					continue
				}
				v2, _, err := dbc.ReadMap("ListDownload", "`NewUrl`", dbs.H{"Lid": Lid, "Status": 2, "Field": "Image"})
				if err != nil {
					continue
				}
				NewUrl, ok := v2["NewUrl"]
				if ok {
					NewUrl, ok := NewUrl.(string)
					if ok && NewUrl != "" {
						list[k]["Image"] = NewUrl
					}
				}
			}
		}

		c.Message("0", "success", gin.H{
			"ruleName":     row.Name,
			"totalContent": totalContent,
			"totalList":    totalList,
			"total":        total,
			"list":         list,
			"searchFields": searchFields,
			"orderFields":  orderFields,
			"isName":       isName,
			"isImage":      isImage,
		})
	}
}

func ShowDownload(c *gin.Context) {
	m := struct {
		Rid      int64
		Table    string
		IsName   int64
		Status   int64
		Page     int64
		PageSize int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	// 验证表名称
	if m.Table != "ListDownload" && m.Table != "ContentDownload" {
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
	if m.Status > 0 {
		h["Status"] = m.Status
	}

	// 数量
	total, err := dbc.Count(m.Table, h)
	if err != nil {
		c.Message("-1", "获取数量失败: "+err.Error())
		return
	}

	// 列表
	if m.PageSize == 0 {
		m.PageSize = 20
	}
	list, _, err := dbc.FindMap(m.Table, "*", h, "Id DESC", m.Page, m.PageSize)
	if err != nil {
		c.Message("-1", "读取表失败: "+err.Error())
		return
	}
	if m.IsName == 1 {
		table := "Content"
		if m.Table != "ListDownload" {
			table = "List"
		}

		// 判断是否存在标题字段
		n, err := model.RuleParamCount(dbs.H{"Rid": m.Rid, "Type": table, "Field": "Name"})
		if err != nil {
			c.Message("-1", err.Error())
			return
		}

		// 循环赋值标题
		if n > 0 {
			tmp := make(map[int64]map[string]interface{})
			for k, v := range list {
				var name interface{}
				Lid, ok := v["Lid"]
				if ok {
					Lid, ok := Lid.(int64)
					if ok {
						v2, ok := tmp[Lid]
						if ok {
							name, _ = v2["Name"]
						} else {
							v2, _, err := dbc.ReadMap(table, "`Name`", dbs.H{"Lid": Lid})
							if err == nil {
								name, _ = v2["Name"]
							}
						}
					}
				}
				list[k]["Name"] = name
			}
		}
	}

	c.Message("0", "success", gin.H{"total": total, "list": list, "ruleName": row.Name})
}
