package function

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
)

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
