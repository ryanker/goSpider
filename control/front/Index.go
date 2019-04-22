package front

import (
	"net/http"

	"github.com/xiuno/gin"

	"../api"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

func Logout(c *gin.Context) {
	api.UserSetCookie(c, "", -1)
	c.Redirect(http.StatusFound, "/")
}

func Index(c *gin.Context) {
	User := api.UserGet(c)
	name := "Index.html"
	if User.Gid != 1 {
		name = "UserIndex.html" // 普通用户首页
	}
	c.HTML(http.StatusOK, name, gin.H{"User": User})
}

func Show(c *gin.Context) {
	User := api.UserGet(c)
	Rid := c.Query("Rid")
	Table := c.Query("Table")
	if Table != "Content" && Table != "ContentDownload" && Table != "ListDownload" {
		Table = "List"
	}
	c.HTML(http.StatusOK, "Show"+Table+".html", gin.H{
		"User": User,
		"Rid":  Rid,
	})
}

func Read(c *gin.Context) {
	User := api.UserGet(c)
	Rid := c.Query("Rid")
	Lid := c.Query("Lid")
	c.HTML(http.StatusOK, "Read.html", gin.H{
		"User": User,
		"Rid":  Rid,
		"Lid":  Lid,
	})
}

func Log(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Log.html", gin.H{"User": User})
}

func User(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "User.html", gin.H{"User": User})
}

func Sys(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Sys.html", gin.H{"User": User})
}

func Table(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Table.html", gin.H{
		"User":  User,
		"Rid":   c.Query("Rid"),
		"Table": c.Query("Table"),
	})
}
