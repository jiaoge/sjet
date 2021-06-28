package context

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shoppehub/sjet/engine"
)

// 模板引擎配置
type TemplateContext struct {
	Vars    *jet.VarMap
	Context *map[string]interface{}

	Module    string
	Page      string
	TemplName string
}

func InitTemplateContext(t *engine.TemplateEngine, c *gin.Context) *TemplateContext {
	vars := make(jet.VarMap)

	handlerGetCtx(&vars, c)

	var context map[string]interface{}
	handlerContext(&vars, &context)

	vars.Set("path", c.Request.URL.Path)

	return &TemplateContext{
		Vars:    &vars,
		Context: &context,
	}
}

func handlerGetCtx(vars *jet.VarMap, c *gin.Context) {
	vars.SetFunc("getCtx", func(a jet.Arguments) reflect.Value {

		key := a.Get(0).String()
		if val, ok := c.GetQuery(key); ok {
			return reflect.ValueOf(val)
		}
		if val, ok := c.GetPostForm(key); ok {
			return reflect.ValueOf(val)
		}
		if val, ok := c.Params.Get(key); ok {
			return reflect.ValueOf(val)
		}
		var body map[string]interface{}

		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			return reflect.ValueOf("")
		}
		return reflect.ValueOf(body[key])
	})

	vars.SetFunc("getURL", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(c.Request.URL)
	})
}

func handlerContext(vars *jet.VarMap, context *map[string]interface{}) {
	vars.SetFunc("context", func(a jet.Arguments) reflect.Value {
		ctx := *context
		ctx[a.Get(0).String()] = a.Get(1).Interface()
		return reflect.ValueOf(&ctx)
	})
}
