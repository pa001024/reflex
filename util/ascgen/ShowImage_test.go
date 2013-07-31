package image

import (
	"fmt"
	"os"
	"testing"

	ct "github.com/daviddengcn/go-colortext"
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

func TestShowColor(t *testing.T) {
	f, err := os.Open("lol.png")
	if err != nil {
		t.Fail()
		return
	}
	err = ShowPngColor(os.Stdout, f, Console{6, 14, 120})
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetConsoleColor(t *testing.T) {
	for b := uint32(0); b <= 5; b++ {
		for g := uint32(0); g <= 5; g++ {
			for r := uint32(0); r <= 5; r++ {
				// fmt.Println("RGB:", r*0x33, g*0x33, b*0x33)
				g, h, j, k := GetConsoleColor(r*256*0x33, g*256*0x33, b*256*0x33, 65536)
				ct.ChangeColor(g, h, j, k)
				if k {
					fmt.Print(j << 1)
				} else {
					fmt.Print(j)
				}
				fmt.Print(" ")
				ct.ResetColor()
			}
		}
	}
}
