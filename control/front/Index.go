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

func Log(c *gin.Context) {
	c.HTML(http.StatusOK, "Log.html", gin.H{})
}
