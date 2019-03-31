package model

import (
	"time"

	"../lib/dbs"
)

type RuleParam struct {
	Pid        int64
	Rid        int64
	CateId     int64
	Type       int64
	Field      string
	FieldType  int64
	Rule       string
	ValueType  int64
	ValueAttr  string
	FilterType int64
	FilterReg  string
	Sort       int64
	IsDown     int64
	CreateDate string
}

func RuleParamMap() (ptr *RuleParam, fields string, args *[]interface{}) {
	row := RuleParam{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Pid":        &row.Pid,
		"Rid":        &row.Rid,
		"CateId":     &row.CateId,
		"Type":       &row.Type,
		"Field":      &row.Field,
		"FieldType":  &row.FieldType,
		"Rule":       &row.Rule,
		"ValueType":  &row.ValueType,
		"ValueAttr":  &row.ValueAttr,
		"FilterType": &row.FilterType,
		"FilterReg":  &row.FilterReg,
		"Sort":       &row.Sort,
		"IsDown":     &row.IsDown,
		"CreateDate": &row.CreateDate,
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
	rows, err := db.Find("RuleParam", fields, h, "Pid DESC", page, pageSize)
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
