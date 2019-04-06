package misc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

// 写入文件日志
func FileLogWrite(FilePath string, format string, args ...interface{}) {
	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	str := fmt.Sprintf("%v | %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...),
	)
	if _, err := f.Write([]byte(str)); err != nil {
		fmt.Println(err)
	}

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}

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
func HttpGetRetry(url string) (bodyByte []byte, i int64, err error) {
	var interval time.Duration = 1 // 间隔
	for i = 1; i <= 6; i++ {
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
	c.Timeout = 5 * time.Second // 请求网页，5 秒足够

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

// 下载文件
func DownloadFile(Url string, DstFile string) (size int64, err error) {
	client := grab.NewClient()
	req, err := grab.NewRequest(DstFile, Url)
	if err != nil {
		return 0, err
	}
	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		return 0, err
	}
	return resp.Size, err
}

func Trim(s string) string {
	return strings.Trim(s, " \t\r\n")
}

func UrlFix(Url string, BaseUrl string) string {
	if Url[0:1] == "/" {
		u, _ := url.Parse(BaseUrl)
		Url = u.Scheme + "://" + u.Host + Url
	}
	return Url
}
