package model

import (
	"time"

	"../lib/dbs"
)

type Rule struct {
	Rid          int64
	Status       int64
	IntervalHour int64
	Name         string
	Brief        string

	ListTable     string
	ListUrl       string
	ListPageStart int64
	ListPageEnd   int64
	ListPageSize  int64
	ListRule      string

	ContentUrl string

	UpdateDate string
	CreateDate string
	Page       int64
}

func RuleMap() (ptr *Rule, fields string, args *[]interface{}) {
	row := Rule{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Rid":           &row.Rid,
		"Status":        &row.Status,
		"IntervalHour":  &row.IntervalHour,
		"Name":          &row.Name,
		"Brief":         &row.Brief,
		"ListTable":     &row.ListTable,
		"ListUrl":       &row.ListUrl,
		"ListPageStart": &row.ListPageStart,
		"ListPageEnd":   &row.ListPageEnd,
		"ListPageSize":  &row.ListPageSize,
		"ListRule":      &row.ListRule,
		"ContentUrl":    &row.ContentUrl,
		"UpdateDate":    &row.UpdateDate,
		"CreateDate":    &row.CreateDate,
	})
	ptr = &row
	args = &scanArr
	return
}

func RuleCreate(h dbs.H) (Rid int64, err error) {
	h["CreateDate"] = time.Now().Format("2006-01-02 15:04:05")
	Rid, err = db.Insert("Rule", h)
	return
}

func RuleUpdate(h dbs.H, Rid int64) (err error) {
	h["UpdateDate"] = time.Now().Format("2006-01-02 15:04:05")
	_, err = db.Update("Rule", h, dbs.H{
		"Rid": Rid,
	})
	return
}

func RuleRead(Rid int64) (row Rule, err error) {
	data, fields, scanArr := RuleMap()
	err = db.Read("Rule", fields, *scanArr, dbs.H{
		"Rid": Rid,
	})
	row = *data
	return
}

func RuleDelete(Rid int64) (err error) {
	row, err := RuleRead(Rid)
	if err != nil {
		return
	}

	// todo: 删除相关数据
	_, err = db.Delete("RuleParam", dbs.H{"Rid": row.Rid})
	if err != nil {
		return
	}

	_, err = db.Delete("Rule", dbs.H{"Rid": row.Rid})
	return
}

func RuleCount(h dbs.H) (n int64, err error) {
	n, err = db.Count("Rule", h)
	return n, err
}

func RuleList(h dbs.H, page, pageSize int64) (list []Rule, err error) {
	data, fields, scanArr := RuleMap()
	rows, err := db.Find("Rule", fields, h, "Rid DESC", page, pageSize)
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
