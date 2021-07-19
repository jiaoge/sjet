package function

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func init() {
	globalFunc["md5"] = md5Func
	globalFunc["base64"] = base64Func
	globalFunc["base64Decode"] = base64DecodeFunc

	globalFunc["substring"] = substringFunc
	globalFunc["indexOf"] = indexOfFunc

}

func substringFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()

	prefix := int32(a.Get(1).Interface().(float64))

	if a.Get(0).Type().Kind() == reflect.Float64 {
		num := value.(float64)
		val, _ := decimal.NewFromFloat(num).Round(prefix).Float64()
		return reflect.ValueOf(val)
	} else {
		str := value.(string)
		return reflect.ValueOf(str[0:prefix])
	}
}

func indexOfFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface().(string)
	key := a.Get(1).Interface().(string)
	return reflect.ValueOf(strings.Index(value, key))
}

func md5Func(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	h := md5.New()
	h.Write([]byte(str))
	val := hex.EncodeToString(h.Sum(nil))
	return reflect.ValueOf(val)
}

func base64Func(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	b := []byte(str)
	sEnc := base64.StdEncoding.EncodeToString(b)
	return reflect.ValueOf(sEnc)
}

func base64DecodeFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logrus.Error(err, str)
	}
	return reflect.ValueOf(sDec)
}
