package front

import (
	"net/http"

	"github.com/xiuno/gin"

	"../api"
)

func Index(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Index.html", gin.H{"User": User})
}

func Item(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Item.html", gin.H{
		"User":  User,
		"Rid":   c.Query("Rid"),
		"Table": c.Query("Table"),
	})
}
func User(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "User.html", gin.H{"User": User})
}

func Log(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Log.html", gin.H{"User": User})
}

func Sys(c *gin.Context) {
	User := api.UserGet(c)
	c.HTML(http.StatusOK, "Sys.html", gin.H{"User": User})
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

func Logout(c *gin.Context) {
	api.UserSetCookie(c, "", 1)
	c.Redirect(http.StatusFound, "/")
}
