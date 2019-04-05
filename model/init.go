package model

import (
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"../lib/dbs"
	"../lib/misc"
)

var err error
var db *dbs.DB

func init() {
	dbs.LogFile = "./log/db.log"
	dbs.ErrorLogFile = "./log/db.error.log"

	dbFile := "./db/data.db"
	db, err = dbs.Open(dbFile)
	if err != nil {
		panic(err)
	}

	// 文件不存在，则创建表
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		s, _ := misc.FileOpen("./install/install.sql")
		_, err = db.Exec(string(s))
		if err != nil {
			panic(err)
		}
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	go cron()
}

// 任务
func cron() {
	for {
		time.Sleep(10 * time.Second)

		// 读取需要采集的规则
		list, err := RuleList(dbs.H{"Status >": 1}, 0, 0)
		if err != nil {
			cronErrorLog("读取采集规则失败: %v", err.Error())
			continue
		}

		for _, row := range list {
			// 打开采集入库数据库
			dbFile := "./db/" + row.DateBase + ".db"
			dbc, err := dbs.Open(dbFile)
			if err != nil {
				cronErrorLog("打开采集入库数据库失败: %v", err.Error())
				continue
			}

			// 读取规则参数
			ListData, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
			if err != nil {
				cronErrorLog("读取列表采集参数失败: %v", err.Error())
				continue
			}
			ContentData, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
			if err != nil {
				cronErrorLog("读取内容采集参数失败: %v", err.Error())
				return
			}

			getListAll(dbc, &ListData, &row)    // 列表页：抓取所有列表
			getContent(dbc, &ContentData, &row) // 内容页：抓取内容页

			// 判断任务状态，进行相应处理
			if row.Status == 2 {
				cronLog("采集一次完成，关闭采集")
				// 采集一次，完成后关闭采集
				err := RuleUpdate(dbs.H{"Status": 1}, row.Rid)
				if err != nil {
					cronErrorLog("更新采集任务状态失败: %v", err.Error())
					continue
				}
			} else if row.Status == 3 {
				// 间隔采集，完成后等待下次采集
				cronLog("间隔采集完成，等待下次采集")
				time.Sleep(time.Duration(row.IntervalHour) * time.Hour)
			}
		}
	}
}

func cronErrorLog(format string, args ...interface{}) {
	misc.FileLogWrite("./log/cronError.log", format, args...)
}

func cronLog(format string, args ...interface{}) {
	misc.FileLogWrite("./log/cron.log", format, args...)
}

// 列表页：抓取所有列表页
func getListAll(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule) {
	// 抓取列表页，入库
	for page := row.ListPageStart; page <= row.ListPageEnd; page += row.ListPageSize {
		getList(dbc, ParamList, row, page)
	}
}

// 列表页：抓取列表页
func getList(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule, page int64) {
	t := time.Now()
	Url := strings.Replace(row.ListUrl, "{page}", strconv.FormatInt(page, 10), -1)
	bodyByte, i, err := misc.HttpGetRetry(Url)
	cronLog("请求链接: %v, 请求次数: %v, 耗时: %v", Url, i, time.Since(t))
	if err != nil {
		cronErrorLog("抓取页面失败: %v", err.Error())
		return
	}

	t2 := time.Now()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		cronErrorLog("解析页面失败: %v, Url:%v", err.Error(), Url)
		return
	}

	se := doc.Find(row.ListRule)
	se.Each(func(i int, s *goquery.Selection) {
		data := dbs.H{}
		data["Status"] = 1 // 待采集
		for _, v := range *ParamList {
			dom := s.Find(v.Rule)

			// 匹配数据
			value := ""
			if v.ValueType == "Html" {
				value, _ = dom.Html()
			} else if v.ValueType == "Text" {
				value = dom.Text()
			} else if v.ValueType == "Attr" {
				value, _ = dom.Attr(v.ValueAttr)
			} else {
				value, _ = dom.Html()
			}

			// 数据过滤
			if v.FilterType == "Trim" {
				value = misc.Trim(value)
			} else if v.FilterType == "Reg" {
				re := regexp.MustCompile(v.FilterRegexp)
				value = re.ReplaceAllString(value, v.FilterRepl)
			}

			data[v.Field] = value
		}

		// 根据Url，判断是否重复
		Url, ok := data["Url"]
		if ok {
			n, err := dbc.Count("List", dbs.H{"Url": Url})
			if err != nil {
				cronErrorLog("列表查询重复入库失败: %v", err.Error())
				return
			}
			if n > 0 {
				cronLog("重复Url: %v", Url)
				return
			}
		}

		// 写入数据库
		id, err := dbc.Insert("List", data)
		if err != nil {
			cronErrorLog("列表写入数据库失败: %v", err.Error())
			return
		}
		cronLog("写入数据库成功: %v", id)
	})
	cronLog("第 %v 页入库完成, 耗时： %v", page, time.Since(t2))
}

// 内容页：抓取内容页
func getContent(dbc *dbs.DB, ContentData *[]RuleParam, row *Rule) {
	rows, err := dbc.Find("List", "Url", dbs.H{}, "ListId DESC", 0, 2)
	if err != nil {
		cronErrorLog("列表读取失败: %v", err.Error())
		return
	}
	for rows.Next() {
		var Url string
		err = rows.Scan(&Url)
		if err != nil {
			cronErrorLog("Url绑定失败: %v", err.Error())
			return
		}

		// 效验链接
		if Url == "" {
			cronErrorLog("内容页 Url 为空")
			continue
		}
		// 修正链接
		if Url[0:1] == "/" {
			u, _ := url.Parse(row.ListUrl)
			Url = u.Scheme + "://" + u.Host + Url
		}

		t := time.Now()
		bodyByte, i, err := misc.HttpGetRetry(Url)
		cronLog("请求链接: %v, 请求次数: %v, 耗时: %v", Url, i, time.Since(t))
		if err != nil {
			cronErrorLog("抓取页面失败: %v", err.Error())
			return
		}

		t2 := time.Now()
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
		if err != nil {
			cronErrorLog("解析页面失败: %v, Url:%v", err.Error(), Url)
			return
		}

		data := dbs.H{}
		data["Url"] = Url
		for _, m := range *ContentData {
			dom := doc.Find(m.Rule).Eq(0)
			value := ""
			if m.ValueType == "Html" {
				value, _ = dom.Html()
			} else if m.ValueType == "Text" {
				value = dom.Text()
			} else if m.ValueType == "Attr" {
				value, _ = dom.Attr(m.ValueAttr)
			} else {
				value, _ = dom.Html()
			}

			if m.FilterType == "Trim" {
				value = misc.Trim(value)
			} else if m.FilterType == "Reg" {
				re := regexp.MustCompile(m.FilterRegexp)
				value = re.ReplaceAllString(value, m.FilterRepl)
			}

			data[m.Field] = value
		}

		// 根据Url，判断是否重复
		n, err := dbc.Count("Content", dbs.H{"Url": Url})
		if err != nil {
			cronErrorLog("内容页查询重复入库失败: %v", err.Error())
			return
		}
		if n > 0 {
			cronLog("重复Url: %v", Url)
			return
		}

		// 写入数据库
		id, err := dbc.Insert("Content", data)
		if err != nil {
			cronErrorLog("列表写入数据库失败: %v", err.Error())
			return
		}
		cronLog("写入数据库成功: %v, 耗时: %v", id, time.Since(t2))
	}
}
