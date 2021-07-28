package function

import (
	"fmt"
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
}
