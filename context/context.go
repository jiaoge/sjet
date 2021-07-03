package context

import (
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shoppehub/sjet/engine"
)

// 模板引擎配置
type TemplateContext struct {
	Vars    *jet.VarMap
	Context *map[string]interface{}

	Template *jet.Template

	Module      string
	Page        string
	TemplName   string
	TempatePath string
}

var TemplateRoot = "pages"

func (ctx *TemplateContext) Namespace() string {
	return strings.ReplaceAll(ctx.TempatePath, "/", "_")
}

func (ctx *TemplateContext) GetTemplateDir() string {
	return strings.Join([]string{TemplateRoot, path.Dir((ctx.TempatePath))}, "/")
}

func (ctx *TemplateContext) FindTemplate(t *engine.TemplateEngine) error {
	// templatePath := strings.Join([]string{TemplateRoot, ctx.Module, ctx.Page, ctx.TemplName}, "/")

	var view *jet.Template
	var err error
	if view, err = t.Views.GetTemplate(ctx.TempatePath); err != nil {
		ctx.TempatePath += "/index"
	}
	if view, err = t.Views.GetTemplate(ctx.TempatePath); err != nil {
		return err
	}

	ctx.Template = view
	return nil
}

// 初始化模板
func InitTemplateContext(t *engine.TemplateEngine, c *gin.Context) *TemplateContext {
	vars := make(jet.VarMap)
	handlerGetCtx(&vars, c)

	context := make(map[string]interface{})
	handlerContext(&vars, &context)

	ctxData := TemplateContext{
		Vars:    &vars,
		Context: &context,
	}
	ctxData.TempatePath = strings.TrimPrefix(c.Request.URL.Path, "/")

	// handlerTemplateFile(c, &ctxData)
	vars.Set("namespace", ctxData.Namespace())

	return &ctxData
}

// 解析模板路径 /:module/:page/:templ
func handlerTemplateFile(c *gin.Context, ctx *TemplateContext) {

	// paths := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
	ctx.TempatePath = strings.TrimPrefix(c.Request.URL.Path, "/")

	module := c.Params.ByName("module")
	if module == "" {
		module = "index"
	}
	page := c.Params.ByName("page")
	if page == "" {
		page = "index"
	}

	templ := c.Params.ByName("templ")
	if templ == "" {
		templ = "index"
	}
	ctx.Module = module
	ctx.Page = page
	ctx.TemplName = templ
}

func getParamInContext(key string, c *gin.Context) interface{} {
	if val, ok := c.GetQuery(key); ok {
		return val
	}
	if val, ok := c.GetPostForm(key); ok {
		return val
	}
	if val, ok := c.Params.Get(key); ok {
		return val
	}
	var body map[string]interface{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err == nil {
		return body[key]
	}
	return ""
}

func handlerGetCtx(vars *jet.VarMap, c *gin.Context) {
	vars.SetFunc("getCtx", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()
		return reflect.ValueOf(getParamInContext(key, c))
	})

	vars.SetFunc("getCtxForInt", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c)

		if val == "" {
			return reflect.ValueOf(0)
		}
		val, _ = strconv.ParseInt(c.Query("curPage"), 10, 64)
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getCtxForFloat", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c)

		if val == "" {
			val = float64(0)
			return reflect.ValueOf(val)
		}
		val, _ = strconv.ParseFloat(c.Query("curPage"), 64)
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getCtxForBool", func(a jet.Arguments) reflect.Value {
		key := a.Get(0).String()

		val := getParamInContext(key, c)

		if val == "" {
			return reflect.ValueOf(false)
		}
		val, _ = strconv.ParseBool(val.(string))
		return reflect.ValueOf(val)
	})

	vars.SetFunc("getURL", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(c.Request.URL)
	})
}

func handlerContext(vars *jet.VarMap, context *map[string]interface{}) {
	vars.SetFunc("context", func(a jet.Arguments) reflect.Value {
		ctx := *context

		if a.NumOfArguments() == 1 {
			if val, ok := ctx[a.Get(0).String()]; ok {
				return reflect.ValueOf(val)
			}
			return reflect.ValueOf("")
		}

		ctx[a.Get(0).String()] = a.Get(1).Interface()
		return reflect.ValueOf("")
	})
}
