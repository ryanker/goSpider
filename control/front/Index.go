package front

import (
	"net/http"

	"github.com/xiuno/gin"

	"../api"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "Index.html", gin.H{})
}

func Item(c *gin.Context) {
	c.HTML(http.StatusOK, "Item.html", gin.H{
		"Rid":   c.Query("Rid"),
		"Table": c.Query("Table"),
	})
}
func User(c *gin.Context) {
	c.HTML(http.StatusOK, "User.html", gin.H{})
}

func Log(c *gin.Context) {
	c.HTML(http.StatusOK, "Log.html", gin.H{})
}

func Sys(c *gin.Context) {
	c.HTML(http.StatusOK, "Sys.html", gin.H{})
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

func Logout(c *gin.Context) {
	api.UserSetCookie(c, "", 1)
	c.Redirect(http.StatusFound, "/")
	// c.HTML(http.StatusOK, "Logout.html", gin.H{})
}
