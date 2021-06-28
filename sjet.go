package sjet

import (
	"github.com/shoppehub/sjet/engine"
	"github.com/shoppehub/sjet/function"
)

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
