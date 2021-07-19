package function

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/conf"
)

func init() {
	globalFunc["urls"] = urlsFunc
}

// 获取urls
func urlsFunc(a jet.Arguments) reflect.Value {

	if a.NumOfArguments() > 0 {
		str := a.Get(0).Interface().(string)
		return reflect.ValueOf(conf.GetString("urls." + str))
	}

	return reflect.ValueOf(conf.GetMap("urls"))
}
