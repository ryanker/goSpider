package model

import (
	"errors"
	"io/ioutil"
	"strings"
	"time"

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
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return errors.New("上传文件到 OSS 失败: " + err.Error())
	}
	return nil
}

func OssGetObject(objectName string) (b []byte, err error) {
	// 获取存储空间
	bucket, err := OssNewBucket()
	if err != nil {
		return b, err
	}

	// 下载文件到流
	body, err := bucket.GetObject(objectName)
	if err != nil {
		return b, errors.New("从 OSS 下载文件失败: " + err.Error())
	}
	defer body.Close()

	b, err = ioutil.ReadAll(body)
	if err != nil {
		return b, err
	}
	return
}

// ===================================================================
func OssBucketList() (list []oss.BucketProperties, err error) {
	client, err := OssNew()
	if err != nil {
		return list, err
	}

	marker := ""
	for {
		lsRes, err := client.ListBuckets(oss.Marker(marker))
		if err != nil {
			return list, err
		}

		// 默认每次返回100条
		for _, bucket := range lsRes.Buckets {
			list = append(list, bucket)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return
}

func OssBucketInfo(bucketName string) (info oss.BucketInfo, err error) {
	client, err := OssNew()
	if err != nil {
		return info, err
	}
	res, err := client.GetBucketInfo(bucketName)
	if err != nil {
		return info, err
	}
	return res.BucketInfo, err
}

func OssObjectList(endpoint, bucketName, prefix string) (files []map[string]interface{}, dirs []string, err error) {
	accessKeyId, ok := Settings["OssAccessKeyId"]
	accessKeySecret, ok2 := Settings["OssAccessKeySecret"]
	if !ok || !ok2 {
		err = errors.New("OSS 配置信息未设置")
		return
	}

	// 优化是否走内网
	Endpoint, ok := Settings["OssEndpoint"]
	if ok && strings.Contains(Endpoint, "-internal.aliyuncs.com") {
		if Endpoint == strings.Replace(endpoint, ".aliyuncs.com", "-internal.aliyuncs.com", -1) {
			endpoint = Endpoint
		}
	}

	// 创建 OSSClient 实例
	var client *oss.Client
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return
	}

	// 获取存储空间
	var bucket *oss.Bucket
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		return
	}

	// 列举所有文件
	marker := oss.Marker("")
	maxKeys := oss.MaxKeys(1000)
	delimiter := oss.Delimiter("/")
	Prefix := oss.Prefix(prefix)
	for {
		var lsRes oss.ListObjectsResult
		lsRes, err = bucket.ListObjects(marker, maxKeys, delimiter, Prefix)
		if err != nil {
			return
		}

		for _, object := range lsRes.Objects {
			files = append(files, map[string]interface{}{
				"Key":          strings.TrimPrefix(object.Key, prefix),
				"Size":         object.Size,
				"LastModified": object.LastModified.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			})
		}

		for _, dir := range lsRes.CommonPrefixes {
			dirs = append(dirs, strings.TrimPrefix(dir, prefix))
		}

		if lsRes.IsTruncated {
			marker = oss.Marker(lsRes.NextMarker)
		} else {
			break
		}
	}
	return
}

// 统计全部 存储用量 文件数量 请求次数 (文件多，耗时特别长)
func OssObjectCount(endpoint, bucketName string) (size, fileNum, requests int64, err error) {
	accessKeyId, ok := Settings["OssAccessKeyId"]
	accessKeySecret, ok2 := Settings["OssAccessKeySecret"]
	if !ok || !ok2 {
		err = errors.New("OSS 配置信息未设置")
		return
	}

	// 创建 OSSClient 实例
	var client *oss.Client
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return
	}

	// 获取存储空间
	var bucket *oss.Bucket
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		return
	}

	// 列举所有文件
	marker := ""
	maxKeys := 1000
	for {
		requests++
		var lsRes oss.ListObjectsResult
		lsRes, err = bucket.ListObjects(oss.Marker(marker), oss.MaxKeys(maxKeys))
		if err != nil {
			return
		}

		for _, object := range lsRes.Objects {
			size += object.Size
			fileNum++
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return
}

// 列出全部文件，可以用于全部删除，或者全部下载
func OssObjectListAll(endpoint, bucketName string) (list []oss.ObjectProperties, err error) {
	accessKeyId, ok := Settings["OssAccessKeyId"]
	accessKeySecret, ok2 := Settings["OssAccessKeySecret"]
	if !ok || !ok2 {
		return list, errors.New("OSS 配置信息未设置")
	}

	// 创建 OSSClient 实例
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return list, errors.New("创建 OSSClient 实例失败: " + err.Error())
	}

	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return list, err
	}

	// 列举所有文件
	marker := ""
	maxKeys := 1000
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker), oss.MaxKeys(maxKeys))
		if err != nil {
			return list, err
		}

		for _, object := range lsRes.Objects {
			list = append(list, object)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return
}
