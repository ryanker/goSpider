package model

import (
	"time"

	"../lib/dbs"
)

type RuleParam struct {
	Pid           int64
	Rid           int64
	CateId        int64
	Type          string
	Field         string
	FieldType     string
	Rule          string
	ValueType     string
	ValueAttr     string
	FilterType    string
	FilterRegexp  string
	FilterRepl    string
	Sort          int64
	IsSearch      int64
	DownType      int64
	DownRule      string
	DownValueType string
	DownValueAttr string
	CreateDate    string
}

func RuleParamMap() (ptr *RuleParam, fields string, args *[]interface{}) {
	row := RuleParam{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Pid":           &row.Pid,
		"Rid":           &row.Rid,
		"CateId":        &row.CateId,
		"Type":          &row.Type,
		"Field":         &row.Field,
		"FieldType":     &row.FieldType,
		"Rule":          &row.Rule,
		"ValueType":     &row.ValueType,
		"ValueAttr":     &row.ValueAttr,
		"FilterType":    &row.FilterType,
		"FilterRegexp":  &row.FilterRegexp,
		"FilterRepl":    &row.FilterRepl,
		"Sort":          &row.Sort,
		"IsSearch":      &row.IsSearch,
		"DownType":      &row.DownType,
		"DownRule":      &row.DownRule,
		"DownValueType": &row.DownValueType,
		"DownValueAttr": &row.DownValueAttr,
		"CreateDate":    &row.CreateDate,
	})
	ptr = &row
	args = &scanArr
	return
}

func RuleParamCreate(h dbs.H) (Pid int64, err error) {
	h["CreateDate"] = time.Now().Format("2006-01-02 15:04:05")
	Pid, err = db.Insert("RuleParam", h)
	return
}

func RuleParamUpdate(h dbs.H, Pid int64) (err error) {
	_, err = db.Update("RuleParam", h, dbs.H{
		"Pid": Pid,
	})
	return
}

func RuleParamRead(Pid int64) (row RuleParam, err error) {
	data, fields, scanArr := RuleParamMap()
	err = db.Read("RuleParam", fields, *scanArr, dbs.H{
		"Pid": Pid,
	})
	row = *data
	return
}

func RuleParamDelete(Pid int64) (err error) {
	row, err := RuleParamRead(Pid)
	if err != nil {
		return
	}

	// todo: 删除相关数据
	_, err = db.Delete("RuleParam", dbs.H{
		"Pid": row.Pid,
	})
	return
}

func RuleParamCount(h dbs.H) (n int64, err error) {
	n, err = db.Count("RuleParam", h)
	return n, err
}

func RuleParamList(h dbs.H, page, pageSize int64) (list []RuleParam, err error) {
	data, fields, scanArr := RuleParamMap()
	rows, err := db.Find("RuleParam", fields, h, "Sort ASC", page, pageSize)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(*scanArr...)
		if err != nil {
			return
		}
		list = append(list, *data)
	}
	return
}
