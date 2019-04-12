package model

import (
	"os"

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
	err = dbLog.Find("Log", fields, *scanArr, h, order, page, pageSize, data, func(row interface{}) {
		if v, ok := row.(*Log); ok {
			list = append(list, *v)
		}
	})
	if err != nil {
		return
	}
	return
}

// 删除数据库
func LogDeleteDB() error {
	dbFile := "./db/log.db"

	// 关闭数据库
	if err := dbLog.Ping(); err == nil {
		err := dbLog.Close()
		if err != nil {
			return err
		}
	}

	// 删除数据库文件
	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		if err = os.Remove(dbFile); err != nil {
			return err
		}
	}

	// 重新打开数据库
	dbLog, err = dbs.Open(dbFile)
	if err != nil {
		return err
	}

	// 重新创建表
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		_, err = dbLog.Exec(logSql)
		if err != nil {
			return err
		}
	}

	if err = dbLog.Ping(); err != nil {
		return err
	}
	return nil
}
