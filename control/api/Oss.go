package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ryanker/gin_v1"

	"../../lib/misc"
	"../../model"
)

func OssFile(c *gin.Context) {
	_, err := UserTokenGet(c)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}

	if model.Settings["OssBucketName"] == "" {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}

	path := strings.Trim(c.Request.URL.Path, "/")
	body, err := model.OssGetObject(path)
	if err != nil {
		c.String(http.StatusNotFound, "404 page not found")
		return
	}

	mp := make(map[string]string)
	mp[".jpg"] = "image/jpeg"
	mp[".jpeg"] = "image/jpeg"
	mp[".gif"] = "image/gif"
	mp[".png"] = "image/png"
	mp[".txt"] = "text/plain"
	mp[".txt"] = "text/plain"
	mp[".htm"] = "text/html"
	mp[".html"] = "text/html"

	ext := filepath.Ext(path)
	contentType, ok := mp[ext]
	extraHeaders := map[string]string{}
	if !ok {
		contentType = "application/octet-stream"
		filename := filepath.Base(path)
		extraHeaders = map[string]string{"Content-Disposition": `attachment; filename="` + filename + `"`}
	}

	c.DataFromReader(http.StatusOK, int64(len(body)), contentType, strings.NewReader(string(body)), extraHeaders)
}

// ================================================
func OssBucketList(c *gin.Context) {
	list, err := model.OssBucketList()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"list": list})
}

func OssBucketInfo(c *gin.Context) {
	m := struct {
		BucketName string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}
	info, err := model.OssBucketInfo(m.BucketName)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"info": info})
}

func OssObjectList(c *gin.Context) {
	m := struct {
		Endpoint   string
		BucketName string
		Prefix     string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}
	files, dirs, err := model.OssObjectList(m.Endpoint, m.BucketName, m.Prefix)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"files": files, "dirs": dirs})
}

func OssObjectCount(c *gin.Context) {
	m := struct {
		Endpoint   string
		BucketName string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确："+err.Error())
		return
	}
	size, fileNum, requests, err := model.OssObjectCount(m.Endpoint, m.BucketName)
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"size": misc.HumanSize(uint64(size)), "fileNum": fileNum, "requests": requests})
}
