package model

import (
	"os"

	"../lib/dbs"
	"../lib/misc"
)

var err error
var db *dbs.DB

func init() {
	dbs.LogFile = "./log/db.log"
	dbs.ErrorLogFile = "./log/db.error.log"

	dbFile := "./db/data.db"
	db, err = dbs.Open(dbFile)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	// 文件不存在，则创建表
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		s, _ := misc.FileOpen("./install/install.sql")
		_, err = db.Exec(string(s))
		if err != nil {
			panic(err)
		}
	}
}
