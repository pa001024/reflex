package image

import (
	"fmt"
	"os"
	"testing"
)

func TestShow(t *testing.T) {
	f, err := os.Open("pic.jpg")
	if err != nil {
		t.Fail()
		return
	}
	err = ShowJpg(os.Stdout, f, Console{6, 14, 120})
	if err != nil {
		fmt.Println(err)
	}
}
