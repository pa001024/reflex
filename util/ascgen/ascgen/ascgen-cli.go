package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	asc "github.com/pa001024/MoeWorker/util/ascgen"
)

func main() {
	pF := flag.String("f", "", "the file you want to render")
	pW := flag.Uint("w", 120, "Console.width")
	pN := flag.Uint("n", 6, "Font.width")
	pH := flag.Uint("h", 14, "Font.height")
	pC := flag.Bool("c", false, "Show Color")
	flag.Parse()

	if *pF == "" {
		flag.Usage()
		return
	}

	ext := path.Ext(*pF)
	console := asc.Console{*pN, *pH, *pW}
	r, err := os.Open(*pF)
	if err != nil {
		fmt.Println(err)
	}
	switch ext {
	default:
		err = asc.ShowPng(os.Stdout, r, console, *pC)
	case ".jpg", ".jpeg":
		err = asc.ShowJpg(os.Stdout, r, console, *pC)
	case ".gif":
		err = asc.ShowGif(os.Stdout, r, console, *pC)
	}
	if err != nil {
		fmt.Println(err)
	}
}
