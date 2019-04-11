package front

import (
	"net/http"

	"github.com/xiuno/gin"
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
