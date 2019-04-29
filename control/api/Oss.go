package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/xiuno/gin"

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
