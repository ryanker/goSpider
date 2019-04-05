package dbs

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	LogInit()
	return &DB{DB: db}, err
}

func (db *DB) Insert(table string, data H) (id int64, err error) {
	kStr, vStr, args := GetSqlInsert(data)
	s := "INSERT INTO `" + table + "`(" + kStr + ") VALUES (" + vStr + ")"
	LogWrite(s, args...)

	stmt, err := db.Prepare(s)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	id, err = res.LastInsertId()
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}
	return
}

func (db *DB) Delete(table string, where H) (n int64, err error) {
	whereStr, args := GetSqlWhere(where)
	s := "DELETE FROM `" + table + "`" + whereStr
	LogWrite(s, args...)

	stmt, err := db.Prepare(s)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	n, err = res.RowsAffected()
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}
	return
}

func (db *DB) Update(table string, data H, where H) (n int64, err error) {
	setStr, args := GetSqlUpdate(data)
	whereStr, args2 := GetSqlWhere(where)
	args = append(args, args2...)

	s := "UPDATE `" + table + "` SET " + setStr + whereStr
	LogWrite(s, args...)

	stmt, err := db.Prepare(s)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	n, err = res.RowsAffected()
	return
}

func (db *DB) Read(table string, fields string, scanArr []interface{}, where H) (err error) {
	whereStr, args := GetSqlWhere(where)
	s := "SELECT " + fields + " FROM `" + table + "`" + whereStr
	LogWrite(s, args...)

	stmt, err := db.Prepare(s)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	err = stmt.QueryRow(args...).Scan(scanArr...)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}
	return
}

func (db *DB) Count(table string, where H) (n int64, err error) {
	whereStr, args := GetSqlWhere(where)
	s := "SELECT COUNT(*) FROM `" + table + "`" + whereStr
	LogWrite(s, args...)

	stmt, err := db.Prepare(s)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}

	err = stmt.QueryRow(args...).Scan(&n)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}
	return
}

func (db *DB) Find(table, fields string, where H, order string, page, pageSize int64) (rows *sql.Rows, err error) {
	whereStr, args := GetSqlWhere(where)
	orderStr := ""
	limitStr := ""
	if order != "" {
		orderStr = " ORDER BY " + order
	}
	if page < 1 {
		page = 1
	}
	if pageSize > 0 {
		start := (page - 1) * pageSize
		limitStr = " LIMIT " + strconv.FormatInt(start, 10) + "," + strconv.FormatInt(pageSize, 10)
	}
	s := "SELECT " + fields + " FROM `" + table + "`" + whereStr + orderStr + limitStr
	LogWrite(s, args...)

	rows, err = db.Query(s, args...)
	if err != nil {
		ErrorLogWrite(err, s, args...)
		return
	}
	return
}

// 绑定到 Map
func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
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
		symbol := k[:1]
		if symbol == "+" || symbol == "-" {
			key := k[1:]
			setStr += "`" + key + "`=`" + key + "`" + symbol + "?, "
		} else {
			setStr += "`" + k + "`=?, "
		}
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
		i := strings.Index(k, " ")
		if i == -1 {
			whereStr += "`" + k + "`=? AND "
		} else {
			key := k[:i]
			symbol := strings.Trim(k[i:], " ")
			whereStr += "`" + key + "` " + symbol + " ? AND "
		}
		args = append(args, v)
	}
	whereStr = strings.TrimSuffix(whereStr, " AND ")
	return
}

var LogFile string
var ErrorLogFile string

var LogIoWriter io.Writer = os.Stdout

func LogInit() {
	if LogFile != "" {
		var err error
		LogIoWriter, err = os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func LogWrite(s string, args ...interface{}) {
	if LogFile == "" {
		return
	}
	_, err := fmt.Fprintf(LogIoWriter, "%v | %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(strings.Replace(s, "?", "%v", -1), args...),
	)
	if err != nil {
		fmt.Println(err)
	}
}

func ErrorLogWrite(e error, s string, args ...interface{}) {
	if ErrorLogFile == "" {
		return
	}

	f, err := os.OpenFile(ErrorLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	str := fmt.Sprintf("%v | ERROR: %v | SQL: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		e.Error(),
		fmt.Sprintf(strings.Replace(s, "?", "%v", -1), args...),
	)
	if _, err := f.Write([]byte(str)); err != nil {
		fmt.Println(err)
	}

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}
