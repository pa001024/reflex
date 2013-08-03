package image

import (
	"bufio"
	"fmt"
	"image"
	"io"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/nfnt/resize"
)

var _ fmt.State

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
			io.WriteString(iw, clist64[int(m):int(m)+1])
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

type rgb1bm struct {
	r, g, b int32
	l       bool
	m       int32
}
type rgb1 struct{ r, g, b int32 }
type rgb2 struct{ r, g, b int32 }
type rgb6 struct{ r, g, b int32 }
type rgb1b struct {
	r, g, b int32
	l       bool
}

func (r rgb1bm) rgb1b() rgb1b {
	return rgb1b{r.r, r.g, r.b, r.l}
}
func (r rgb1bm) rgb6() rgb6 {
	d := int32(31)
	if r.l {
		d = 63
	}
	return rgb6{r.r * d, r.g * d, r.b * d}
}
func (r rgb1b) ccolor() (ct.Color, bool) { return rgb1_c[rgb1{r.r, r.g, r.b}], r.l }
func (r rgb1b) rgb6() rgb6 {
	d := int32(31)
	if r.l {
		d = 63
	}
	return rgb6{r.r * d, r.g * d, r.b * d}
}
func (r rgb2) rgb1b() rgb1b {
	return rgb1b{(r.r + 1) / 2, (r.g + 1) / 2, (r.b + 1) / 2, (r.r+r.g+r.b) > 3 || r.r > 1 || r.g > 1 || r.b > 1}
}
func (r rgb6) rgb2() rgb2         { return rgb2{(r.r + 16) / 31, (r.g + 16) / 31, (r.b + 16) / 31} }
func (r rgb6) sub(d rgb6) rgb6    { return rgb6{r.r - d.r, r.g - d.g, r.b - d.b} }
func (r rgb6) times(t int32) rgb6 { return rgb6{r.r * t, r.g * t, r.b * t} }
func (r rgb6) avg(d rgb6) rgb6    { return rgb6{(r.r + d.r) / 2, (r.g + d.g) / 2, (r.b + d.b) / 2} }

func GetConsoleColor(r, g, b uint32) (m int32, fg ct.Color, fgl bool, bg ct.Color, bgl bool) {
	// 实色
	h := rgb6{int32(r), int32(g), int32(b)}
	// 背景色使用最接近的颜色
	j := h.rgb2().rgb1b().rgb6()
	bg, bgl = j.rgb2().rgb1b().ccolor()
	// 仿色:
	// 统计最高混合比: 字符'M' 50% (可能浮动) 最低 字符' ' 0%
	// 理论可用率66% 最高精确度 6位色
	// 设: 实色 = h , 底色 = j , 补色 = f , 表色 = l , 补色混合比 = m
	// 则: h = (f + j) / 2
	//     f = 2h - j
	//     f ≈ l * m
	//     l = [rgb2(0,0,0),rgb2(2,2,2)]
	//     m = [0,63]
	// f := h.times(2).sub(j.rgb6())
	// 纯色优化
	if h == j {
		m = 63
		return
	}
	min, minV := int32(0x7fffffff), rgb1bm{}
cop:
	for v, _ := range rgb1_c {
		for m = int32(0); m < 31; m++ {
			n := 31 - m // 差距越小灰度索引越大 32 = 64*50% 最高混合比
			for _, i := range []int32{31, 63} {
				f := (rgb6{v.r * i * n / 63, v.g * i * n / 63, v.b * i * n / 63}).avg(j)
				sr, sg, sb := (h.r - f.r), (h.g - f.g), (h.b - f.b)
				sub := sr*sr + sg*sg + sb*sb // 差方
				// fmt.Printf("h=%v f=%v v=%v i=%v srgb=(%v,%v,%v) m=%v n=%v sub=%v\n", h, f, v, i, sr, sg, sb, m, n, sub)
				// 完美混合
				if sub == 0 {
					minV = rgb1bm{v.r, v.g, v.b, i > 31, m}
					break cop // 249,240,225 -> 62 60 56 -> 63 63 63 -(-)> -1 -3 -7 -> 0 31 31 * 13 / 31
				}
				if min > sub {
					min = sub
					minV = rgb1bm{v.r, v.g, v.b, i > 31, m}
				}
			}
		}
	}
	// fmt.Println(h, j, minV)
	m = minV.m * 2
	fg, fgl = minV.rgb1b().ccolor()
	return
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
