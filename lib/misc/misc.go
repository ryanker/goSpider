package misc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
	var interval time.Duration = 2 // 间隔
	for i = 1; i <= 6; i++ {
		bodyByte, err = HttpGet(url)
		if err != nil {
			time.Sleep(interval * time.Second)
			interval *= 1 // 每次重试，延迟双倍时间，增加重试成功率
			continue
		}
		return
	}
	return
}

// 抓取网页
func HttpGet(url string) (bodyByte []byte, err error) {
	c := &http.Client{}
	c.Timeout = 20 * time.Second // 请求网页，20 秒足够

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
	err = os.MkdirAll(filepath.Dir(DstFile), os.ModePerm)
	if err != nil {
		return
	}
	var req *grab.Request
	var resp *grab.Response
	for i := 1; i <= 5; i++ {
		client := grab.NewClient()
		client.UserAgent = "Mozilla/5.0 AppleWebKit/537.36 Chrome/73.0.3683.86 Safari/537.36"
		client.HTTPClient.Timeout = 20 * time.Duration(i) * time.Second
		req, err = grab.NewRequest(DstFile, Url)
		if err != nil {
			continue
		}
		resp = client.Do(req)
		err = resp.Err()
		if err != nil {
			continue
		}

		// 重新获取文件大小
		if resp.Size < 1 {
			fileInfo, _ := os.Stat(DstFile)
			resp.Size = fileInfo.Size()
		}

		return resp.Size, err
	}
	return resp.Size, err
}

func HumanSize(n uint64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	} else if n < 1024*1024 {
		return fmt.Sprintf("%.2f K", float64(n)/1024)
	} else if n < 1024*1024*1024 {
		return fmt.Sprintf("%.2f M", float64(n)/(1024*1024))
	} else if n < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f G", float64(n)/(1024*1024*1024))
	} else if n < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f T", float64(n)/(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2f P", float64(n)/(1024*1024*1024*1024*1024))
	}
}

func Trim(s string) string {
	return strings.Trim(s, " \t\r\n")
}

func UrlFix(Url string, BaseUrl string) string {
	if Url[0:1] == "/" {
		u, _ := url.Parse(BaseUrl)
		Url = u.Scheme + "://" + u.Host + Url
	} else if Url[0:4] != "http" {
		u, _ := url.Parse(BaseUrl)
		Url = u.Scheme + "://" + u.Host + filepath.Dir(u.Path) + "/" + Url
	}
	return Url
}

func UrlExt(Url string) string {
	u, _ := url.Parse(Url)
	return filepath.Ext(u.Path)
}

func StrClear(value, FilterType, FilterRegexp, FilterRepl string) string {
	if FilterType == "Reg" {
		re := regexp.MustCompile(FilterRegexp)
		if FilterRepl == "" {
			value = re.FindString(value)
		} else {
			value = re.ReplaceAllString(value, FilterRepl)
		}
	}
	return Trim(value)
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Base16Encode(s string) string {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return fmt.Sprintf("%s", dst)
}

func Base16Decode(s string) (string, error) {
	src := []byte(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", dst[:n]), nil
}
