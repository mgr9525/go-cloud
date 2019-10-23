package gocloud

import "gopkg.in/macaron.v1"

var mapController = make([]IController, 0)

type IController interface {
	GetPath() string
	Routes()
	Mid() []macaron.Handler
}

/*type Controller struct {
	path		string
}
func (e*Controller)Init(ph string){
	e.path=ph
	MapController[ph]=e
}
func (e*Controller)Routes(){

}*/

func RegController(c IController) {
	mapController = append(mapController, c)
}

func runController() {
	for _, v := range mapController {
		/*if len(v.GetPath())<=0 {
			v.Routes()
		}*/

		/*mids:=make([]macaron.Handler,0)
		if v.Mids()!=nil {
			mids=v.Mids()
		}*/
		mid := v.Mid()
		hld := func() {
			v.Routes()
		}
		if mid == nil {
			Web.Group(v.GetPath(), hld)
		} else {
			Web.Group(v.GetPath(), hld, mid...)
		}
	}
}
