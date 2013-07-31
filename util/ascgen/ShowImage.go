package image

import (
	"bufio"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/nfnt/resize"
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

func ShowGif(iw io.Writer, r io.Reader, console Console, useColor bool) (err error) {
	img, err := gif.Decode(r)
	if err != nil {
		return
	}
	if useColor {
		ShowColor(iw, img, console)
	} else {
		Show(iw, img, console)
	}
	return
}
func ShowJpg(iw io.Writer, r io.Reader, console Console, useColor bool) (err error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return
	}
	if useColor {
		ShowColor(iw, img, console)
	} else {
		Show(iw, img, console)
	}
	return
}
func ShowPng(iw io.Writer, r io.Reader, console Console, useColor bool) (err error) {
	img, err := png.Decode(r)
	if err != nil {
		return
	}
	if useColor {
		ShowColor(iw, img, console)
	} else {
		Show(iw, img, console)
	}
	return
}

// 将图片渲染到文字(有颜色版)
func ShowColor(iw io.Writer, img image.Image, console Console) {
	const clist = "MMMMMMMMMM@@@@@@@@@$$$$$$$$$WWWWWWWWW#######BBBBBBBB000000000888888888mmmmmmmmmmmmHHHHHHHHHHHHZZZZZZZZZZEEEEEEEEaaaaaaa2222222SSSSSSSXXXXXXXXXXXnnnnnnnn777777777jjjjjjjlllllllvvvvvvvvzzzzzzzzrrrrrrrrriiiiiiiii;;;;;;;;::::::::::,,,,,,,...........           " //"       ......,,,,,,,,-------========:::::::;;;;;;;/////////////aaaaaaaammmmmmmbbbbbbbbkkkkkkkkDDDDDEEEHHHHRRRRQQQQQQQ$$$$$$%%%%%%%%"
	// Lanczos3缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
	h := uint(float64(s.Y) / float64(s.X) * float64(w) * console.Ratio())
	m := resize.Resize(w, h, img, resize.Lanczos3)
	for y := 0; y < int(h); y++ {
		for x := (0); x < int(w); x++ {
			r, g, b, a := m.At(x, y).RGBA()
			ct.ChangeColor(GetConsoleColor(r, g, b, a))
			cindex := (r + g + b) / 3 * a / 65536 / 256
			io.WriteString(iw, clist[cindex:cindex+1])
		}
		io.WriteString(iw, "\n")
	}
	ct.ResetColor()
}

func GetConsoleColor(r, g, b, a uint32) (fg ct.Color, fgl bool, bg ct.Color, bgl bool) {
	r, g, b = (r*a/65536/256+32)/64, (g*a/65536/256+32)/64, (b*a/65536/256+32)/64 // 从16位颜色转换到3位颜色 Alpha混合底色为黑
	type truecolor struct{ r, g, b uint32 }
	type fbcolor struct {
		c ct.Color
		b bool
	}
	c2i := map[truecolor]fbcolor{
		truecolor{0, 0, 0}: {ct.Black, false},
		truecolor{1, 0, 0}: {ct.Red, false},
		truecolor{0, 1, 0}: {ct.Green, false},
		truecolor{0, 0, 1}: {ct.Blue, false},
		truecolor{1, 1, 0}: {ct.Yellow, false},
		truecolor{0, 1, 1}: {ct.Cyan, false},
		truecolor{1, 0, 1}: {ct.Magenta, false},
		truecolor{1, 1, 1}: {ct.White, false},
	}
	sr, sg, sb := (r+1)/2, (g+1)/2, (b+1)/2                        //去掉1位
	bl := sr > 1 && sg > 1 || sr > 1 && sb > 1 || sg > 1 && sb > 1 // 计算亮度
	sr, sg, sb = (r+2)/4, (g+2)/4, (b+2)/4                         //去掉2位
	bgi := c2i[truecolor{sr, sg, sb}]
	bg, bgl = bgi.c, bl
	fg, fgl = ct.White, true
	return
}

// 将图片渲染到文字
func Show(iw io.Writer, img image.Image, console Console) {
	const clist = "MMMMMMMMMM@@@@@@@@@$$$$$$$$$WWWWWWWWW#######BBBBBBBB000000000888888888mmmmmmmmmmmmHHHHHHHHHHHHZZZZZZZZZZEEEEEEEEaaaaaaa2222222SSSSSSSXXXXXXXXXXXnnnnnnnn777777777jjjjjjjlllllllvvvvvvvvzzzzzzzzrrrrrrrrriiiiiiiii;;;;;;;;::::::::::,,,,,,,...........           " //"       ......,,,,,,,,-------========:::::::;;;;;;;/////////////aaaaaaaammmmmmmbbbbbbbbkkkkkkkkDDDDDEEEHHHHRRRRQQQQQQQ$$$$$$%%%%%%%%"
	// Lanczos3缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
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
