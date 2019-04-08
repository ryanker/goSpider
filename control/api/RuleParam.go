package api

import (
	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

func RuleParamCreate(c *gin.Context) {
	m := model.RuleParam{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	_, err = model.RuleParamCreate(dbs.H{
		"Rid":           m.Rid,
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

	c.Message("0", "添加成功")
}

func RuleParamUpdate(c *gin.Context) {
	m := model.RuleParam{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	_, err = model.RuleParamRead(m.Pid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	err = model.RuleParamUpdate(dbs.H{
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
	}, m.Pid)
	if err != nil {
		c.Message("-1", "更新数据库失败："+err.Error())
		return
	}

	c.Message("0", "修改成功")
}

func RuleParamRead(c *gin.Context) {
	m := model.RuleParam{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	row, err := model.RuleParamRead(m.Pid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", row)
}

func RuleParamDelete(c *gin.Context) {
	m := model.RuleParam{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	err = model.RuleParamDelete(m.Pid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "删除成功")
}

func RuleParamList(c *gin.Context) {
	m := model.RuleParam{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}
	if m.Rid > 0 {
		h["Rid"] = m.Rid
	}
	if m.Type != "" {
		h["Type"] = m.Type
	}

	total, err := model.RuleParamCount(h)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	list, err := model.RuleParamList(h, 0, 0)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}
