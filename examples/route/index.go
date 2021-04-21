package route

import (
	"github.com/gin-gonic/gin"
	gocloud "github.com/mgr9525/go-cloud"
)

type IndexController struct{}

func (IndexController) GetPath() string {
	return ""
}
func (IndexController) GetMid() gin.HandlerFunc {
	return nil
}
func (c *IndexController) Routes(g gin.IRoutes) {
	g.Any("/", c.index)
	g.Any("/test", gocloud.JsonHandle(c.test))
}
func (IndexController) index(c *gin.Context) {
	c.HTML(200, "index.html", map[string]interface{}{"Name": "123"})
}
func (IndexController) test(c *gin.Context) {

}
