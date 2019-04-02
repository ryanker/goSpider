package api

import "github.com/xiuno/gin"

func DatabaseSet(c *gin.Context) {
	m := struct {
		Rid    int64
		Status int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}
}
