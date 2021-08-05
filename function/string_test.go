package function

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestString(t *testing.T) {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	ti := primitive.DateTime(1626921098228)

	// tt := time.Unix(ti, 0)

	// primitive.DateTime.Time()

	fmt.Println(ti.Time().Format("2006-01-02"))

	//Golang 实现 float64 转 string
	var km = 9900.101
	str := fmt.Sprintf("%f", km)

	fmt.Println(str)

	strKm := strconv.FormatFloat(km, 'f', -1, 64)
	fmt.Println("StrKm = ", strKm)
}
