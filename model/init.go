package model

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
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
			ParamList, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "List"}, 0, 0)
			if err != nil {
				cronErrorLog("读取列表采集参数失败: %v", err.Error())
				continue
			}
			ParamContent, err := RuleParamList(dbs.H{"Rid": row.Rid, "Type": "Content"}, 0, 0)
			if err != nil {
				cronErrorLog("读取内容采集参数失败: %v", err.Error())
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
				cronErrorLog("列表页查询重复入库失败: %v", err.Error())
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
			cronErrorLog("列表页写入数据库失败: %v", err.Error())
			return
		}
		cronLog("列表页写入数据库成功: %v", id)
	})
	cronLog("Url: %v 入库完成, 耗时： %v", Url, time.Since(t2))
}

// 内容页：抓取内容页
func getContent(dbc *dbs.DB, ParamContent *[]RuleParam, row *Rule) {
	t := time.Now() // 记时开始
	repeatNum := 0  // 重复入库数
	errorNum := 0   // 错误入库数

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("List", where)
	if err != nil {
		cronErrorLog("下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(1000)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		rows, err := dbc.Find("List", "Url", where, "ListId DESC", page, pageSize)
		if err != nil {
			cronErrorLog("列表读取失败: %v", err.Error())
			return
		}
		var list []string
		for rows.Next() {
			var Url string
			err = rows.Scan(&Url)
			if err != nil {
				cronErrorLog("Url绑定失败: %v", err.Error())
				return
			}
			list = append(list, Url)
		}

		for _, Url := range list {
			// 效验链接
			if Url == "" {
				errorNum++
				cronErrorLog("内容页 Url 为空")
				continue
			}
			Url = misc.UrlFix(Url, row.ListUrl) // 修正链接

			// 根据Url，判断是否重复
			n, err := dbc.Count("Content", dbs.H{"Url": Url})
			if err != nil {
				errorNum++
				cronErrorLog("内容页查询重复入库失败: %v", err.Error())
				continue
			}
			if n > 0 {
				repeatNum++
				// cronLog("重复Url: %v", Url)
				continue
			}

			// 抓取页面
			t2 := time.Now()
			bodyByte, i, err := misc.HttpGetRetry(Url)
			if err != nil {
				errorNum++
				cronErrorLog("抓取页面失败, 耗时: %v, Error: %v", time.Since(t2), err.Error())
				continue
			}
			cronLog("抓取页面完成, 请求次数: %v, 耗时: %v, Url: %v", i, time.Since(t2), Url)

			// 解析代码
			t3 := time.Now()
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyByte)))
			if err != nil {
				errorNum++
				cronErrorLog("解析页面失败: %v, Url:%v", err.Error(), Url)
				continue
			}

			// 匹配内容
			data := dbs.H{}
			data["Url"] = Url
			for _, m := range *ParamContent {
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

			// 写入数据库
			id, err := dbc.Insert("Content", data)
			if err != nil {
				errorNum++
				cronErrorLog("内容页写入数据库失败: %v", err.Error())
				continue
			}
			cronLog("内容页写入数据库成功: %v, 耗时: %v", id, time.Since(t3))
		}
	}
	cronLog("内容页下载资源完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
}

// 下载：分析列表页下载地址
func downInitList(dbc *dbs.DB, ParamList *[]RuleParam, row *Rule) {
	for _, v := range *ParamList {
		if v.DownType == 1 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			type st struct {
				ListId int64
				Url    string
			}

			where := dbs.H{}
			n, err := dbc.Count("List", where)
			if err != nil {
				cronErrorLog("列表统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(1000)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				rows, err := dbc.Find("List", "`ListId`,`"+v.Field+"`", where, "ListId ASC", page, pageSize)
				if err != nil {
					cronErrorLog("列表读取失败: %v", err.Error())
					return
				}
				var list []st
				for rows.Next() {
					st := st{}
					err = rows.Scan(&st.ListId, &st.Url)
					if err != nil {
						cronErrorLog("绑定失败: %v", err.Error())
						continue
					}
					list = append(list, st)
				}

				for _, lv := range list {
					// 效验链接
					if lv.Url == "" {
						errorNum++
						cronErrorLog("下载 Url 为空")
						continue
					}
					lv.Url = misc.UrlFix(lv.Url, row.ListUrl) // 修正链接

					// 根据Url，判断是否重复
					n, err := dbc.Count("ListDownload", dbs.H{"ListId": lv.ListId, "Field": v.Field, "OldUrl": lv.Url})
					if err != nil {
						errorNum++
						cronErrorLog("分析下载地址排重查询失败: %v", err.Error())
						continue
					}
					if n > 0 {
						repeatNum++
						// cronLog("分析下载地址重复ListId: %v, Url: %v", lv.ListId, lv.Url)
						continue
					}

					// 写入数据库
					_, err = dbc.Insert("ListDownload", dbs.H{
						"ListId": lv.ListId,
						"Status": 1,
						"Field":  v.Field,
						"OldUrl": lv.Url,
					})
					if err != nil {
						errorNum++
						cronErrorLog("列表页下载地址写入数据库失败: %v", err.Error())
						continue
					}
				}
			}

			cronLog("列表页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
		} else if v.DownType == 2 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			type st struct {
				ListId int64
				Html   string
			}

			where := dbs.H{}
			n, err := dbc.Count("List", where)
			if err != nil {
				cronErrorLog("列表统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(1000)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				rows, err := dbc.Find("List", "`ListId`,`"+v.Field+"`", where, "ListId ASC", page, pageSize)
				if err != nil {
					cronErrorLog("列表读取失败: %v", err.Error())
					return
				}
				var list []st
				for rows.Next() {
					st := st{}
					err = rows.Scan(&st.ListId, &st.Html)
					if err != nil {
						cronErrorLog("绑定失败: %v", err.Error())
						continue
					}
					list = append(list, st)
				}

				for _, lv := range list {
					// 分析图片
					dom, err := goquery.NewDocumentFromReader(strings.NewReader(lv.Html))
					if err != nil {
						errorNum++
						cronErrorLog("解析代码失败: %v, Html:%v", err.Error(), lv.Html)
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
							cronErrorLog("下载 Url 为空")
							return
						}
						Url = misc.UrlFix(Url, row.ListUrl) // 修正链接

						// 根据Url，判断是否重复
						n, err := dbc.Count("ListDownload", dbs.H{"ListId": lv.ListId, "Field": v.Field, "OldUrl": Url})
						if err != nil {
							errorNum++
							cronErrorLog("分析下载地址排重查询失败: %v", err.Error())
							return
						}
						if n > 0 {
							repeatNum++
							// cronLog("分析下载地址重复ListId: %v, Url: %v", lv.ListId, lv.Url)
							return
						}

						// 写入数据库
						_, err = dbc.Insert("ListDownload", dbs.H{
							"ListId": lv.ListId,
							"Status": 1,
							"Field":  v.Field,
							"OldUrl": Url,
							"Sort":   i,
						})
						if err != nil {
							errorNum++
							cronErrorLog("列表页下载地址写入数据库失败: %v", err.Error())
							return
						}
					})

				}
			}

			cronLog("列表页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
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

			type st struct {
				ContentId int64
				Url       string
			}

			where := dbs.H{}
			n, err := dbc.Count("Content", where)
			if err != nil {
				cronErrorLog("内容统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(1000)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				rows, err := dbc.Find("Content", "`ContentId`,`"+v.Field+"`", where, "ContentId ASC", page, pageSize)
				if err != nil {
					cronErrorLog("内容读取失败: %v", err.Error())
					return
				}
				var list []st
				for rows.Next() {
					st := st{}
					err = rows.Scan(&st.ContentId, &st.Url)
					if err != nil {
						cronErrorLog("绑定失败: %v", err.Error())
						continue
					}
					list = append(list, st)
				}

				for _, lv := range list {
					// 效验链接
					if lv.Url == "" {
						errorNum++
						cronErrorLog("下载 Url 为空")
						continue
					}
					lv.Url = misc.UrlFix(lv.Url, row.ContentUrl) // 修正链接

					// 根据Url，判断是否重复
					n, err := dbc.Count("ContentDownload", dbs.H{"ContentId": lv.ContentId, "Field": v.Field, "OldUrl": lv.Url})
					if err != nil {
						errorNum++
						cronErrorLog("分析下载地址排重查询失败: %v", err.Error())
						continue
					}
					if n > 0 {
						repeatNum++
						// cronLog("分析下载地址重复ContentId: %v, Url: %v", lv.ContentId, lv.Url)
						continue
					}

					// 写入数据库
					_, err = dbc.Insert("ContentDownload", dbs.H{
						"ContentId": lv.ContentId,
						"Status":    1,
						"Field":     v.Field,
						"OldUrl":    lv.Url,
					})
					if err != nil {
						errorNum++
						cronErrorLog("内容页下载地址写入数据库失败: %v", err.Error())
						continue
					}
				}
			}

			cronLog("内容页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
		} else if v.DownType == 2 {
			t := time.Now() // 记时开始
			repeatNum := 0  // 重复入库数
			errorNum := 0   // 错误入库数

			type st struct {
				ContentId int64
				Html      string
			}

			where := dbs.H{}
			n, err := dbc.Count("Content", where)
			if err != nil {
				cronErrorLog("内容统计数量失败: %v", err.Error())
				return
			}

			pageSize := int64(1000)
			pageMax := int64(math.Ceil(float64(n / pageSize)))

			for page := int64(1); page <= pageMax; page++ {
				rows, err := dbc.Find("Content", "`ContentId`,`"+v.Field+"`", where, "ContentId ASC", page, pageSize)
				if err != nil {
					cronErrorLog("内容读取失败: %v", err.Error())
					return
				}
				var list []st
				for rows.Next() {
					st := st{}
					err = rows.Scan(&st.ContentId, &st.Html)
					if err != nil {
						cronErrorLog("绑定失败: %v", err.Error())
						continue
					}
					list = append(list, st)
				}

				for _, lv := range list {
					// 分析图片
					dom, err := goquery.NewDocumentFromReader(strings.NewReader(lv.Html))
					if err != nil {
						errorNum++
						cronErrorLog("解析代码失败: %v, Html:%v", err.Error(), lv.Html)
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
							cronErrorLog("下载 Url 为空")
							return
						}
						Url = misc.UrlFix(Url, row.ContentUrl) // 修正链接

						// 根据Url，判断是否重复
						n, err := dbc.Count("ContentDownload", dbs.H{"ContentId": lv.ContentId, "Field": v.Field, "OldUrl": Url})
						if err != nil {
							errorNum++
							cronErrorLog("分析下载地址排重查询失败: %v", err.Error())
							return
						}
						if n > 0 {
							repeatNum++
							// cronLog("分析下载地址重复ContentId: %v, Url: %v", lv.ContentId, lv.Url)
							return
						}

						// 写入数据库
						_, err = dbc.Insert("ContentDownload", dbs.H{
							"ContentId": lv.ContentId,
							"Status":    1,
							"Field":     v.Field,
							"OldUrl":    Url,
							"Sort":      i,
						})
						if err != nil {
							errorNum++
							cronErrorLog("内容页下载地址写入数据库失败: %v", err.Error())
							return
						}
					})

				}
			}

			cronLog("内容页下载地址入库完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
		}
	}
}

// 下载：列表页资源下载
func downList(dbc *dbs.DB, row *Rule) {
	t := time.Now() // 记时开始
	repeatNum := 0  // 重复入库数
	errorNum := 0   // 错误入库数

	fields := "`Id`,`ListId`,`Field`,`OldUrl`,`NewUrl`,`Sort`"
	type st struct {
		Id     int64
		ListId int64
		Field  string
		OldUrl string
		NewUrl string
		Sort   int64
	}

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("ListDownload", where)
	if err != nil {
		cronErrorLog("下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(1000)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		rows, err := dbc.Find("ListDownload", fields, where, "Id ASC", 0, pageSize)
		if err != nil {
			cronErrorLog("下载列表读取失败: %v", err.Error())
			return
		}

		var list []st
		for rows.Next() {
			st := st{}
			err = rows.Scan(&st.Id, &st.ListId, &st.Field, &st.OldUrl, &st.NewUrl, &st.Sort)
			if err != nil {
				cronErrorLog("绑定失败: %v", err.Error())
				return
			}
			list = append(list, st)
		}

		for _, lv := range list {
			// 存放目录
			path := fmt.Sprintf("/upload/%s/list/%s/%03d/%d/", row.DateBase, lv.Field,
				int64(math.Floor(float64(lv.ListId/1000))), lv.ListId)
			dir := "." + path
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				cronErrorLog("创建目录失败: %v", err.Error())
				return
			}

			// 存放文件名
			filename := fmt.Sprintf("%02d%v", lv.Sort+1, filepath.Ext(lv.OldUrl))

			// 统计已下载过的数量
			if lv.NewUrl != "" {
				repeatNum++
			}

			// 下载文件
			t2 := time.Now() // 记时开始
			FileSize, err := misc.DownloadFile(lv.OldUrl, dir+filename)
			if err != nil {
				errorNum++
				cronErrorLog("下载文件失败: %v, 大小: %v, 耗时: %v, File: %v, Url: %v", err.Error(), FileSize, time.Since(t2), dir+filename, lv.OldUrl)
				continue
			}
			cronLog("下载文件完成, 大小: %v, 耗时: %v, File: %v, Url: %v", FileSize, time.Since(t2), dir+filename, lv.OldUrl)

			// 更新
			_, err = dbc.Update("ListDownload", dbs.H{
				"Status":       2, // 下载完成
				"NewUrl":       path + filename,
				"FileSize":     FileSize,
				"DownloadDate": time.Now().Format("2006-01-02 15:04:05"),
			}, dbs.H{
				"Id": lv.Id,
			})
			if err != nil {
				errorNum++
				cronErrorLog("下载完成更新数据库失败: %v", err.Error())
			}
		}
	}

	cronLog("列表页下载资源完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
}

// 下载：内容页资源下载
func downContent(dbc *dbs.DB, row *Rule) {
	t := time.Now() // 记时开始
	repeatNum := 0  // 重复入库数
	errorNum := 0   // 错误入库数

	fields := "`Id`,`ContentId`,`Field`,`OldUrl`,`NewUrl`,`Sort`"
	type st struct {
		Id        int64
		ContentId int64
		Field     string
		OldUrl    string
		NewUrl    string
		Sort      int64
	}

	where := dbs.H{"Status": 1}
	n, err := dbc.Count("ContentDownload", where)
	if err != nil {
		cronErrorLog("下载列表统计数量失败: %v", err.Error())
		return
	}

	pageSize := int64(1000)
	pageMax := int64(math.Ceil(float64(n / pageSize)))

	for page := int64(1); page <= pageMax; page++ {
		rows, err := dbc.Find("ContentDownload", fields, where, "Id ASC", 0, pageSize)
		if err != nil {
			cronErrorLog("下载列表读取失败: %v", err.Error())
			return
		}

		var list []st
		for rows.Next() {
			st := st{}
			err = rows.Scan(&st.Id, &st.ContentId, &st.Field, &st.OldUrl, &st.NewUrl, &st.Sort)
			if err != nil {
				cronErrorLog("绑定失败: %v", err.Error())
				return
			}
			list = append(list, st)
		}

		for _, lv := range list {
			// 存放目录
			path := fmt.Sprintf("/upload/%s/content/%s/%03d/%d/", row.DateBase, lv.Field,
				int64(math.Floor(float64(lv.ContentId/1000))), lv.ContentId)
			dir := "." + path
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				cronErrorLog("创建目录失败: %v", err.Error())
				return
			}

			// 存放文件名
			filename := fmt.Sprintf("%02d%v", lv.Sort+1, filepath.Ext(lv.OldUrl))

			// 统计已下载过的数量
			if lv.NewUrl != "" {
				repeatNum++
			}

			// 下载文件
			t2 := time.Now() // 记时开始
			FileSize, err := misc.DownloadFile(lv.OldUrl, dir+filename)
			if err != nil {
				errorNum++
				cronErrorLog("下载文件失败: %v, 大小: %v, 耗时: %v, File: %v, Url: %v", err.Error(), FileSize, time.Since(t2), dir+filename, lv.OldUrl)
				continue
			}
			cronLog("下载文件完成, 大小: %v, 耗时: %v, File: %v, Url: %v", FileSize, time.Since(t2), dir+filename, lv.OldUrl)

			// 更新
			_, err = dbc.Update("ContentDownload", dbs.H{
				"Status":       2, // 下载完成
				"NewUrl":       path + filename,
				"FileSize":     FileSize,
				"DownloadDate": time.Now().Format("2006-01-02 15:04:05"),
			}, dbs.H{
				"Id": lv.Id,
			})
			if err != nil {
				errorNum++
				cronErrorLog("下载完成更新数据库失败: %v", err.Error())
			}
		}
	}

	cronLog("内容页下载资源完成, 总数: %v, 重复: %v, 错误: %v, 耗时: %v", n, repeatNum, errorNum, time.Since(t))
}
