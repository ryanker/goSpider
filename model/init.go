package model

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"../lib/dbs"
	"../lib/misc"
)

type ListUrl struct {
	Lid int64
	Url string
}

type ListHtml struct {
	Lid  int64
	Html string
}

type ListDown struct {
	Id     int64
	Lid    int64
	Field  string
	OldUrl string
	NewUrl string
	Sort   int64
}

var err error
var db *dbs.DB
var dbLog *dbs.DB

var logSql = `CREATE TABLE Log
(
  LogId      INTEGER PRIMARY KEY AUTOINCREMENT,   -- 日志ID
  Status     INTEGER        NOT NULL DEFAULT '0', -- 日志状态 1:普通日志 2:错误日志
  Runtime    DECIMAL(10, 4) NOT NULL DEFAULT '0', -- 执行耗时
  Message    TEXT                    DEFAULT '',  -- 日志内容
  CreateDate DATETIME                DEFAULT CURRENT_TIMESTAMP
);`

var loc = "?_loc=Asia%2FShanghai"

func init() {
	dbs.LogFile = "./log/db.log"
	dbs.ErrorLogFile = "./log/db.error.log"

	InitDB()    // 打开主库
	InitDbLog() // 打开日志库

	go cron()
}

// 打开主库
func InitDB() {
	dbFile := "./db/data.db"
	db, err = dbs.Open(dbFile + loc)
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
}

// 打开日志库
func InitDbLog() {
	dbFile := "./db/log.db"
	dbLog, err = dbs.Open(dbFile + loc)
	if err != nil {
		panic(err)
	}

	// 文件不存在，则创建表
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		_, err = dbLog.Exec(logSql)
		if err != nil {
			panic(err)
		}
	}

	if err = dbLog.Ping(); err != nil {
		panic(err)
	}
}

func cronErrorLog(t time.Duration, format string, args ...interface{}) {
	_, err := dbLog.Insert("Log", dbs.H{
		"Status":     2,
		"Runtime":    fmt.Sprintf("%.4f", float64(t)/float64(time.Second)),
		"Message":    fmt.Sprintf(format, args...),
		"CreateDate": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		misc.FileLogWrite("./log/cronError.log", "cronErrorLog error: "+err.Error())
	}
}

func cronLog(t time.Duration, format string, args ...interface{}) {
	_, err := dbLog.Insert("Log", dbs.H{
		"Status":     1,
		"Runtime":    fmt.Sprintf("%.4f", float64(t)/float64(time.Second)),
		"Message":    fmt.Sprintf(format, args...),
		"CreateDate": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		misc.FileLogWrite("./log/cronError.log", "cronLog error: "+err.Error())
	}
}

// 任务
func cron() {
	for {
		time.Sleep(1 * time.Minute)

		// 读取需要采集的规则
		list, err := RuleList(dbs.H{"Status >": 1}, 0, 0)
		if err != nil {
			cronErrorLog(0, "读取采集规则失败: %v", err.Error())
			continue
		}

		for _, row := range list {
			// 打开采集入库数据库
			dbFile := "./db/" + row.Database + ".db"
			dbc, err := dbs.Open(dbFile + loc)
			if err != nil {
				cronErrorLog(0, "打开采集入库数据库失败: %v", err.Error())
				continue
			}

			// 读取规则参数
			ParamList, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
			if err != nil {
				cronErrorLog(0, "读取列表采集参数失败: %v", err.Error())
				continue
			}
			ParamContent, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
			if err != nil {
				cronErrorLog(0, "读取内容采集参数失败: %v", err.Error())
				return
			}

			if row.IsList == 1 {
				getListAll(dbc, &ParamList, &row) // 列表页：抓取所有列表
			}
			if row.IsContent == 1 {
				getContent(dbc, &ParamContent, &row) // 内容页：抓取内容页
			}
			if row.IsListDownAna == 1 {
				downInitList(dbc, &ParamList, &row) // 下载：分析列表页下载地址
			}
			if row.IsContentDownAna == 1 {
				downInitContent(dbc, &ParamContent, &row) // 下载：分析内容页下载地址
			}
			if row.IsListDownRun == 1 {
				downList(dbc, &row) // 下载：列表页资源下载
			}
			if row.IsContentDownRun == 1 {
				downContent(dbc, &row) // 下载：内容页资源下载
			}

			// 判断任务状态，进行相应处理
			if row.Status == 2 {
				cronLog(0, "采集一次完成，关闭采集")
				// 采集一次，完成后关闭采集
				err := RuleUpdate(dbs.H{"Status": 1}, row.Rid)
				if err != nil {
					cronErrorLog(0, "更新采集任务状态失败: %v", err.Error())
					continue
				}
			} else if row.Status == 3 {
				// 间隔采集，完成后等待下次采集
				cronLog(0, "间隔采集完成，等待下次采集")
				time.Sleep(time.Duration(row.IntervalHour) * time.Hour)
			}
		}
	}
}

// 列表页：抓取所有列表页
func getListAll(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule) {
	// 特殊列表 Url 抓取
	if row.ListSpecialUrl != "" {
		Urls := strings.Split(row.ListSpecialUrl, "\n")
		for _, Url := range Urls {
			Url = misc.Trim(Url)
			if Url != "" {
				getList(dbc, ParamList, row, Url)
			}
		}
	}

	// 规律列表 Url 抓取
	for page := row.ListPageStart; page <= row.ListPageEnd; page += row.ListPageSize {
		Url := strings.Replace(row.ListUrl, "{page}", strconv.FormatInt(page, 10), -1)
		getList(dbc, ParamList, row, Url)
	}
}

// 列表页：抓取列表页
func getList(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule, Url string) {
	// 抓取页面
	t2 := time.Now()
	bodyByte, i, err := misc.HttpGetRetry(Url)
	if err != nil {
		cronErrorLog(time.Since(t2), "抓取列表页失败: %v, Url: %v", err.Error(), Url)
		return
	}
	cronLog(time.Since(t2), "抓取列表页成功, 请求次数: %v, Url: %v", i, Url)

	// 解析代码
	t3 := time.Now()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
	if err != nil {
		cronErrorLog(0, "解析页面失败: %v, Url:%v", err.Error(), Url)
		return
	}

	// 代码匹配
	doc.Find(row.ListRule).Each(func(i int, s *goquery.Selection) {
		data := dbs.H{}
		data["Status"] = 1 // 待采集
		for _, v := range *ParamList {
			dom := s.Find(v.Rule)

			// 匹配字段
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
			value = misc.StrClear(value, v.FilterType, v.FilterRegexp, v.FilterRepl)
			data[v.Field] = value
		}

		// 根据Url，判断是否重复
		Url, ok := data["Url"]
		if ok {
			n, err := dbc.Count("List", dbs.H{"Url": Url})
			if err != nil {
				cronErrorLog(0, "列表页查询重复入库失败: %v", err.Error())
				return
			}
			if n > 0 {
				cronLog(0, "重复Url: %v", Url)
				return
			}
		}

		// 写入数据库
		id, err := dbc.Insert("List", data)
		if err != nil {
			cronErrorLog(0, "列表页写入数据库失败: %v", err.Error())
			return
		}
		cronLog(0, "列表页写入数据库成功: %v", id)
	})
	cronLog(time.Since(t3), "Url: %v 抓取完成", Url)
}

// 内容页：抓取内容页
func getContent(dbc *dbs.DB, ParamContent *[]RuleParam, row *Rule) {
	t := time.Now() // 记时开始
	repeatNum := 0  // 重复入库数
	newNum := 0     // 新增入库数
	errorNum := 0   // 错误入库数

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("List", where)
	if err != nil {
		cronErrorLog(0, "下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(100)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		st := ListUrl{}
		var list []ListUrl
		err = dbc.Find("List", "`Lid`,`Url`", []interface{}{&st.Lid, &st.Url}, where, "Lid DESC", page, pageSize, func() {
			list = append(list, st)
		})
		if err != nil {
			cronErrorLog(0, "列表读取失败: %v", err.Error())
			return
		}

		for _, lv := range list {
			// 效验链接
			if lv.Url == "" {
				errorNum++
				cronErrorLog(0, "内容页 Url 为空")
				continue
			}
			lv.Url = misc.UrlFix(lv.Url, row.ListUrl) // 修正链接

			// 根据Url，判断是否重复
			n, err := dbc.Count("Content", dbs.H{"Url": lv.Url})
			if err != nil {
				errorNum++
				cronErrorLog(0, "内容页查询重复入库失败: %v", err.Error())
				continue
			}
			if n > 0 {
				repeatNum++
				// cronLog(0, "重复Url: %v", Url)
				continue
			}

			// 抓取页面
			t2 := time.Now()
			bodyByte, i, err := misc.HttpGetRetry(lv.Url)
			if err != nil {
				errorNum++
				cronErrorLog(time.Since(t2), "抓取内容页失败, Error: %v, Url: %v", err.Error(), lv.Url)
				continue
			}
			cronLog(time.Since(t2), "抓取内容页完成, 请求次数: %v, Url: %v", i, lv.Url)

			// 解析代码
			t3 := time.Now()
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
			if err != nil {
				errorNum++
				cronErrorLog(0, "解析页面失败: %v, Url:%v", err.Error(), lv.Url)
				continue
			}

			// 匹配字段
			data := dbs.H{}
			data["Lid"] = lv.Lid
			data["Url"] = lv.Url
			for _, v := range *ParamContent {
				dom := doc.Find(v.Rule).Eq(0)
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
				value = misc.StrClear(value, v.FilterType, v.FilterRegexp, v.FilterRepl)
				data[v.Field] = value
			}

			// 写入数据库
			id, err := dbc.Insert("Content", data)
			if err != nil {
				errorNum++
				cronErrorLog(0, "内容页写入数据库失败: %v", err.Error())
				continue
			}
			newNum++
			cronLog(time.Since(t3), "内容页写入数据库成功: %v", id)
		}
	}
	cronLog(time.Since(t), "内容页下载资源完成, 总数: %v, 重复: %v, 新增: %v, 错误: %v", n, repeatNum, newNum, errorNum)
}

// 下载：分析列表页下载地址
func downInitList(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule) {
	for _, v := range *ParamList {
		if v.DownType == 1 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			where := dbs.H{}
			n, err := dbc.Count("List", where)
			if err != nil {
				cronErrorLog(0, "列表统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(100)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				st := ListUrl{}
				var list []ListUrl
				err := dbc.Find("List", "`Lid`,`"+v.Field+"`", []interface{}{&st.Lid, &st.Url}, where, "Lid ASC", page, pageSize, func() {
					list = append(list, st)
				})
				if err != nil {
					cronErrorLog(0, "列表读取失败: %v", err.Error())
					return
				}

				for _, lv := range list {
					// 效验链接
					if lv.Url == "" {
						errorNum++
						cronErrorLog(0, "下载 Url 为空")
						continue
					}
					lv.Url = misc.UrlFix(lv.Url, row.ListUrl) // 修正链接

					// 根据Url，判断是否重复
					n, err := dbc.Count("ListDownload", dbs.H{"Lid": lv.Lid, "Field": v.Field, "OldUrl": lv.Url})
					if err != nil {
						errorNum++
						cronErrorLog(0, "分析下载地址排重查询失败: %v", err.Error())
						continue
					}
					if n > 0 {
						repeatNum++
						// cronLog(0, "分析下载地址重复Lid: %v, Url: %v", lv.Lid, lv.Url)
						continue
					}

					// 存放路径
					NumDir := int64(math.Floor(float64(lv.Lid / 1000)))
					NewUrl := fmt.Sprintf("/upload/%s/list/%s/%03d/%d%s", row.Database, v.Field, NumDir, lv.Lid, filepath.Ext(lv.Url))

					// 写入数据库
					_, err = dbc.Insert("ListDownload", dbs.H{
						"Lid":    lv.Lid,
						"Status": 1,
						"Field":  v.Field,
						"OldUrl": lv.Url,
						"NewUrl": NewUrl,
					})
					if err != nil {
						errorNum++
						cronErrorLog(0, "列表页下载地址写入数据库失败: %v", err.Error())
						continue
					}
				}
			}

			cronLog(time.Since(t), "列表页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v", n, repeatNum, errorNum)
		} else if v.DownType == 2 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			where := dbs.H{}
			n, err := dbc.Count("List", where)
			if err != nil {
				cronErrorLog(0, "列表统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(100)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				st := ListHtml{}
				var list []ListHtml
				err = dbc.Find("List", "`Lid`,`"+v.Field+"`", []interface{}{&st.Lid, &st.Html}, where, "Lid ASC", page, pageSize, func() {
					list = append(list, st)
				})
				if err != nil {
					cronErrorLog(0, "列表读取失败: %v", err.Error())
					return
				}

				for _, lv := range list {
					// 分析图片
					dom, err := goquery.NewDocumentFromReader(strings.NewReader(lv.Html))
					if err != nil {
						errorNum++
						cronErrorLog(0, "解析代码失败: %v, Html:%v", err.Error(), lv.Html)
						continue
					}

					// 下载列表
					dom.Find(v.DownRule).Each(func(i int, s *goquery.Selection) {
						Url := ""
						if v.DownValueType == "Text" {
							Url = s.Text()
						} else if v.DownValueType == "Attr" {
							Url, _ = s.Attr(v.DownValueAttr)
						}
						if Url == "" {
							errorNum++
							cronErrorLog(0, "下载 Url 为空")
							return
						}
						Url = misc.UrlFix(Url, row.ListUrl) // 修正链接

						// 根据Url，判断是否重复
						n, err := dbc.Count("ListDownload", dbs.H{"Lid": lv.Lid, "Field": v.Field, "OldUrl": Url})
						if err != nil {
							errorNum++
							cronErrorLog(0, "分析下载地址排重查询失败: %v", err.Error())
							return
						}
						if n > 0 {
							repeatNum++
							// cronLog(0, "分析下载地址重复Lid: %v, Url: %v", lv.Lid, lv.Url)
							return
						}

						// 存放路径
						NumDir := int64(math.Floor(float64(lv.Lid / 1000)))
						Sort := i + 1
						NewUrl := fmt.Sprintf("/upload/%s/list/%s/%03d/%d/%02d%s", row.Database, v.Field, NumDir, lv.Lid, Sort, filepath.Ext(Url))

						// 写入数据库
						_, err = dbc.Insert("ListDownload", dbs.H{
							"Lid":    lv.Lid,
							"Status": 1,
							"Field":  v.Field,
							"OldUrl": Url,
							"NewUrl": NewUrl,
							"Sort":   i,
						})
						if err != nil {
							errorNum++
							cronErrorLog(0, "列表页下载地址写入数据库失败: %v", err.Error())
							return
						}
					})

				}
			}

			cronLog(time.Since(t), "列表页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v", n, repeatNum, errorNum)
		}
	}
}

// 下载：分析内容页下载地址
func downInitContent(dbc *dbs.DB, ParamContent *[]RuleParam, row *Rule) {
	for _, v := range *ParamContent {
		if v.DownType == 1 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			where := dbs.H{}
			n, err := dbc.Count("Content", where)
			if err != nil {
				cronErrorLog(0, "内容统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(100)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				st := ListUrl{}
				var list []ListUrl
				err = dbc.Find("Content", "`Lid`,`"+v.Field+"`", []interface{}{&st.Lid, &st.Url}, where, "Lid ASC", page, pageSize, func() {
					list = append(list, st)
				})
				if err != nil {
					cronErrorLog(0, "内容读取失败: %v", err.Error())
					return
				}

				for _, lv := range list {
					// 效验链接
					if lv.Url == "" {
						errorNum++
						cronErrorLog(0, "下载 Url 为空")
						continue
					}
					lv.Url = misc.UrlFix(lv.Url, row.ContentUrl) // 修正链接

					// 根据Url，判断是否重复
					n, err := dbc.Count("ContentDownload", dbs.H{"Lid": lv.Lid, "Field": v.Field, "OldUrl": lv.Url})
					if err != nil {
						errorNum++
						cronErrorLog(0, "分析下载地址排重查询失败: %v", err.Error())
						continue
					}
					if n > 0 {
						repeatNum++
						// cronLog(0, "分析下载地址重复Lid: %v, Url: %v", lv.Lid, lv.Url)
						continue
					}

					// 存放路径
					NumDir := int64(math.Floor(float64(lv.Lid / 1000)))
					NewUrl := fmt.Sprintf("/upload/%s/list/%s/%03d/%d%s", row.Database, v.Field, NumDir, lv.Lid, filepath.Ext(lv.Url))

					// 写入数据库
					_, err = dbc.Insert("ContentDownload", dbs.H{
						"Lid":    lv.Lid,
						"Status": 1,
						"Field":  v.Field,
						"OldUrl": lv.Url,
						"NewUrl": NewUrl,
					})
					if err != nil {
						errorNum++
						cronErrorLog(0, "内容页下载地址写入数据库失败: %v", err.Error())
						continue
					}
				}
			}

			cronLog(time.Since(t), "内容页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v", n, repeatNum, errorNum)
		} else if v.DownType == 2 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			where := dbs.H{}
			n, err := dbc.Count("Content", where)
			if err != nil {
				cronErrorLog(0, "内容统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(100)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				st := ListHtml{}
				var list []ListHtml
				err = dbc.Find("Content", "`Lid`,`"+v.Field+"`", []interface{}{&st.Lid, &st.Html}, where, "Lid ASC", page, pageSize, func() {
					list = append(list, st)
				})
				if err != nil {
					cronErrorLog(0, "内容读取失败: %v", err.Error())
					return
				}

				for _, lv := range list {
					// 分析图片
					dom, err := goquery.NewDocumentFromReader(strings.NewReader(lv.Html))
					if err != nil {
						errorNum++
						cronErrorLog(0, "解析代码失败: %v, Html:%v", err.Error(), lv.Html)
						continue
					}

					// 下载内容
					dom.Find(v.DownRule).Each(func(i int, s *goquery.Selection) {
						Url := ""
						if v.DownValueType == "Text" {
							Url = s.Text()
						} else if v.DownValueType == "Attr" {
							Url, _ = s.Attr(v.DownValueAttr)
						}
						if Url == "" {
							errorNum++
							cronErrorLog(0, "下载 Url 为空")
							return
						}
						Url = misc.UrlFix(Url, row.ContentUrl) // 修正链接

						// 根据Url，判断是否重复
						n, err := dbc.Count("ContentDownload", dbs.H{"Lid": lv.Lid, "Field": v.Field, "OldUrl": Url})
						if err != nil {
							errorNum++
							cronErrorLog(0, "分析下载地址排重查询失败: %v", err.Error())
							return
						}
						if n > 0 {
							repeatNum++
							// cronLog(0, "分析下载地址重复Lid: %v, Url: %v", lv.Lid, lv.Url)
							return
						}

						// 存放路径
						NumDir := int64(math.Floor(float64(lv.Lid / 1000)))
						Sort := i + 1
						NewUrl := fmt.Sprintf("/upload/%s/list/%s/%03d/%d/%02d%s", row.Database, v.Field, NumDir, lv.Lid, Sort, filepath.Ext(Url))

						// 写入数据库
						_, err = dbc.Insert("ContentDownload", dbs.H{
							"Lid":    lv.Lid,
							"Status": 1,
							"Field":  v.Field,
							"OldUrl": Url,
							"NewUrl": NewUrl,
							"Sort":   Sort,
						})
						if err != nil {
							errorNum++
							cronErrorLog(0, "内容页下载地址写入数据库失败: %v", err.Error())
							return
						}
					})

				}
			}

			cronLog(time.Since(t), "内容页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v", n, repeatNum, errorNum)
		}
	}
}

// 下载：列表页资源下载
func downList(dbc *dbs.DB, row *Rule) {
	t := time.Now() // 记时开始
	errorNum := 0   // 错误入库数

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("ListDownload", where)
	if err != nil {
		cronErrorLog(0, "下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(100)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		st := ListDown{}
		var list []ListDown
		fields := "`Id`,`Lid`,`Field`,`OldUrl`,`NewUrl`,`Sort`"
		scanArr := []interface{}{&st.Id, &st.Lid, &st.Field, &st.OldUrl, &st.NewUrl, &st.Sort}
		err = dbc.Find("ListDownload", fields, scanArr, where, "Id ASC", 0, pageSize, func() {
			list = append(list, st)
		})
		if err != nil {
			cronErrorLog(0, "下载列表读取失败: %v", err.Error())
			return
		}

		for _, lv := range list {
			// 下载文件
			t2 := time.Now() // 记时开始
			FileSize, err := misc.DownloadFile(lv.OldUrl, "."+lv.NewUrl)
			if err != nil {
				errorNum++
				cronErrorLog(time.Since(t2), "下载文件失败: %v, 大小: %v, File: %v, Url: %v", err.Error(), FileSize, lv.NewUrl, lv.OldUrl)

				// 更新状态为下载失败
				_, err = dbc.Update("ListDownload", dbs.H{"Status": 3}, dbs.H{"Id": lv.Id})
				if err != nil {
					errorNum++
					cronErrorLog(0, "更新数据库失败: %v", err.Error())
				}
				continue
			}
			cronLog(time.Since(t2), "下载文件完成, 大小: %v, File: %v, Url: %v", FileSize, lv.NewUrl, lv.OldUrl)

			// 更新
			_, err = dbc.Update("ListDownload", dbs.H{
				"Status":       2, // 下载完成
				"FileSize":     FileSize,
				"DownloadDate": time.Now().Format("2006-01-02 15:04:05"),
			}, dbs.H{
				"Id": lv.Id,
			})
			if err != nil {
				errorNum++
				cronErrorLog(0, "下载完成更新数据库失败: %v", err.Error())
			}
		}
	}

	cronLog(time.Since(t), "列表页下载资源完成, 总数: %v, 错误: %v", n, errorNum)
}

// 下载：内容页资源下载
func downContent(dbc *dbs.DB, row *Rule) {
	t := time.Now() // 记时开始
	errorNum := 0   // 错误入库数

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("ContentDownload", where)
	if err != nil {
		cronErrorLog(0, "下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(100)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		st := ListDown{}
		var list []ListDown
		fields := "`Id`,`Lid`,`Field`,`OldUrl`,`NewUrl`,`Sort`"
		scanArr := []interface{}{&st.Id, &st.Lid, &st.Field, &st.OldUrl, &st.NewUrl, &st.Sort}
		err := dbc.Find("ContentDownload", fields, scanArr, where, "Id ASC", 0, pageSize, func() {
			list = append(list, st)
		})
		if err != nil {
			cronErrorLog(0, "下载列表读取失败: %v", err.Error())
			return
		}

		for _, lv := range list {
			// 下载文件
			t2 := time.Now() // 记时开始
			FileSize, err := misc.DownloadFile(lv.OldUrl, "."+lv.NewUrl)
			if err != nil {
				errorNum++
				cronErrorLog(time.Since(t2), "下载文件失败: %v, 大小: %v, File: %v, Url: %v", err.Error(), FileSize, lv.NewUrl, lv.OldUrl)

				// 更新状态为下载失败
				_, err = dbc.Update("ContentDownload", dbs.H{"Status": 3}, dbs.H{"Id": lv.Id})
				if err != nil {
					errorNum++
					cronErrorLog(0, "更新数据库失败: %v", err.Error())
				}
				continue
			}
			cronLog(time.Since(t2), "下载文件完成, 大小: %v, File: %v, Url: %v", FileSize, lv.NewUrl, lv.OldUrl)

			// 更新
			_, err = dbc.Update("ContentDownload", dbs.H{
				"Status":       2, // 下载完成
				"FileSize":     FileSize,
				"DownloadDate": time.Now().Format("2006-01-02 15:04:05"),
			}, dbs.H{
				"Id": lv.Id,
			})
			if err != nil {
				errorNum++
				cronErrorLog(0, "下载完成更新数据库失败: %v", err.Error())
			}
		}
	}

	cronLog(time.Since(t), "内容页下载资源完成, 总数: %v, 错误: %v", n, errorNum)
}
