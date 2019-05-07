package mysql

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertImg(t *testing.T) {
	img := Img{12345, "test.jpg", 1, "9487529347589", "test", time.Now()}
	InsertImg(&img)
}

func TestQueryByAlias(t *testing.T) {
	img, err := QueryByAlias("9487529347581")
	if err != nil {
		fmt.Printf("query failed, err:%v", err)
	}
	fmt.Println(img.ID)
}

func TestQueryMulti(t *testing.T) {
	QueryMulti()
}
