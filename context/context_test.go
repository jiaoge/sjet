package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {

	str := "/index"

	paths := strings.Split(strings.TrimPrefix(str, "/"), "/")

	pathSize := len(paths)
	fmt.Println(pathSize, paths)

	for i := 1; i < pathSize; i++ {
		fmt.Println(i, paths[i])
	}

	fmt.Println(path.Dir("a/b/c"))

}

func TestHelloWorld(t *testing.T) {
	body := gin.H{
		"Hello": "World",
	}

	router := SetupRouter()

	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["Hello"]

	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["Hello"], value)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {

		fmt.Println("path", c.Request.URL.Path)

		c.JSON(http.StatusOK, gin.H{
			"Hello": "World",
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
