package filter

import (
	"regexp"
)

type FilterMoegirlwiki struct { // 萌娘百科
	IFilter
}

var (
	reg_pre = regexp.MustCompile(
		`^[\s\S]+?\s== ?.+ ?==\n|` + // 去除第一个区块前的所有内容 替换为空 如果没有任何区块则不会替换
			`{{[\s\S]+?}}|` + // 去除所有{{}}标签
			`\[\[(?:File|分类):[\s\S]+?\]\]|` + // 提取并替换所有WIKI标签(提取图片和分类)
			`\[\[([^|\s]+?)\]\]|` + // 替换[[明文|词条名]]词条链接为明文 即$1
			`\[\[[\s\S]+?\]\]|` + // 替换[[明文]]词条链接为明文
			`'''(.+?)'''|` + // 替换粗体为明文 即$2
			`\n.{1,12}[：:].+` + // 去除说明体文字 如 姓名: xxx
			// `\n(?:本名|姓名|身高|机体|声优|种别|种类|萌点|年龄)[：:].+` +
			``)
	reg_content = regexp.MustCompile(
		`== (?:基本信息|登场作品) ==[\s\S]+?(==)|` + // 去除基本信息和登场作品区块
			`===? ?(.+?) ?===?|` + // 替换小标题为明文 即$1
			`<s>.+</s>|` + // 去除所有<s></s>标签
			`\* ?|` + // 去除列表
			`<[\w "%%=/\-:]+>` + // 去除HTML标签如<br /><references />
			``)
	reg_final = regexp.MustCompile(
		`^\s+|` + // 去除开头空白
			`(\n)\n+|` + // 去除连续两个以上的换行符
			`[^\S\n]\s+` + // 去除两个以上的空白符 行文中分词的单个空格不会被替换 中文空格不会被替换
			``)

	// {{[\s\S]+?}}|\[\[[\s\S]+?\]\]|'''(.+?)'''|== (.+?) ==|<references />|\s+
)

func (this *FilterMoegirlwiki) Process(src string) (dst string) {
	dst = reg_final.ReplaceAllString(reg_content.ReplaceAllString(reg_pre.ReplaceAllString(src, "$1$2"), "$1"), "")
	ds := []rune(dst)
	if len(ds) > 100 {
		dst = string(ds[0:100]) + "..."
	}
	return dst
}
