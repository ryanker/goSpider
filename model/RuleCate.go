package model

import (
	"time"

	"../lib/dbs"
)

type RuleCate struct {
	CateId     int64
	Name       string
	Brief      string
	Url        string
	DateBase   string
	CreateDate string
	Page       int64
}

func RuleCateMap() (ptr *RuleCate, fields string, args *[]interface{}) {
	row := RuleCate{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"CateId":     &row.CateId,
		"Name":       &row.Name,
		"Brief":      &row.Brief,
		"Url":        &row.Url,
		"DateBase":   &row.DateBase,
		"CreateDate": &row.CreateDate,
	})
	ptr = &row
	args = &scanArr
	return
}

func RuleCateCreate(h dbs.H) (CateId int64, err error) {
	h["CreateDate"] = time.Now().Format("2006-01-02 15:04:05")
	CateId, err = db.Insert("RuleCate", h)
	return
}

func RuleCateUpdate(h dbs.H, CateId int64) (err error) {
	_, err = db.Update("RuleCate", h, dbs.H{
		"CateId": CateId,
	})
	return
}

func RuleCateRead(CateId int64) (row RuleCate, err error) {
	data, fields, scanArr := RuleCateMap()
	err = db.Read("RuleCate", fields, *scanArr, dbs.H{
		"CateId": CateId,
	})
	row = *data
	return
}

func RuleCateDelete(CateId int64) (err error) {
	row, err := RuleCateRead(CateId)
	if err != nil {
		return
	}

	// todo: 删除相关数据
	_, err = db.Delete("RuleCate", dbs.H{
		"CateId": row.CateId,
	})
	return
}

func RuleCateCount(h dbs.H) (n int64, err error) {
	n, err = db.Count("RuleCate", h)
	return n, err
}

func RuleCateList(h dbs.H, page, pageSize int64) (list []RuleCate, err error) {
	data, fields, scanArr := RuleCateMap()
	rows, err := db.Find("RuleCate", fields, h, "CateId DESC", page, pageSize)
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
