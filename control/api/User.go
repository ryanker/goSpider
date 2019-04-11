package api

import (
	"math/rand"
	"strconv"

	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../lib/misc"
	"../../model"
)

func UserCreate(c *gin.Context) {
	m := model.User{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	if m.Name == "" {
		c.Message("-1", "用户名不能为空")
		return
	}
	if m.Mobile == "" {
		c.Message("-1", "手机不能为空")
		return
	}
	if m.Password == "" {
		c.Message("-1", "密码不能为空")
		return
	}
	if m.Gid == 0 {
		m.Gid = 2
	}

	// 判断是否重名
	u, _ := model.UserReadByName(m.Name)
	if u.Uid > 0 {
		c.Message("-1", "用户名已存在")
		return
	}
	u, _ = model.UserReadByMobile(m.Mobile)
	if u.Uid > 0 {
		c.Message("-1", "手机号已存在")
		return
	}

	// password 为 md5 以后的数据
	m.Salt = strconv.Itoa(rand.Intn(999999))
	m.Password = misc.Md5(m.Password + m.Salt)
	m.CreateIP = c.ClientIP()

	uid, err := model.UserCreate(dbs.H{
		"Gid":        m.Gid,
		"Name":       m.Name,
		"Email":      m.Email,
		"Mobile":     m.Mobile,
		"Password":   m.Password,
		"Salt":       m.Salt,
		"LoginNum":   0,
		"LastIP":     m.LastIP,
		"LastDate":   "0001-01-01 08:00:00",
		"CreateIP":   m.CreateIP,
		"UpdateDate": "0001-01-01 08:00:00",
	})
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "添加成功", gin.H{"uid": uid})
}

func UserUpdate(c *gin.Context) {
	m := model.User{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	if m.Uid == 0 {
		c.Message("-1", "uid 不能为空")
		return
	}
	if m.Name == "" {
		c.Message("-1", "用户名不能为空")
		return
	}
	if m.Mobile == "" {
		c.Message("-1", "手机不能为空")
		return
	}
	if m.Gid == 0 {
		m.Gid = 2
	}

	// 判断是否重名
	u, _ := model.UserReadByName(m.Name)
	if u.Uid > 0 && u.Uid != m.Uid {
		c.Message("-1", "用户名已存在")
		return
	}
	u, _ = model.UserReadByMobile(m.Mobile)
	if u.Uid > 0 && u.Uid != m.Uid {
		c.Message("-1", "手机号已存在")
		return
	}

	h := dbs.H{
		"Gid":      m.Gid,
		"Name":     m.Name,
		"Email":    m.Email,
		"Mobile":   m.Mobile,
		"Password": u.Password,
		"Salt":     u.Salt,
	}
	if m.Password != "" {
		if len(m.Password) == 32 && m.Password != misc.Md5("") {
			h["Salt"] = strconv.Itoa(rand.Intn(999999))
			h["Password"] = misc.Md5(m.Password + u.Salt)
		}
	}

	err = model.UserUpdate(h, m.Uid)
	if err != nil {
		c.Message("-1", "更新数据库失败："+err.Error())
		return
	}

	c.Message("0", "修改成功")
}

func UserRead(c *gin.Context) {
	m := model.User{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	u, err := model.UserRead(m.Uid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	u.UserFormatSafe()

	c.Message("0", "success", u)
}

func UserDelete(c *gin.Context) {
	m := model.User{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	err = model.UserDelete(m.Uid)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "删除成功")

}

func UserList(c *gin.Context) {
	m := struct {
		model.User
		Page int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}

	// 筛选条件
	h := dbs.H{}
	if m.Gid > 0 {
		h["Gid"] = m.Gid
	}
	if m.Name != "" {
		h["Name LIKE"] = "%" + m.Name + "%"
	}
	if m.Email != "" {
		h["Name LIKE"] = "%" + m.Email + "%"
	}
	if m.Mobile != "" {
		h["Name LIKE"] = "%" + m.Mobile + "%"
	}

	total, err := model.UserCount(h)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	list, err := model.UserList(h, m.Page, 20)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "success", gin.H{"total": total, "list": list})
}
