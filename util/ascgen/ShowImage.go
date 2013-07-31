package image

import (
	"bufio"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type Console struct {
	FontW, FontH uint // 单个字符的像素大小 = 屏幕像素大小/字符行数
	ScreenW      uint // 屏幕能够显示的字符行数/列数
}

func (c Console) Width() uint {
	return c.FontW * c.ScreenW
}
func (c Console) Ratio() float64 {
	return float64(c.FontW) / float64(c.FontH)
}

func ShowJpg(iw io.Writer, r io.Reader, console Console) (err error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return
	}
	Show(iw, img, console)
	return
}
func ShowPng(iw io.Writer, r io.Reader, console Console) (err error) {
	img, err := png.Decode(r)
	if err != nil {
		return
	}
	Show(iw, img, console)
	return
}

func Show(iw io.Writer, img image.Image, console Console) {
	const clist = "MMMMMMMMMM@@@@@@@@@$$$$$$$$$WWWWWWWWW#######BBBBBBBB000000000888888888mmmmmmmmmmmmHHHHHHHHHHHHZZZZZZZZZZEEEEEEEEaaaaaaa2222222SSSSSSSXXXXXXXXXXXnnnnnnnn777777777jjjjjjjlllllllvvvvvvvvzzzzzzzzrrrrrrrrriiiiiiiii;;;;;;;;::::::::::,,,,,,,...........           " //"       ......,,,,,,,,-------========:::::::;;;;;;;/////////////aaaaaaaammmmmmmbbbbbbbbkkkkkkkkDDDDDEEEHHHHRRRRQQQQQQQ$$$$$$%%%%%%%%"
	// Lanczos3缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 原结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
	h := uint(float64(s.Y) / float64(s.X) * float64(w) * console.Ratio())
	m := resize.Resize(w, h, img, resize.Lanczos3)
	bw := bufio.NewWriter(iw)
	for y := 0; y < int(h); y++ {
		for x := (0); x < int(w); x++ {
			r, g, b, a := m.At(x, y).RGBA()
			cindex := (r + g + b) / 3 * a / 65536 / 256
			io.WriteString(bw, clist[cindex:cindex+1])
		}
		io.WriteString(bw, "\n")
	}
	bw.Flush()
}
