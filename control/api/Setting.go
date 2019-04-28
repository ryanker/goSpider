package api

import (
	"github.com/xiuno/gin"

	"../../model"
)

func SettingList(c *gin.Context) {
	list, err := model.SettingList()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"list": list})
}

func SettingSet(c *gin.Context) {
	m := struct {
		Key   string
		Value string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}
	err = model.SettingSet(m.Key, m.Value)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "保存成功")
}

func SettingGet(c *gin.Context) {
	m := struct {
		Key string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}
	Value, err := model.SettingGet(m.Key)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"Value": Value})
}
