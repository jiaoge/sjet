package sjet

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/sjet/context"
	"github.com/shoppehub/sjet/engine"

	"github.com/shoppehub/sjet/function"
)

var customFunc = make(map[string]CustomFunc)

// 创建文件路径的模板引擎
func CreateWithFile(dir string) *engine.TemplateEngine {
	e := engine.CreateWithFile(dir)
	function.InitGlobalFunc(e)
	return e
}

// 创建内存存储的模板引擎
func CreateWithMem() *engine.TemplateEngine {
	e := engine.CreateWithMem()
	function.InitGlobalFunc(e)
	return e
}

func RenderHTMLTemplate(eng *engine.TemplateEngine, c *gin.Context) {

	templateContext := context.InitTemplateContext(eng, c)

	err := templateContext.FindTemplate(eng)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	for key, v := range customFunc {
		templateContext.Vars.SetFunc(key, v(c))
	}

}

type CustomFunc func(c *gin.Context) jet.Func

// 注册自定义的函数
func RegCustomFunc(funcName string, v CustomFunc) {
	customFunc[funcName] = v
}
