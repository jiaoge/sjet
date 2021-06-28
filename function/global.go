package function

import (
	"fmt"
	"math"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/sjet/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 初始化全局函数
func InitGlobalFunc(t *engine.TemplateEngine) {
	// 把数字转换为int数组
	t.Views.AddGlobalFunc("numArray", numArrayFunc)
	// 支持把数据转换为字符串，比如 objectId
	t.Views.AddGlobalFunc("string", stringFunc)

	t.Views.AddGlobalFunc("map", mapFunc)
	t.Views.AddGlobalFunc("put", putFunc)
	t.Views.AddGlobalFunc("append", appendFunc)

	t.Views.AddGlobalFunc("aggregate", aggregateFunc)
	t.Views.AddGlobalFunc("m", mFunc)
	t.Views.AddGlobalFunc("d", dFunc)

	t.Views.AddGlobalFunc("ceil", putFunc)
	t.Views.AddGlobalFunc("floor", putFunc)

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

func stringFunc(a jet.Arguments) reflect.Value {
	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}

	name := a.Get(0).Type().Name()

	switch name {
	case "ObjectID":
		oid := a.Get(0).Interface().(primitive.ObjectID)
		return reflect.ValueOf(oid.Hex())
	case "int":
		return reflect.ValueOf(fmt.Sprint(a.Get(0).Interface().(int)))
	}

	return reflect.ValueOf(a.Get(0).Interface())
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
	if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]bson.M), a.Get(2).Interface().(bson.M))
			m[a.Get(1).String()] = val
		} else {
			val := []bson.M{a.Get(2).Interface().(bson.M)}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	} else {
		m := a.Get(0).Interface().(map[string]interface{})
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]interface{}), a.Get(2).Interface())
			m[a.Get(1).String()] = val
		} else {
			val := []interface{}{a.Get(2).Interface()}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	}
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
	p := mongo.Pipeline{}
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface().(bson.D))
	}
	m := reflect.ValueOf(p)
	return m
}
