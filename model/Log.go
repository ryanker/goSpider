package model

import (
	"../lib/dbs"
)

type Log struct {
	LogId      int64
	Status     int64
	Runtime    float64
	Message    string
	CreateDate string
}

func LogMap() (ptr *Log, fields string, args *[]interface{}) {
	row := Log{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"LogId":      &row.LogId,
		"Status":     &row.Status,
		"Runtime":    &row.Runtime,
		"Message":    &row.Message,
		"CreateDate": &row.CreateDate,
	})
	ptr = &row
	args = &scanArr
	return
}

func LogCount(h dbs.H) (n int64, err error) {
	n, err = dbLog.Count("Log", h)
	return n, err
}

func LogList(h dbs.H, order string, page, pageSize int64) (list []Log, err error) {
	data, fields, scanArr := LogMap()
	rows, err := dbLog.Find("Log", fields, h, order, page, pageSize)
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
