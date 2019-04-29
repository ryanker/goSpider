package api

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/xiuno/gin"

	"../../model"
)

func OssRead(c *gin.Context) error {
	if model.Settings["OssBucketName"] == "" {
		return errors.New("OssBucketName is empty")
	}

	path := strings.Trim(c.Request.URL.Path, "/")
	body, err := model.OssGetObject(path)
	if err != nil {
		return err
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
	c.Abort()
	return nil
}
