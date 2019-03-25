package main

import (
	"./lib/dbs"
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xiuno/gin"
)

var err error
var db *dbs.DB

type Com99reImgList struct {
	Pid        int64
	Url        string
	Title      string
	ImgUrl     string
	ImgUrlNew  string
	ImgNum     int64
	Views      int64
	Date       string
	CreateDate string
}

type Com99reImgPost struct {
	Pid         int64
	Url         string
	Title       string
	Content     string
	Description string
	ImgNum      int64
	Views       int64
	Date        string
	Author      string
	AuthorHtml  string
	Cate        string
	CateHtml    string
	Tags        string
	TagsHtml    string
	CreateDate  string
}

type Com99reImgPostData struct {
	DataId     int64
	Pid        string
	ImgUrl     string
	ImgUrlNew  string
	CreateDate string
}

func main() {
	runtime.GOMAXPROCS(3)

	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.Static("/upload", "./upload")
	r.SetFuncMap(template.FuncMap{
		"htmlTags": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	r.LoadHTMLGlob("./view/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/read", func(c *gin.Context) {
		Pid := c.DefaultQuery("Pid", "")
		c.HTML(http.StatusOK, "read.html", gin.H{
			"Pid": Pid,
		})
	})

	r.GET("/img/list", func(c *gin.Context) {
		// 映射结构体
		scanF := func() (ptr *Com99reImgList, fields string, args *[]interface{}) {
			row := Com99reImgList{}
			fields, scanArr := dbs.GetSqlRead(dbs.H{
				"Pid":        &row.Pid,
				"Url":        &row.Url,
				"Title":      &row.Title,
				"ImgUrl":     &row.ImgUrl,
				"ImgNum":     &row.ImgNum,
				"Views":      &row.Views,
				"Date":       &row.Date,
				"CreateDate": &row.CreateDate,
			})
			ptr = &row
			args = &scanArr
			return
		}
		data, fields, scanArr := scanF()

		// 读取多条(到结构体)
		rows, err := db.Find("Com99reImgList", fields, dbs.H{}, dbs.H{}, 1, 20)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "-1", "message": "读取数据库失败"})
		}

		var list []Com99reImgList
		for rows.Next() {
			err = rows.Scan(*scanArr...)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": "-1", "message": "数据库 Scan 失败"})
			}
			list = append(list, *data)
		}

		c.JSON(http.StatusOK, gin.H{"code": "0", "list": list})
	})

	r.GET("/img/read", func(c *gin.Context) {
		// 映射结构体
		scanF := func() (ptr *Com99reImgPost, fields string, args *[]interface{}) {
			row := Com99reImgPost{}
			fields, scanArr := dbs.GetSqlRead(dbs.H{
				"Pid":         &row.Pid,
				"Url":         &row.Url,
				"Title":       &row.Title,
				"Content":     &row.Content,
				"Description": &row.Description,
				// "ImgNum":      &row.ImgNum,
				// "Views":       &row.Views,
				"Date":       &row.Date,
				"Author":     &row.Author,
				"Cate":       &row.Cate,
				"Tags":       &row.Tags,
				"CreateDate": &row.CreateDate,
			})
			ptr = &row
			args = &scanArr
			return
		}
		data, fields, scanArr := scanF()

		// 读取(到结构体)
		Pid := c.DefaultQuery("Pid", "")
		err = db.Read("Com99reImgPost", fields, *scanArr, dbs.H{
			"Pid": Pid,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "-1", "message": "读取数据库失败"})
		}

		c.JSON(http.StatusOK, gin.H{"code": "0", "data": *data})
	})

	err := r.Run("0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
}

func init() {
	dbFile := "./db/99re.db"
	db, err = dbs.Open(dbFile)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	// 创建表
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		s, _ := ReadAll("./install/99re.sql")
		_, err = db.Exec(string(s))
		if err != nil {
			panic(err)
		}
	}

	go cron()
}

func cron() {
	for {
		getListAll()
		// getContent()
		time.Sleep(1 * time.Hour)
	}
}

// 根据列表数据库，抓取内容页
func getContent() {
	// 映射结构体
	scanF := func() (ptr *Com99reImgList, fields string, args *[]interface{}) {
		row := Com99reImgList{}
		fields, scanArr := dbs.GetSqlRead(dbs.H{
			"Pid": &row.Pid,
			"Url": &row.Url,
		})
		ptr = &row
		args = &scanArr
		return
	}
	data, fields, scanArr := scanF()

	// 读取多条(到结构体)
	rows, err := db.Find("Com99reImgList", fields, dbs.H{}, dbs.H{}, 1, 20)
	if err != nil {
		panic(err)
	}

	var list []Com99reImgList
	for rows.Next() {
		err = rows.Scan(*scanArr...)
		if err != nil {
			panic(err)
		}
		list = append(list, *data)
	}

	for _, row := range list {
		bodyByte, err := HttpGet(row.Url)
		if err != nil {
			continue
			// panic(err)
		}

		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
		if err != nil {
			log.Fatalln(err)
		}

		size := dom.Find("div.models-slider-big").Size()
		fmt.Println("size: ", size)

		title, _ := dom.Find("div.wrap-title").Find("div.title").Html()
		content, _ := dom.Find("div.models-slider-big").Html()

		u := dom.Find("ul.information").Find("span.value")
		views := u.Eq(0).Text()
		date := u.Eq(1).Text()

		desc := dom.Find("div.description-block").Find("div.wrap-overflow")
		description := desc.Eq(0).Find("div.desc-video").Text()
		imgNum := desc.Eq(1).Find("div.desc-video").Text()

		eq2 := desc.Eq(2).Find("div.name-avtor-dwnl")
		author := eq2.Text()
		authorHtml, _ := eq2.Html()

		eq3 := desc.Eq(3)
		cate := eq3.Text()
		cateHtml, _ := eq3.Html()

		eq4 := desc.Eq(4)
		tags := eq4.Text()
		tagsHtml, _ := eq4.Html()

		// 判断是否已经抓取过
		n, err := db.Count("Com99reImgPost", dbs.H{"url": row.Url})
		if err != nil {
			panic(err)
		}
		if n > 0 {
			// 更新
			_, err = db.Update("Com99reImgPost", dbs.H{
				"Title":       Trim(title),
				"Content":     Trim(content),
				"Description": Trim(description),
				"ImgNum":      Trim(imgNum),
				"Views":       Trim(views),
				"Date":        Trim(date),
				"Author":      Trim(author),
				"AuthorHtml":  Trim(authorHtml),
				"Cate":        Trim(cate),
				"CateHtml":    Trim(cateHtml),
				"Tags":        Trim(tags),
				"TagsHtml":    Trim(tagsHtml),
			}, dbs.H{
				"Pid": row.Pid,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println("Update: ", row.Pid)
			continue
		}

		// 插入
		id, err := db.Insert("Com99reImgPost", dbs.H{
			"Pid":         row.Pid,
			"Url":         row.Url,
			"Title":       Trim(title),
			"Content":     Trim(content),
			"Description": Trim(description),
			"ImgNum":      Trim(imgNum),
			"Views":       Trim(views),
			"Date":        Trim(date),
			"Author":      Trim(author),
			"AuthorHtml":  Trim(authorHtml),
			"Cate":        Trim(cate),
			"CateHtml":    Trim(cateHtml),
			"Tags":        Trim(tags),
			"TagsHtml":    Trim(tagsHtml),
			"CreateDate":  time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Insert: ", id)
	}
}

// 抓取列表页 && 下载图片
func getListAll() {
	// 抓取列表页，入库
	// for page := 1; page <= 268; page++ {
	// 	getList(page)
	// }

	// 下载列表页图片
	downloadListImg()
}

// 抓取列表页，入库
func getList(page int) {
	// url demo: https://99a22.com/albums/1/
	bodyByte, err := HttpGet("https://99a22.com/albums/" + strconv.Itoa(page) + "/")
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(bodyByte))

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		log.Fatalln(err)
	}

	size := dom.Find("div.thumb-content").Size()
	fmt.Println("size: ", size)

	dom.Find("div.thumb-content").Each(func(i int, selection *goquery.Selection) {
		// fmt.Println(selection.Text())
		// fmt.Println(selection.Html())

		a := selection.Find("a")
		url, _ := a.Attr("href")
		title, _ := a.Attr("title")
		imgUrl, _ := a.Find("img").Attr("src")
		imgNum := a.Find("span.duration").Text()
		views := a.Find("span.views").Text()
		date := a.Find("span.data").Text()

		// fmt.Println(url + " " + title)

		// 判断是否已经抓取过
		url = "https://99a22.com" + url
		n, err := db.Count("Com99reImgList", dbs.H{"url": url})
		if err != nil {
			panic(err)
		}
		if n > 0 {
			return
		}

		// 插入
		id, err := db.Insert("Com99reImgList", dbs.H{
			"Url":        url,
			"Title":      title,
			"ImgUrl":     imgUrl,
			"ImgNum":     strings.TrimRight(imgNum, " photos"),
			"Views":      strings.TrimRight(views, " 次浏览"),
			"Date":       date,
			"CreateDate": time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Insert: ", id)
	})
}

// 下载列表页图片
func downloadListImg() {
	// 映射结构体
	scanF := func() (ptr *Com99reImgList, fields string, args *[]interface{}) {
		row := Com99reImgList{}
		fields, scanArr := dbs.GetSqlRead(dbs.H{
			"Pid":    &row.Pid,
			"ImgUrl": &row.ImgUrl,
		})
		ptr = &row
		args = &scanArr
		return
	}
	data, fields, scanArr := scanF()

	// 读取多条(到结构体)
	rows, err := db.Find("Com99reImgList", fields, dbs.H{}, dbs.H{}, 1, 20)
	if err != nil {
		panic(err)
	}

	var list []Com99reImgList
	for rows.Next() {
		err = rows.Scan(*scanArr...)
		if err != nil {
			panic(err)
		}
		list = append(list, *data)
	}

	path := "/upload/com99re/img/preview/"
	dir := "." + path
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range list {
		bodyByte, err := HttpGet(row.ImgUrl)
		if err != nil {
			continue
			// panic(err)
		}

		// 存放文件名
		filename := fmt.Sprintf("%x", md5.Sum([]byte(row.ImgUrl))) + filepath.Ext(row.ImgUrl)
		SaveFile(dir+filename, bodyByte)

		// 更新
		_, err = db.Update("Com99reImgList", dbs.H{
			"ImgUrlNew": path + filename,
		}, dbs.H{
			"Pid": row.Pid,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Update: ", row.Pid)
	}
}

// 保存列表图片
func SaveFile(filename string, bodyByte []byte) {
	fmt.Println("Save: ", filename)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(bodyByte); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func HttpGet(url string) (bodyByte []byte, err error) {
	for i := 1; i <= 5; i++ {
		fmt.Println(url, "times: ", i)
		bodyByte, err = httpGet(url)
		if err != nil {
			time.Sleep(1)
			continue
		}
		return
	}
	return
}

func httpGet(url string) (bodyByte []byte, err error) {
	var c = &http.Client{}
	c.Timeout = 10 * time.Second

	resp, err := c.Get(url)
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
