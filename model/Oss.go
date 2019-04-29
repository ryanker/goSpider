package model

import (
	"errors"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OssNew() (client *oss.Client, err error) {
	endpoint, ok := Settings["OssEndpoint"]
	accessKeyId, ok2 := Settings["OssAccessKeyId"]
	accessKeySecret, ok3 := Settings["OssAccessKeySecret"]
	if !ok || !ok2 || !ok3 {
		return client, errors.New("OSS 配置信息未设置")
	}

	// 创建 OSSClient 实例
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return client, errors.New("创建 OSSClient 实例失败: " + err.Error())
	}
	return
}

func OssNewBucket() (bucket *oss.Bucket, err error) {
	bucketName, ok := Settings["OssBucketName"]
	if !ok {
		return bucket, errors.New("OSS 配置信息 bucketName 未填写")
	}

	// 创建 OSSClient 实例
	var client *oss.Client
	client, err = OssNew()
	if err != nil {
		return bucket, err
	}

	// 获取存储空间
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		return bucket, errors.New("获取 OSS 存储空间失败: " + err.Error())
	}
	return
}

func OssUpload(objectName, localFileName string) error {
	// 获取存储空间
	bucket, err := OssNewBucket()
	if err != nil {
		return err
	}

	// 上传文件
	err = bucket.PutObjectFromFile(strings.Trim(objectName, "/"), localFileName)
	if err != nil {
		return errors.New("上传文件到 OSS 失败: " + err.Error())
	}
	return nil
}

func OssGetObject(objectName string) (body io.ReadCloser, err error) {
	// 获取存储空间
	bucket, err := OssNewBucket()
	if err != nil {
		return body, err
	}

	// 下载文件到流
	body, err = bucket.GetObject(strings.Trim(objectName, "/"))
	if err != nil {
		return body, errors.New("从 OSS 下载文件失败: " + err.Error())
	}
	body.Close() // 数据读取完成后，获取的流必须关闭
	return
}
