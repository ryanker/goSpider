package api

import (
	"encoding/base64"
	"encoding/json"
	"strings"

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

	if m.Name == "" {
		c.Message("-1", "采集规则名称不能为空")
		return
	}
	if m.Database == "" {
		c.Message("-1", "数据库名称不能为空")
		return
	}
	if m.Database == "data" || m.Database == "log" {
		c.Message("-1", "数据库名称已经被使用")
		return
	}
	n, err := model.RuleCount(dbs.H{"Database": m.Database})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	if n > 0 {
		c.Message("-1", "数据库名称已经被使用")
		return
	}

	Rid, err := model.RuleCreate(dbs.H{
		"Status":           1,
		"IntervalHour":     1,
		"Name":             m.Name,
		"Brief":            m.Brief,
		"Database":         m.Database,
		"Cookie":           m.Cookie,
		"Charset":          m.Charset,
		"Timeout":          20,
		"ListSpecialUrl":   m.ListSpecialUrl,
		"ListUrl":          m.ListUrl,
		"ListPageStart":    m.ListPageStart,
		"ListPageEnd":      m.ListPageEnd,
		"ListPageSize":     m.ListPageSize,
		"ListRule":         m.ListRule,
		"ContentUrl":       m.ContentUrl,
		"IsList":           1,
		"IsListDownAna":    1,
		"IsListDownRun":    1,
		"IsContent":        1,
		"IsContentDownAna": 1,
		"IsContentDownRun": 1,
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

	if m.Name == "" {
		c.Message("-1", "采集规则名称不能为空")
		return
	}
	if m.Database == "" {
		c.Message("-1", "数据库名称不能为空")
		return
	}
	if m.Database == "data" || m.Database == "log" {
		c.Message("-1", "数据库名称已经被使用")
		return
	}
	row, err := model.RuleReadByDatabase(m.Database)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	if row.Rid != m.Rid {
		c.Message("-1", "数据库名称已经被使用")
		return
	}

	_, err = model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	err = model.RuleUpdate(dbs.H{
		"Name":             m.Name,
		"Brief":            m.Brief,
		"Database":         m.Database,
		"Cookie":           m.Cookie,
		"Charset":          m.Charset,
		"ListSpecialUrl":   m.ListSpecialUrl,
		"ListUrl":          m.ListUrl,
		"ListPageStart":    m.ListPageStart,
		"ListPageEnd":      m.ListPageEnd,
		"ListPageSize":     m.ListPageSize,
		"ListRule":         m.ListRule,
		"ContentUrl":       m.ContentUrl,
		"IsList":           m.IsList,
		"IsListDownAna":    m.IsListDownAna,
		"IsListDownRun":    m.IsListDownRun,
		"IsContent":        m.IsContent,
		"IsContentDownAna": m.IsContentDownAna,
		"IsContentDownRun": m.IsContentDownRun,
	}, m.Rid)
	if err != nil {
		c.Message("-1", "更新数据库失败："+err.Error())
		return
	}

	c.Message("0", "修改成功")
}

func RuleUpdateCron(c *gin.Context) {
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
		"Status":           m.Status,
		"IntervalHour":     m.IntervalHour,
		"Timeout":          m.Timeout,
		"IsList":           m.IsList,
		"IsListDownAna":    m.IsListDownAna,
		"IsListDownRun":    m.IsListDownRun,
		"IsContent":        m.IsContent,
		"IsContentDownAna": m.IsContentDownAna,
		"IsContentDownRun": m.IsContentDownRun,
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
	m := struct {
		model.Rule
		Page int64
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

func RuleExport(c *gin.Context) {
	m := model.Rule{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 读取规则
	row, err := model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	// 读取规则参数
	ParamList, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	ParamContent, err := model.RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	RuleSt := struct {
		Rule         model.Rule
		ParamList    []model.RuleParam
		ParamContent []model.RuleParam
	}{row, ParamList, ParamContent}
	b, err := json.Marshal(RuleSt)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	data := "[GO:BEGIN]\n" + base64.StdEncoding.EncodeToString(b) + "\n[GO:END]"
	c.Message("0", "success", data)
}

func RuleImport(c *gin.Context) {
	m := struct {
		Data string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	m.Data = strings.TrimPrefix(m.Data, "[GO:BEGIN]\n")
	m.Data = strings.TrimSuffix(m.Data, "\n[GO:END]")
	m.Data = strings.Trim(m.Data, " \t\r\n")
	b, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		c.Message("-1", "解析 encoded 失败："+err.Error())
		return
	}

	F := struct {
		Rule         model.Rule
		ParamList    []model.RuleParam
		ParamContent []model.RuleParam
	}{}
	err = json.Unmarshal(b, &F)
	if err != nil {
		c.Message("-1", "解析 Json 失败："+err.Error())
		return
	}

	Rid, err := model.RuleCreate(dbs.H{
		"Status":           1,
		"IntervalHour":     F.Rule.IntervalHour,
		"Name":             F.Rule.Name,
		"Brief":            F.Rule.Brief,
		"Database":         F.Rule.Database,
		"Cookie":           F.Rule.Cookie,
		"Charset":          F.Rule.Charset,
		"Timeout":          F.Rule.Timeout,
		"ListSpecialUrl":   F.Rule.ListSpecialUrl,
		"ListUrl":          F.Rule.ListUrl,
		"ListPageStart":    F.Rule.ListPageStart,
		"ListPageEnd":      F.Rule.ListPageEnd,
		"ListPageSize":     F.Rule.ListPageSize,
		"ListRule":         F.Rule.ListRule,
		"ContentUrl":       F.Rule.ContentUrl,
		"IsList":           1,
		"IsListDownAna":    1,
		"IsListDownRun":    1,
		"IsContent":        1,
		"IsContentDownAna": 1,
		"IsContentDownRun": 1,
	})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	for _, m := range F.ParamList {
		_, err = model.RuleParamCreate(dbs.H{
			"Rid":           Rid,
			"Type":          m.Type,
			"Field":         m.Field,
			"FieldType":     m.FieldType,
			"Brief":         m.Brief,
			"Rule":          m.Rule,
			"ValueType":     m.ValueType,
			"ValueAttr":     m.ValueAttr,
			"FilterType":    m.FilterType,
			"FilterRegexp":  m.FilterRegexp,
			"FilterRepl":    m.FilterRepl,
			"Sort":          m.Sort,
			"IsSearch":      m.IsSearch,
			"DownType":      m.DownType,
			"DownRule":      m.DownRule,
			"DownValueType": m.DownValueType,
			"DownValueAttr": m.DownValueAttr,
			"DownFileType":  m.DownFileType,
			"DownTimeout":   m.DownTimeout,
		})
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	}

	for _, m := range F.ParamContent {
		_, err = model.RuleParamCreate(dbs.H{
			"Rid":           Rid,
			"Type":          m.Type,
			"Field":         m.Field,
			"FieldType":     m.FieldType,
			"Brief":         m.Brief,
			"Rule":          m.Rule,
			"ValueType":     m.ValueType,
			"ValueAttr":     m.ValueAttr,
			"FilterType":    m.FilterType,
			"FilterRegexp":  m.FilterRegexp,
			"FilterRepl":    m.FilterRepl,
			"Sort":          m.Sort,
			"IsSearch":      m.IsSearch,
			"DownType":      m.DownType,
			"DownRule":      m.DownRule,
			"DownValueType": m.DownValueType,
			"DownValueAttr": m.DownValueAttr,
			"DownFileType":  m.DownFileType,
			"DownTimeout":   m.DownTimeout,
		})
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	}

	c.Message("0", "导入成功")
}
