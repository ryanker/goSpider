package model

import (
	"time"

	"../lib/dbs"
)

type Rule struct {
	Rid          int64
	Status       int64
	IntervalHour int64

	Name  string
	Brief string

	DateBase string
	Cookie   string
	Charset  string
	Timeout  int64

	ListSpecialUrl string
	ListUrl        string
	ListPageStart  int64
	ListPageEnd    int64
	ListPageSize   int64
	ListRule       string

	ContentUrl string

	IsList           int64
	IsListDownAna    int64
	IsListDownRun    int64
	IsContent        int64
	IsContentDownAna int64
	IsContentDownRun int64

	LastStartDate time.Time `time_format:"sql_datetime" time_location:"UTC"`
	LastEndData   time.Time `time_format:"sql_datetime" time_location:"UTC"`
	NextStartDate time.Time `time_format:"sql_datetime" time_location:"UTC"`

	UpdateDate time.Time `time_format:"sql_datetime" time_location:"UTC"`
	CreateDate time.Time `time_format:"sql_datetime" time_location:"UTC"`
}

func RuleMap() (ptr *Rule, fields string, args *[]interface{}) {
	row := Rule{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Rid":              &row.Rid,
		"Status":           &row.Status,
		"IntervalHour":     &row.IntervalHour,
		"Name":             &row.Name,
		"Brief":            &row.Brief,
		"DateBase":         &row.DateBase,
		"Cookie":           &row.Cookie,
		"Charset":          &row.Charset,
		"Timeout":          &row.Timeout,
		"ListSpecialUrl":   &row.ListSpecialUrl,
		"ListUrl":          &row.ListUrl,
		"ListPageStart":    &row.ListPageStart,
		"ListPageEnd":      &row.ListPageEnd,
		"ListPageSize":     &row.ListPageSize,
		"ListRule":         &row.ListRule,
		"ContentUrl":       &row.ContentUrl,
		"IsList":           &row.IsList,
		"IsListDownAna":    &row.IsListDownAna,
		"IsListDownRun":    &row.IsListDownRun,
		"IsContent":        &row.IsContent,
		"IsContentDownAna": &row.IsContentDownAna,
		"IsContentDownRun": &row.IsContentDownRun,
		"LastStartDate":    &row.LastStartDate,
		"LastEndData":      &row.LastEndData,
		"NextStartDate":    &row.NextStartDate,
		"UpdateDate":       &row.UpdateDate,
		"CreateDate":       &row.CreateDate,
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
	err = db.Find("Rule", fields, *scanArr, h, "Rid DESC", page, pageSize, func() {
		list = append(list, *data)
	})
	if err != nil {
		return
	}
	return
}
