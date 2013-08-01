package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	asc "github.com/pa001024/MoeWorker/util/ascgen"
)

func Show(iw io.Writer, r io.Reader, console asc.Console, useColor, simple bool) (err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return
	}
	switch {
	case useColor && simple:
		asc.ShowSimpleColor(iw, img, console)
	case useColor && !simple:
		asc.ShowColor(iw, img, console)
	default:
		asc.Show(iw, img, console)
	}
	return
}

func main() {
	pF := flag.String("i", "", "the file you want to render")
	pW := flag.Uint("w", 120, "Console.width")
	pN := flag.Uint("n", 6, "Font.width")
	pH := flag.Uint("h", 14, "Font.height")
	pC := flag.Bool("c", false, "Show Color")
	pS := flag.Bool("s", false, "Show Color Simple")
	flag.Parse()

	if *pF == "" {
		flag.Usage()
		return
	}

	console := asc.Console{*pN, *pH, *pW}
	r, err := os.Open(*pF)
	if err != nil {
		fmt.Println(err)
	}
	Show(os.Stdout, r, console, *pC, *pS)
	if err != nil {
		fmt.Println(err)
	}
}
