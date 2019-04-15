package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/xiuno/gin"

	"./control/api"
	"./control/front"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func StartGin() {
	WebLogFile := "./log/web.log"
	WebErrorLogFile := "./log/webError.log"

	// 日志
	webLog, err := os.OpenFile(WebLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(webLog, os.Stdout)

	// 错误日志
	webErrorLog, err := os.OpenFile(WebErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	gin.DefaultErrorWriter = io.MultiWriter(webErrorLog, os.Stdout)

	r := gin.Default()
	r.Static("/static", "./static")

	// 附件目录，登录后才能查看
	upload := r.Group("/", func(c *gin.Context) {
		err := api.UserTokenGet(c)
		if err != nil {
			// c.Writer.Header().Add("Author", "goSpider 1.0")
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
			return
		}
	})
	upload.Static("/upload", "./upload")

	// 模板
	r.SetFuncMap(template.FuncMap{
		"htmlTags": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	r.LoadHTMLGlob("./view/*")

	// ========== 不用登录 ==========
	r.GET("/Login", front.Login)
	r.GET("/Logout", front.Logout)
	r.POST("/UserLogin", api.UserLogin)

	// ========== 前台页面 ==========
	app := r.Group("/", func(c *gin.Context) {
		err := api.UserTokenGet(c)
		if err != nil {
			c.Redirect(http.StatusFound, "/Login")
			return
		}
	})
	app.GET("/", front.Index)
	app.GET("/Log", front.Log)
	app.GET("/User", front.User)
	app.GET("/Sys", front.Sys)
	app.GET("/Item", front.Item)
	app.GET("/Show", front.Show)

	// ========== 后台接口 ==========
	admin := r.Group("/", func(c *gin.Context) {
		err := api.UserTokenGet(c)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	})

	// User
	admin.POST("/UserCreate", api.UserCreate)
	admin.POST("/UserUpdate", api.UserUpdate)
	admin.POST("/UserRead", api.UserRead)
	admin.POST("/UserDelete", api.UserDelete)
	admin.POST("/UserList", api.UserList)

	// Rule
	admin.POST("/RuleCreate", api.RuleCreate)
	admin.POST("/RuleUpdate", api.RuleUpdate)
	admin.POST("/RuleUpdateCron", api.RuleUpdateCron)
	admin.POST("/RuleRead", api.RuleRead)
	admin.POST("/RuleDelete", api.RuleDelete)
	admin.POST("/RuleList", api.RuleList)
	admin.POST("/RuleExport", api.RuleExport)
	admin.POST("/RuleImport", api.RuleImport)

	// RuleParam
	admin.POST("/RuleParamCreate", api.RuleParamCreate)
	admin.POST("/RuleParamUpdate", api.RuleParamUpdate)
	admin.POST("/RuleParamRead", api.RuleParamRead)
	admin.POST("/RuleParamDelete", api.RuleParamDelete)
	admin.POST("/RuleParamList", api.RuleParamList)

	// HttpGet
	admin.POST("/HttpGetPage", api.HttpGetPage)
	admin.POST("/HttpGetList", api.HttpGetList)
	admin.POST("/HttpGetListRule", api.HttpGetListRule)
	admin.POST("/HttpGetContentRule", api.HttpGetContentRule)
	admin.POST("/HttpGetListRuleDown", api.HttpGetListRuleDown)
	admin.POST("/HttpGetContentRuleDown", api.HttpGetContentRuleDown)

	// Database
	admin.POST("/DatabaseCreate", api.DatabaseCreate)

	// Item
	admin.POST("/ItemList", api.ItemList)

	// Show
	admin.POST("/ShowList", api.ShowList)

	// Log
	admin.POST("/LogList", api.LogList)
	admin.POST("/LogDeleteDB", api.LogDeleteDB)

	// 内存信息 && 磁盘信息
	admin.POST("/MemStatsInfo", api.MemStatsInfo)
	admin.POST("/DiskInfo", api.DiskInfo)

	err = r.Run("0.0.0.0:3333")
	if err != nil {
		panic(err)
	}
}

func main() {
	ConfigRuntime()
	StartGin()
}
