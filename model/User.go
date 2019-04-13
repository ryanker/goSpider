package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"../lib/dbs"
	"../lib/misc"
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

type UserToken struct {
	Uid      int64  // 用户ID
	Password string // 用户密码MD5
	LastIP   string // 最后登录IP
	LastDate string // 最后登录时间
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

func UserTokenEncode(Uid int64, Password string, LastIP string) string {
	s := fmt.Sprintf("%v,%v,%v,%v", Uid, Password, LastIP, time.Now().Format("2006-01-02 15:04:05"))
	return misc.Base16Encode(s)
}

func UserTokenDecode(s string) (token UserToken, err error) {
	if s == "" {
		return token, errors.New("token is empty")
	}
	s, err = misc.Base16Decode(s)
	if err != nil {
		return
	}
	arr := strings.Split(s, ",")
	if len(arr) != 4 {
		return token, errors.New("token is error")
	}
	uid, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return token, err
	}
	token.Uid = uid
	token.Password = arr[1]
	token.LastIP = arr[2]
	token.LastDate = arr[3]
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

func UserReadByName(Name string) (row User, err error) {
	data, fields, scanArr := UserMap()
	err = db.Read("User", fields, *scanArr, dbs.H{
		"Name": Name,
	})
	row = *data
	return
}

func UserReadByEmail(Email string) (row User, err error) {
	data, fields, scanArr := UserMap()
	err = db.Read("User", fields, *scanArr, dbs.H{
		"Email": Email,
	})
	row = *data
	return
}

func UserReadByMobile(Mobile string) (row User, err error) {
	data, fields, scanArr := UserMap()
	err = db.Read("User", fields, *scanArr, dbs.H{
		"Mobile": Mobile,
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
	err = db.Find("User", fields, *scanArr, h, "Uid DESC", page, pageSize, func() {
		data.UserFormatSafe() // 格式化
		list = append(list, *data)
	})
	if err != nil {
		return
	}
	return
}

// 干掉敏感信息
func (u *User) UserFormatSafe() {
	u.Password = ""
	u.Salt = ""
}
