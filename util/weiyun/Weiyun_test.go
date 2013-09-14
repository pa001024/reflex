package weiyun

import (
	"fmt"
	"testing"
)

func TestGetShareFileTrueUrl(t *testing.T) {
	data := []string{
		"79ff0e64bd9c88cbd53a0ae72c14401c",
		"http://url.cn/PVX6FN",
		"http://share.weiyun.com/79ff0e64bd9c88cbd53a0ae72c14401c",
	}
	for _, v := range data {
		dl, err := ParseAndDL(v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(dl)
	}
}
