package engine

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/conf"
)

// 模板引擎配置
type TemplateEngine struct {
	Views  *jet.Set
	IsFs   *bool
	Loader *TemplateLoader
}

// 模板文件加载器
type TemplateLoader interface {
	jet.Loader
	Set(templatePath, contents string)
	Delete(templatePath string)
}

// 根据内存创建
func CreateWithMem() *TemplateEngine {
	return create(false, "")
}

// 根据本地文件创建
func CreateWithFile(dir string) *TemplateEngine {
	return create(true, dir)
}

// 创建新实例
func create(isFs bool, dir string) *TemplateEngine {
	var views *jet.Set
	ecache := ECache{}

	var loader TemplateLoader
	if isFs {
		loader = NewOSFileSystemLoader(dir, views, &ecache)
	} else {
		loader = jet.NewInMemLoader()
	}

	extension := jet.WithTemplateNameExtensions([]string{".jet", ".html.jet", ".jet.html"})
	cache := jet.WithCache(&ecache)
	if conf.GetBool("production") {
		views = jet.NewSet(
			loader,
			jet.InDevelopmentMode(),
			cache,
			extension,
		)
	} else {
		views = jet.NewSet(
			loader,
			cache,
			extension,
		)
	}

	return &TemplateEngine{
		Views:  views,
		IsFs:   &isFs,
		Loader: &loader,
	}
}
