package model

import (
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
		time.Sleep(5 * time.Second)

		// 读取需要采集的规则
		list, err := RuleList(dbs.H{"Status >": 0}, 0, 0)
		if err != nil {
			cronErrorLog("读取采集规则失败: " + err.Error())
			continue
		}

		for _, row := range list {
			// 打开采集入库数据库
			dbFile := "./db/" + row.DateBase + ".db"
			dbc, err := dbs.Open(dbFile)
			if err != nil {
				cronErrorLog("打开采集入库数据库失败: " + err.Error())
				continue
			}

			// 读取规则参数
			ListData, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
			if err != nil {
				cronErrorLog("读取列表采集参数失败: " + err.Error())
				continue
			}
			// ContentData, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
			// if err != nil {
			// 	cronErrorLog("读取内容采集参数失败: " + err.Error())
			// 	return
			// }

			getListAll(dbc, &ListData, &row) // 列表：抓取所有列表

			// 判断任务状态，进行相应处理
			if row.Status == 1 {
				cronLog("采集一次完成，关闭采集")
				// 采集一次，完成后关闭采集
				err := RuleUpdate(dbs.H{"Status": 0}, row.Rid)
				if err != nil {
					cronErrorLog("更新采集任务状态失败: " + err.Error())
					continue
				}
			} else if row.Status == 2 {
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
	url := strings.Replace(row.ListUrl, "{page}", strconv.FormatInt(page, 10), -1)
	bodyByte, i, err := misc.HttpGetRetry(url)
	cronLog("请求链接: %v, 请求次数: %v, 耗时: %v", url, i, time.Since(t))
	if err != nil {
		cronErrorLog("抓取页面失败: " + err.Error() + ", Url: " + url)
		return
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		cronErrorLog("解析页面失败: " + err.Error() + ", Url: " + url)
		return
	}

	se := dom.Find(row.ListRule)
	se.Each(func(i int, sel *goquery.Selection) {
		data := dbs.H{}
		for _, v := range *ParamList {
			sel2 := sel.Find(v.Rule)

			// 匹配数据
			value := ""
			if v.ValueType == "Html" {
				value, _ = sel2.Html()
			} else if v.ValueType == "Text" {
				value = sel2.Text()
			} else if v.ValueType == "Attr" {
				value, _ = sel2.Attr(v.ValueAttr)
			} else {
				value, _ = sel2.Html()
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
				cronErrorLog("列表查询重复入库失败: " + err.Error())
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
			cronErrorLog("列表写入数据库失败: " + err.Error())
			return
		}
		cronLog("写入数据库成功: " + strconv.FormatInt(id, 10))
	})
}
