package main

import (
	"fmt"
	"html/template"
	"io"
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

	app := gin.Default()
	app.Delims("${", "}")
	app.Static("/upload", "./upload")
	app.SetFuncMap(template.FuncMap{
		"htmlTags": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	app.LoadHTMLGlob("./view/*")

	// 内存信息 && 磁盘信息
	app.GET("/MemStatsInfo", api.MemStatsInfo)
	app.GET("/DiskInfo", api.DiskInfo)

	// ========== Front ==========
	app.GET("/", front.Index)
	app.GET("/Item", front.Item)
	app.GET("/User", front.User)
	app.GET("/Log", front.Log)
	app.GET("/Sys", front.Sys)

	// ========== Api ==========
	// User
	app.POST("/UserLogin", api.UserLogin)
	app.POST("/UserCreate", api.UserCreate)
	app.POST("/UserUpdate", api.UserUpdate)
	app.POST("/UserRead", api.UserRead)
	app.POST("/UserDelete", api.UserDelete)
	app.POST("/UserList", api.UserList)

	// Rule
	app.POST("/RuleCreate", api.RuleCreate)
	app.POST("/RuleUpdate", api.RuleUpdate)
	app.POST("/RuleUpdateCron", api.RuleUpdateCron)
	app.POST("/RuleRead", api.RuleRead)
	app.POST("/RuleDelete", api.RuleDelete)
	app.POST("/RuleList", api.RuleList)
	app.POST("/RuleExport", api.RuleExport)
	app.POST("/RuleImport", api.RuleImport)

	// RuleParam
	app.POST("/RuleParamCreate", api.RuleParamCreate)
	app.POST("/RuleParamUpdate", api.RuleParamUpdate)
	app.POST("/RuleParamRead", api.RuleParamRead)
	app.POST("/RuleParamDelete", api.RuleParamDelete)
	app.POST("/RuleParamList", api.RuleParamList)

	// HttpGet
	app.POST("/HttpGetPage", api.HttpGetPage)
	app.POST("/HttpGetList", api.HttpGetList)
	app.POST("/HttpGetListRule", api.HttpGetListRule)
	app.POST("/HttpGetContentRule", api.HttpGetContentRule)
	app.POST("/HttpGetListRuleDown", api.HttpGetListRuleDown)
	app.POST("/HttpGetContentRuleDown", api.HttpGetContentRuleDown)

	// Database
	app.POST("/DatabaseCreate", api.DatabaseCreate)

	// Item
	app.POST("/ItemList", api.ItemList)

	// Log
	app.POST("/LogList", api.LogList)
	app.POST("/LogDeleteDB", api.LogDeleteDB)

	err = app.Run("0.0.0.0:3333")
	if err != nil {
		panic(err)
	}
}

func main() {
	ConfigRuntime()
	StartGin()
}
