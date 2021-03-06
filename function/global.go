package function

import (
	"math"
	"net/url"
	"reflect"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/sjet/engine"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var globalFunc = make(map[string]jet.Func)

// 初始化全局函数
func InitGlobalFunc(t *engine.TemplateEngine) {

	for k, v := range globalFunc {
		t.Views.AddGlobalFunc(k, v)
	}

	// 把数字转换为int数组
	t.Views.AddGlobalFunc("numArray", numArrayFunc)
	// 支持把数据转换为字符串，比如 objectId
	t.Views.AddGlobalFunc("oid", oidFunc)

	t.Views.AddGlobalFunc("formatUrlPath", formatUrlPathFunc)

	t.Views.AddGlobalFunc("map", mapFunc)
	t.Views.AddGlobalFunc("put", putFunc)
	t.Views.AddGlobalFunc("append", appendFunc)

	t.Views.AddGlobalFunc("array", arrayFunc)

	t.Views.AddGlobalFunc("aggregate", aggregateFunc)
	t.Views.AddGlobalFunc("pipeline", aggregateFunc)

	t.Views.AddGlobalFunc("m", mFunc)
	t.Views.AddGlobalFunc("d", dFunc)

	t.Views.AddGlobalFunc("parseInt", parseIntFunc)
	t.Views.AddGlobalFunc("ceil", ceilFunc)
	t.Views.AddGlobalFunc("floor", floorFunc)

	t.Views.AddGlobalFunc("log", logFunc)

	t.Views.AddGlobalFunc("exit", exitFunc)
}

func oidFunc(a jet.Arguments) reflect.Value {
	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}
	oid, _ := primitive.ObjectIDFromHex(a.Get(0).String())
	return reflect.ValueOf(oid)
}

func formatUrlPathFunc(a jet.Arguments) reflect.Value {
	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}
	u, _ := url.Parse(a.Get(0).Interface().(string))
	return reflect.ValueOf(u.Path)
}

// 把数字转换为int数组
func numArrayFunc(a jet.Arguments) reflect.Value {
	var total int
	k := a.Get(0).Kind()
	switch k {
	case reflect.Float64:
		total = int(a.Get(0).Float())
	default:
		total = int(a.Get(0).Int())
	}

	nums := make([]int64, total)
	for i := 0; i < total; i++ {
		nums[i] = int64(i + 1)
	}
	return reflect.ValueOf(nums)
}

func mapFunc(a jet.Arguments) reflect.Value {
	if a.NumOfArguments()%2 > 0 {
		return reflect.ValueOf(make(map[string]interface{}))
	}
	m := reflect.ValueOf(make(map[string]interface{}, a.NumOfArguments()/2))
	for i := 0; i < a.NumOfArguments(); i += 2 {

		m.SetMapIndex(a.Get(i), a.Get(i+1))
	}
	return m
}

func putFunc(a jet.Arguments) reflect.Value {
	name := a.Get(0).Type().Name()

	if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		m[a.Get(1).String()] = a.Get(2).Interface()
		return reflect.ValueOf(m)
	} else {
		m := a.Get(0).Interface().(map[string]interface{})
		m[a.Get(1).String()] = a.Get(2).Interface()
		return reflect.ValueOf(m)
	}
}

func appendFunc(a jet.Arguments) reflect.Value {
	name := a.Get(0).Type().Name()
	kind := a.Get(0).Type().Kind()

	if name == "D" {
		m := a.Get(0).Interface().(bson.D)
		e := bson.E{}
		e.Key = a.Get(1).Interface().(string)
		e.Value = a.Get(2).Interface()
		m = append(m, e)
		return reflect.ValueOf(m)
	} else if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]bson.M), a.Get(2).Interface().(bson.M))
			m[a.Get(1).String()] = val
		} else {
			val := []bson.M{a.Get(2).Interface().(bson.M)}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	} else if kind == reflect.Map {
		m := a.Get(0).Interface().(map[string]interface{})
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]interface{}), a.Get(2).Interface())
			m[a.Get(1).String()] = val
		} else {
			val := []interface{}{a.Get(2).Interface()}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	} else if kind == reflect.Slice {
		m := a.Get(0).Interface().([]interface{})
		m = append(m, a.Get(1).Interface())
		return reflect.ValueOf(m)
	}
	return reflect.ValueOf("")
}

func parseIntFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	val, _ := strconv.ParseInt(value.(string), 10, 64)
	return reflect.ValueOf(val)
}

func ceilFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	return reflect.ValueOf(int(math.Ceil(value.(float64))))
}

func floorFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	return reflect.ValueOf(int(math.Floor(value.(float64))))
}

func mFunc(a jet.Arguments) reflect.Value {
	d := bson.M{}
	for i := 0; i < a.NumOfArguments(); i += 2 {
		d[a.Get(i).String()] = a.Get(i + 1).Interface()
	}
	m := reflect.ValueOf(d)
	return m
}

func dFunc(a jet.Arguments) reflect.Value {
	d := bson.D{}
	for i := 0; i < a.NumOfArguments(); i += 2 {
		d = append(d, bson.E{
			Key:   a.Get(i).String(),
			Value: a.Get(i + 1).Interface(),
		})
	}
	m := reflect.ValueOf(d)
	return m
}

func aggregateFunc(a jet.Arguments) reflect.Value {
	var p []bson.D
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface().(bson.D))
	}
	m := reflect.ValueOf(p)
	return m
}

func arrayFunc(a jet.Arguments) reflect.Value {
	var p []interface{}
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface())
	}
	m := reflect.ValueOf(p)
	return m
}

func logFunc(a jet.Arguments) reflect.Value {

	level := a.Get(0).Interface().(string)
	logVal := a.Get(1).Interface()
	switch level {
	case "err":
		logrus.Error(logVal)
	case "info":
		logrus.Info(logVal)
	case "warn":
		logrus.Warn(logVal)
	}

	return reflect.ValueOf("")
}

func exitFunc(a jet.Arguments) reflect.Value {
	panic("exit::::")
}
