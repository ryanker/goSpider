package dbs

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

type DB struct {
	*sql.DB
}

type H map[string]interface{}

func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, err
}

func (db *DB) Insert(table string, data H) (id int64, err error) {
	kStr, vStr, args := GetSqlInsert(data)
	stmt, err := db.Prepare("INSERT INTO `" + table + "`(" + kStr + ") VALUES (" + vStr + ")")
	if err != nil {
		return
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return
	}

	id, err = res.LastInsertId()
	return
}

func (db *DB) Delete(table string, where H) (n int64, err error) {
	whereStr, args := GetSqlWhere(where)
	stmt, err := db.Prepare("DELETE FROM `" + table + "`" + whereStr)
	if err != nil {
		return
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return
	}

	n, err = res.RowsAffected()
	return
}

func (db *DB) Update(table string, data H, where H) (n int64, err error) {
	setStr, args := GetSqlUpdate(data)
	whereStr, args2 := GetSqlWhere(where)
	stmt, err := db.Prepare("UPDATE `" + table + "` SET " + setStr + whereStr)
	if err != nil {
		return
	}

	args = append(args, args2...)
	res, err := stmt.Exec(args...)
	if err != nil {
		return
	}

	n, err = res.RowsAffected()
	return
}

func (db *DB) Read(table string, fields string, scanArr []interface{}, where H) (err error) {
	whereStr, args := GetSqlWhere(where)
	stmt, err := db.Prepare("SELECT " + fields + " FROM `" + table + "`" + whereStr)
	if err != nil {
		return
	}

	err = stmt.QueryRow(args...).Scan(scanArr...)
	return
}

func (db *DB) Count(table string, where H) (n int64, err error) {
	whereStr, args := GetSqlWhere(where)
	stmt, err := db.Prepare("SELECT COUNT(*) FROM `" + table + "`" + whereStr)
	if err != nil {
		return
	}

	err = stmt.QueryRow(args...).Scan(&n)
	return
}

func (db *DB) Find(table string, fields string, where H, order string, page int64, pageSize int64) (rows *sql.Rows, err error) {
	whereStr, args := GetSqlWhere(where)
	orderStr := ""
	limitStr := ""
	if order != "" {
		orderStr = " ORDER BY " + order
	}
	if page > 0 && pageSize > 0 {
		start := (page - 1) * pageSize
		limitStr = " LIMIT " + strconv.FormatInt(start, 10) + "," + strconv.FormatInt(pageSize, 10)
	} else if page == 0 && pageSize > 0 {
		limitStr = " LIMIT " + strconv.FormatInt(pageSize, 10)
	}
	s := "SELECT " + fields + " FROM `" + table + "`" + whereStr + orderStr + limitStr
	rows, err = db.Query(s, args...)
	if err != nil {
		return
	}
	return
}

func GetSqlRead(dest H) (fStr string, scanArr []interface{}) {
	for k, v := range dest {
		fStr += "`" + k + "`, "
		scanArr = append(scanArr, v)
	}
	fStr = strings.TrimSuffix(fStr, ", ")
	return
}

func GetSqlInsert(data H) (kStr, vStr string, args []interface{}) {
	for k, v := range data {
		kStr += "`" + k + "`, "
		vStr += "?, "
		args = append(args, v)
	}
	kStr = strings.TrimSuffix(kStr, ", ")
	vStr = strings.TrimSuffix(vStr, ", ")
	return
}

func GetSqlUpdate(data H) (setStr string, args []interface{}) {
	for k, v := range data {
		setStr += "`" + k + "`=?, "
		args = append(args, v)
	}
	setStr = strings.TrimSuffix(setStr, ", ")
	return
}

func GetSqlWhere(where H) (whereStr string, args []interface{}) {
	if len(where) == 0 {
		return
	}
	whereStr = " WHERE "
	for k, v := range where {
		whereStr += "`" + k + "`=? AND "
		args = append(args, v)
	}
	whereStr = strings.TrimSuffix(whereStr, " AND ")
	return
}
