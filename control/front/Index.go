package front

import (
	"net/http"

	"github.com/xiuno/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "Index.html", gin.H{})
}
