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

	// OSS 文件
	r.GET("/file/*filepath", api.OssFile)

	// 附件目录，登录后才能查看
	upload := r.Group("/", func(c *gin.Context) {
		_, err := api.UserTokenGet(c)
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

	// ==============================================================================================================
	// ========== 所以用户页面 ==========
	app := r.Group("/", func(c *gin.Context) {
		_, err := api.UserTokenGet(c)
		if err != nil {
			c.Redirect(http.StatusFound, "/Login")
			return
		}
	})
	app.GET("/", front.Index)
	app.GET("/Show", front.Show)
	app.GET("/Read", front.Read)

	// ========== 所以用户接口 ==========
	appApi := r.Group("/", func(c *gin.Context) {
		_, err := api.UserTokenGet(c)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	})
	appApi.POST("/UserRuleList", api.RuleList)
	appApi.POST("/ShowRead", api.ShowRead)
	appApi.POST("/ShowList", api.ShowList)
	appApi.POST("/ShowDownload", api.ShowDownload)

	// ==============================================================================================================
	// ========== 管理员页面 ==========
	admin := r.Group("/", func(c *gin.Context) {
		_, err := api.UserTokenGetByAdmin(c)
		if err != nil {
			c.Redirect(http.StatusFound, "/Login")
			return
		}
	})
	admin.GET("/Log", front.Log)
	admin.GET("/User", front.User)
	admin.GET("/Oss", front.Oss)
	admin.GET("/Sys", front.Sys)
	admin.GET("/Setting", front.Setting)
	admin.GET("/Table", front.Table)

	// ========== 管理员接口 ==========
	adminApi := r.Group("/", func(c *gin.Context) {
		_, err := api.UserTokenGetByAdmin(c)
		if err != nil {
			c.Message("-1", err.Error())
			return
		}
	})

	// User
	adminApi.POST("/UserCreate", api.UserCreate)
	adminApi.POST("/UserUpdate", api.UserUpdate)
	adminApi.POST("/UserRead", api.UserRead)
	adminApi.POST("/UserDelete", api.UserDelete)
	adminApi.POST("/UserList", api.UserList)

	// Rule
	adminApi.POST("/RuleCreate", api.RuleCreate)
	adminApi.POST("/RuleUpdate", api.RuleUpdate)
	adminApi.POST("/RuleUpdateCron", api.RuleUpdateCron)
	adminApi.POST("/RuleRead", api.RuleRead)
	adminApi.POST("/RuleDelete", api.RuleDelete)
	adminApi.POST("/RuleList", api.RuleList)
	adminApi.POST("/RuleExport", api.RuleExport)
	adminApi.POST("/RuleImport", api.RuleImport)

	// RuleParam
	adminApi.POST("/RuleParamCreate", api.RuleParamCreate)
	adminApi.POST("/RuleParamUpdate", api.RuleParamUpdate)
	adminApi.POST("/RuleParamRead", api.RuleParamRead)
	adminApi.POST("/RuleParamDelete", api.RuleParamDelete)
	adminApi.POST("/RuleParamList", api.RuleParamList)

	// HttpGet
	adminApi.POST("/HttpGetPage", api.HttpGetPage)
	adminApi.POST("/HttpGetList", api.HttpGetList)
	adminApi.POST("/HttpGetListRule", api.HttpGetListRule)
	adminApi.POST("/HttpGetContentRule", api.HttpGetContentRule)
	adminApi.POST("/HttpGetListRuleDown", api.HttpGetListRuleDown)
	adminApi.POST("/HttpGetContentRuleDown", api.HttpGetContentRuleDown)

	// Database
	adminApi.POST("/DatabaseCreate", api.DatabaseCreate)

	// Table
	adminApi.POST("/TableList", api.TableList)

	// Log
	adminApi.POST("/LogList", api.LogList)
	adminApi.POST("/LogDeleteDB", api.LogDeleteDB)

	// Setting
	adminApi.POST("/SettingSave", api.SettingSave)
	adminApi.POST("/SettingSet", api.SettingSet)
	adminApi.POST("/SettingGet", api.SettingGet)
	adminApi.POST("/SettingList", api.SettingList)
	adminApi.POST("/SettingInit", api.SettingInit)

	// Oss
	adminApi.POST("/OssBucketList", api.OssBucketList)
	adminApi.POST("/OssObjectList", api.OssObjectList)

	// 内存信息 && 磁盘信息
	adminApi.POST("/MemStatsInfo", api.MemStatsInfo)
	adminApi.POST("/DiskInfo", api.DiskInfo)

	err = r.Run("0.0.0.0:3333")
	if err != nil {
		panic(err)
	}
}

func main() {
	ConfigRuntime()
	StartGin()
}
