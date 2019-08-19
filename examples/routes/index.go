package routes

import (
	"examples/dao"
	"examples/models"
	"fmt"
	"github.com/mgr9525/go-cloud"
	"gopkg.in/macaron.v1"
)

func IndexHandler(c *macaron.Context) {
	// Check auto-login.
	/*uname := c.GetCookie(setting.CookieUserName)
	if len(uname) != 0 {
		c.Redirect(setting.AppSubURL + "/user/login")
		return
	}*/

	defer gocloud.RuisRecovers("IndexHandler", func() {
		c.PlainText(500, []byte("server error!"))
	})

	if gocloud.Db != nil {
		usr := dao.GetUserById(1)
		fmt.Printf("user:%v\n", usr)
		usrs := dao.UserDao.FindOne(&map[string]interface{}{"name": "root"})
		fmt.Printf("users:%v\n", usrs)
		c.Data["User"] = usrs

		//test
		usrt := dao.FindUser(&map[string]interface{}{"nick": "管理员"})
		fmt.Printf("usert:%v\n", usrt)

		plist := dao.UserDao.FindList(gocloud.NewParamMap())
		list := *(plist.(*[]models.User))
		fmt.Printf("userList:%v\nlist[0].Name=%s\n", list, list[0].Name)
		pPage := dao.UserDao.FindPage(gocloud.NewParamMap(), 1, nil)
		pageList := *(pPage.List.(*[]models.User))
		fmt.Printf("pPage:%v\n", pPage)
		fmt.Printf("pageList:%v\n", pageList)
	}

	c.Data["PageIsHome"] = true
	c.HTML(200, "index")
}
