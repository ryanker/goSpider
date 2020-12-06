package api

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ryanker/gin_v1"

	"../../model"
)

func SettingSave(c *gin.Context) {
	m := struct {
		OssEndpoint        string
		OssAccessKeyId     string
		OssAccessKeySecret string
		OssBucketName      string
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	mp := make(map[string]string)
	mp["OssEndpoint"] = m.OssEndpoint
	mp["OssAccessKeyId"] = m.OssAccessKeyId
	mp["OssAccessKeySecret"] = m.OssAccessKeySecret
	mp["OssBucketName"] = m.OssBucketName

	for k, v := range mp {
		err = model.SettingSet(k, v)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	}

	// 创建 OSSClient 实例
	client, err := oss.New(m.OssEndpoint, m.OssAccessKeyId, m.OssAccessKeySecret)
	if err != nil {
		c.Message("-1", "创建 OSSClient 实例失败: "+err.Error())
		return
	}

	if m.OssBucketName != "" {
		// 检测存储空间是否存在
		isExist, err := client.IsBucketExist(m.OssBucketName)
		if err != nil {
			c.Message("-1", "检测存储空间是否存在失败: "+err.Error())
			return
		}

		if !isExist {
			// 创建存储空间
			err = client.CreateBucket(m.OssBucketName)
			if err != nil {
				c.Message("-1", "创建存储空间失败: "+err.Error())
				return
			}
		}
	}

	// 加载配置信息到内存
	err = model.SettingInit()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}

	c.Message("0", "保存成功")
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

func SettingList(c *gin.Context) {
	data, err := model.SettingList()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success", gin.H{"data": data})
}

func SettingInit(c *gin.Context) {
	err := model.SettingInit()
	if err != nil {
		c.Message("-1", err.Error())
		return
	}
	c.Message("0", "success")
}
