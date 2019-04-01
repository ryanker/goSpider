package misc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// 打开文件
func FileOpen(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// 保存文件
func FileSave(name string, bodyByte []byte) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.Write(bodyByte); err != nil {
		fmt.Println(err)
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}

// 抓取网页（支持重试）
func HttpGetRetry(url string) (bodyByte []byte, err error) {
	var interval time.Duration = 1 // 间隔
	for i := 1; i <= 6; i++ {
		bodyByte, err = HttpGet(url)
		if err != nil {
			time.Sleep(interval * time.Second)
			interval *= 2 // 每次重试，延迟双倍时间，增加重试成功率
			continue
		}
		return
	}
	return
}

// 抓取网页
func HttpGet(url string) (bodyByte []byte, err error) {
	c := &http.Client{}
	c.Timeout = 20 * time.Second // 请求网页，20秒足够

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 AppleWebKit/537.36 Chrome/73.0.3683.86 Safari/537.36")
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyByte, err = ioutil.ReadAll(resp.Body)
	return
}

func Trim(s string) string {
	return strings.Trim(s, " \t\r\n")
}
