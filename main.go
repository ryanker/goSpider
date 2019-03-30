package main

import (
	"html/template"
	"io"
	"os"
	"runtime"

	"github.com/xiuno/gin"

	"./api"
)

func main() {
	runtime.GOMAXPROCS(3)

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
	app.Delims("{[{", "}]}")
	app.Static("/upload", "./upload")
	app.SetFuncMap(template.FuncMap{
		"htmlTags": func(s string) template.HTML {
			return template.HTML(s)
		},
	})
	app.LoadHTMLGlob("./view/*")

	// Rule
	app.POST("/RuleCreate", api.RuleCreate)
	app.POST("/RuleUpdate", api.RuleUpdate)
	app.POST("/RuleRead", api.RuleRead)
	app.POST("/RuleDelete", api.RuleDelete)
	app.POST("/RuleList", api.RuleList)

	err = app.Run("0.0.0.0:8888")
	if err != nil {
		panic(err)
	}
}
