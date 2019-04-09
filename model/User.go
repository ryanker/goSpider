package model

import (
	"time"

	"../lib/dbs"
)

type User struct {
	Uid        int64
	Gid        int64
	Name       string
	Email      string
	Mobile     string
	Password   string
	Salt       string
	LoginNum   int64
	LastIP     string
	LastDate   string
	CreateIP   string
	CreateDate string
	UpdateDate string
}

func UserMap() (ptr *User, fields string, args *[]interface{}) {
	row := User{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Uid":        &row.Uid,
		"Gid":        &row.Gid,
		"Name":       &row.Name,
		"Email":      &row.Email,
		"Mobile":     &row.Mobile,
		"Password":   &row.Password,
		"Salt":       &row.Salt,
		"LoginNum":   &row.LoginNum,
		"LastIP":     &row.LastIP,
		"LastDate":   &row.LastDate,
		"CreateIP":   &row.CreateIP,
		"UpdateDate": &row.UpdateDate,
		"CreateDate": &row.CreateDate,
	})
	ptr = &row
	args = &scanArr
	return
}

func UserCreate(h dbs.H) (Uid int64, err error) {
	h["CreateDate"] = time.Now().Format("2006-01-02 15:04:05")
	Uid, err = db.Insert("User", h)
	return
}

func UserUpdate(h dbs.H, Uid int64) (err error) {
	h["UpdateDate"] = time.Now().Format("2006-01-02 15:04:05")
	_, err = db.Update("User", h, dbs.H{
		"Uid": Uid,
	})
	return
}

func UserRead(Uid int64) (row User, err error) {
	data, fields, scanArr := UserMap()
	err = db.Read("User", fields, *scanArr, dbs.H{
		"Uid": Uid,
	})
	row = *data
	return
}

func UserDelete(Uid int64) (err error) {
	row, err := UserRead(Uid)
	if err != nil {
		return
	}

	_, err = db.Delete("User", dbs.H{"Uid": row.Uid})
	return
}

func UserCount(h dbs.H) (n int64, err error) {
	n, err = db.Count("User", h)
	return n, err
}

func UserList(h dbs.H, page, pageSize int64) (list []User, err error) {
	data, fields, scanArr := UserMap()
	rows, err := db.Find("User", fields, h, "Uid DESC", page, pageSize)
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

	// 格式化
	for k := range list {
		list[k].UserFormatSafe()
	}
	return
}

// 干掉敏感信息
func (u *User) UserFormatSafe() {
	u.Password = ""
	u.Salt = ""
}
