package routes

import (
	"github.com/mgr9525/go-cloud"
	"gopkg.in/macaron.v1"
)

func IndexHandler(c *macaron.Context, contJSON gocloud.ContJSON) {
	// Check auto-login.
	/*uname := c.GetCookie(setting.CookieUserName)
	if len(uname) != 0 {
		c.Redirect(setting.AppSubURL + "/user/login")
		return
	}*/

	defer gocloud.RuisRecovers("IndexHandler", func() {
		c.PlainText(500, []byte("server error!"))
	})

	c.Data["PageIsHome"] = true
	c.HTML(200, "index")
}
