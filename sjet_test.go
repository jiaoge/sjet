package sjet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/conf"
	"github.com/shoppehub/sjet/context"
	"github.com/shoppehub/sjet/engine"
)

func TestInt(t *testing.T) {
	conf.Init("")

	engine := CreateWithMem()

	router := SetupRouter(engine)

	w := performRequest(router, "GET", "/")
	fmt.Println(w.Body.String())
}

func SetupRouter(engine *engine.TemplateEngine) *gin.Engine {
	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {

		templateContext := context.InitTemplateContext(engine, c)

		template := `

		{{substring(2.455,2)}}

		{{substring("2.455",2)}}

		`

		result, _ := RenderMemTemplate(engine, templateContext, c, "demo", template)

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})
	return router
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
