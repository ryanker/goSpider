package main

import (
	"../../dbs"
	"fmt"
	"os"
	"time"
)

type User struct {
	Uid        int64
	Gid        int64
	Name       string
	CreateDate string
}

var err error
var db *dbs.DB

func main() {
	os.Remove("./test.db")

	// 开启日志
	dbs.LogFile = "./db.log"
	dbs.ErrorLogFile = "./db.error.log"

	db, err = dbs.Open("./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 创建表
	_, err = db.Exec(`DROP TABLE IF EXISTS user;
CREATE TABLE user
(
  uid        INTEGER PRIMARY KEY AUTOINCREMENT,
  gid        INTEGER NOT NULL DEFAULT '0',
  name       TEXT             DEFAULT '',
  createDate DATETIME         DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		panic(err)
	}

	// 插入
	uid, err := db.Insert("user", dbs.H{
		"gid":        1,
		"name":       "admin",
		"createDate": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert: ", uid)

	uid, err = db.Insert("user", dbs.H{
		"gid":        1,
		"name":       "admin2",
		"createDate": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert: ", uid)

	// 更新
	n, err := db.Update("user", dbs.H{
		"gid":  "2",
		"name": "test",
	}, dbs.H{
		"uid": uid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Update: ", n)

	// 统计数量
	n, err = db.Count("user", dbs.H{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Count: ", n)

	// ===========================================================================
	// 映射结构体
	scanF := func() (ptr *User, fields string, args *[]interface{}) {
		row := User{}
		fields, scanArr := dbs.GetSqlRead(dbs.H{
			"uid":        &row.Uid,
			"gid":        &row.Gid,
			"name":       &row.Name,
			"createDate": &row.CreateDate,
		})
		ptr = &row
		args = &scanArr
		return
	}
	data, fields, scanArr := scanF()

	// 读取(到结构体)
	err = db.Read("user", fields, *scanArr, dbs.H{
		"uid": uid,
	})
	if err != nil {
		panic(err)
	}
	u := *data
	fmt.Printf("Read: %+v\n", u)

	// 读取多条(到结构体)
	rows, err := db.Find("user", fields, dbs.H{}, "", 1, 20)
	if err != nil {
		panic(err)
	}
	var list []User
	for rows.Next() {
		err = rows.Scan(*scanArr...)
		if err != nil {
			panic(err)
		}
		list = append(list, *data)
	}
	fmt.Println("List: ", list)

	// 删除
	n, err = db.Delete("user", dbs.H{
		"uid": uid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete: ", n)
}
