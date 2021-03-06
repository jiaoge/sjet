package context

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
)

func handlerRoute(vars *jet.VarMap, c *gin.Context) {
	vars.SetFunc("redirect", func(a jet.Arguments) reflect.Value {
		rpath := a.Get(0).Interface().(string)
		panic("redirect::::" + rpath)
	})
}
