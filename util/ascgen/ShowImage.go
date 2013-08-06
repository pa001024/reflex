package image

import (
	"bufio"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/nfnt/resize"
)

// 从文件读取
func ShowFile(iw io.Writer, r io.Reader, console Console, useColor, simple bool) (err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return
	}
	switch {
	case useColor && simple:
		ShowSimpleColor(iw, img, console)
	case useColor && !simple:
		ShowColor(iw, img, console)
	default:
		Show(iw, img, console)
	}
	return
}

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

const clist64 = "MMMNNBBEEFRRWW9956OGQDKU33ICJXV7i11jjllllrr;;;;:::,,,,.....     "

// 将图片渲染到文字(有颜色版)
func ShowColor(iw io.Writer, img image.Image, console Console) {
	// Lanczos3缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
	h := uint(float64(s.Y) / float64(s.X) * float64(w) * console.Ratio())
	m := resize.Resize(w, h, img, resize.Bicubic)
	for y := 0; y < int(h); y++ {
		for x := (0); x < int(w); x++ {
			r, g, b, a := m.At(x, y).RGBA()
			r, g, b, a = r>>10, g>>10, b>>10, a>>10 // 转换到6位色 0~63
			r, g, b = r*a/63, g*a/63, b*a/63        // Alpha混合
			m, fg, fgl, bg, bgl := GetConsoleColor(r, g, b)
			ct.ChangeColor(fg, fgl, bg, bgl)
			// cindex := (r + g + b) / 3
			io.WriteString(iw, clist64[m:m+1])
		}
		io.WriteString(iw, "\n")
	}
	ct.ResetColor()
}

var (
	rgb1_c = map[rgb1]ct.Color{
		rgb1{0, 0, 0}: ct.Black,
		rgb1{1, 1, 1}: ct.White,
		rgb1{1, 0, 0}: ct.Red,
		rgb1{0, 1, 0}: ct.Green,
		rgb1{0, 0, 1}: ct.Blue,
		rgb1{1, 1, 0}: ct.Yellow,
		rgb1{0, 1, 1}: ct.Cyan,
		rgb1{1, 0, 1}: ct.Magenta,
	}
)

type rgb1 struct{ r, g, b uint32 }
type rgb2 struct{ r, g, b uint32 }
type rgb6 struct{ r, g, b uint32 }
type rgb1b struct {
	r, g, b uint32
	l       bool
}

func (r rgb1b) ccolor() (ct.Color, bool) { return rgb1_c[rgb1{r.r, r.g, r.b}], r.l }
func (r rgb1b) rgb6() rgb6 {
	d := uint32(31)
	if r.l {
		d = 63
	}
	return rgb6{r.r * d, r.g * d, r.b * d}
}
func (r rgb2) rgb1b() rgb1b {
	return rgb1b{(r.r + 1) / 2, (r.g + 1) / 2, (r.b + 1) / 2, (r.r+r.g+r.b) > 3 || r.r > 1 || r.g > 1 || r.b > 1}
}
func (r rgb6) rgb2() rgb2          { return rgb2{(r.r + 16) / 32, (r.g + 16) / 32, (r.b + 16) / 32} }
func (r rgb6) sub(d rgb6) rgb6     { return rgb6{r.r - d.r, r.g - d.g, r.b - d.b} }
func (r rgb6) times(t uint32) rgb6 { return rgb6{r.r * t, r.g * t, r.b * t} }

func GetConsoleColor(r, g, b uint32) (m int, fg ct.Color, fgl bool, bg ct.Color, bgl bool) {
	// 从6位颜色转换到2位颜色
	h := rgb6{r, g, b}
	// 背景色使用最接近的颜色
	j := h.rgb2().rgb1b()
	bg, bgl = j.ccolor()
	// 仿色:
	// 统计最高混合比: 字符'M' 50% 最低 字符' ' 0%
	// 理论可用率66% 最高精确度 6位色
	// 设: 实色h = rgb6(12,56,31) = rgb2(0,2,1) = rgb1b(0,1,0,false) = rgb6(0,31,0)
	//     补色f rgb6(x,y,z)
	//     补色混合比 = (m/n)
	// 则: h = (f * (m/n) + h.rgb1b/2
	//   (rgb6(x,y,z)*n + rgb1b(0,1,0,false))/2 = rgb6(12,56,31)
	//   rgb6(x,y,z)*n = rgb6(12,56,31)*2 - rgb6(0,31,0)
	//   rgb6(x,y,z)*(m/63) = rgb6(24,81,62)
	//   rgb2(1,2,2)*(m/63) = rgb1b(0,1,1,true)*(m/63) = rgb6(0,63,63)/63*m
	//   {
	//      rgb6(x,y,z) = rgb6(0,63,63) = h*2-h.rgb1b
	//      m = (56+31)/2=44 = (Z[h](h.?>1)/C[h](h.?>1))
	//   }
	f := h.times(2).sub(j.rgb6())
	fc := f.rgb2().rgb1b()
	m, n := 0, 0
	if fc.r > 0 {
		m += int(f.r)
		n++
	}
	if fc.g > 0 {
		m += int(f.g)
		n++
	}
	if fc.b > 0 {
		m += int(f.b)
		n++
	}
	if n > 0 {
		m = m * 2 / n // 因为最高混合比(50%)的限制 先*2 然后削掉
	}
	m = 63 - m
	if m < 0 {
		m = 0
	}
	if m > 63 {
		m = 63
	}
	fg, fgl = fc.ccolor()
	return
}

// 实用小函数 模拟 b?l:r
func sw(b bool, l, r uint32) uint32 {
	if b {
		return l
	}
	return r
}

// 将图片渲染到彩色文字
func ShowSimpleColor(iw io.Writer, img image.Image, console Console) {
	// Bicubic缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
	h := uint(float64(s.Y) / float64(s.X) * float64(w) * console.Ratio())
	m := resize.Resize(w, h, img, resize.Bicubic)
	for y := 0; y < int(h); y++ {
		for x := (0); x < int(w); x++ {
			r, g, b, a := m.At(x, y).RGBA()
			r, g, b, a = r>>10, g>>10, b>>10, a>>10 // 转换到6位色 0~63
			r, g, b = r*a/63, g*a/63, b*a/63        // Alpha混合
			_, _, _, bg, bgl := GetConsoleColor(r, g, b)
			ct.ChangeColor(bg, bgl, ct.None, false)
			cindex := (r + g + b) / 3
			io.WriteString(iw, clist64[cindex:cindex+1])
		}
		io.WriteString(iw, "\n")
	}
	ct.ResetColor()
}

// 将图片渲染到文字
func Show(iw io.Writer, img image.Image, console Console) {
	// Bicubic缩放图片到控制台屏幕大小
	w := console.ScreenW - 1 // 标宽 换行符竟然也算一个?
	s := img.Bounds().Size()
	// 标高 = 结果高 * 宽高比 = (原高 / 原宽 * 标宽) * 宽高比
	h := uint(float64(s.Y) / float64(s.X) * float64(w) * console.Ratio())
	m := resize.Resize(w, h, img, resize.Bicubic)
	bw := bufio.NewWriter(iw)
	for y := 0; y < int(h); y++ {
		for x := (0); x < int(w); x++ {
			r, g, b, a := m.At(x, y).RGBA()
			cindex := (r + g + b) / 3 * a / 65535 / 256 / 4
			io.WriteString(bw, clist64[cindex:cindex+1])
		}
		io.WriteString(bw, "\n")
	}
	bw.Flush()
}
