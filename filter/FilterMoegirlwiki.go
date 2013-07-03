package filter

import (
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
	"net/url"
	"regexp"
)

type FilterMoegirlwiki struct { // 萌娘百科
	IFilter
	Filter

	WikiUrl string `json:"wiki_url"`
}

var (
	rep_mw_round1 = regexp.MustCompile(
		`^[\s\S]+?\s==? ?(?:基本介绍|简介|簡介|.{2,4}设定|.{2,4}介绍) ?==\n|` + // 去除第一个区块前的所有内容 替换为空 如果没有任何区块则不会替换
			`{{(?:Bd|lang[-\|].+?)\|(.+?)\|*}}|` + // 替换生日/lang为明文 $1
			// `{{.+}}|{{[^{]+?}}|` + // 去除所有{{}}标签 - 因为Go的正则不支持平衡组(泥马也就.NET支持了)已经换成util.RemoveBlock()
			`{\|[\s\S]+?\|}|` + // 去除所有{||}标签
			`\[\[([^|\s]+?)\]\]|` + // 替换[[明文|词条名]]词条链接为明文 即$2
			`\[\[[\s\S]+?\]\]|` + // 替换[[明文]]词条链接为明文
			`'''(.+?)'''|` + // 替换粗体为明文 即$3
			`<s>.+?</s>|` + // 去除所有<s></s>标签
			`<\w+>(.+?)</\w+>` + // 替换所有HTML标签为明文 即$4
			``)
	rep_mw_round2 = regexp.MustCompile(
		`(?:\n|^).{1,12}[：:].+|` + // 去除说明体文字 如 姓名: xxx  `\n(?:本名|姓名|身高|机体|声优|种别|种类|萌点|年龄)[：:].+` +
			`{{[\s\S]+?}}|` + // 去除所有{{}}标签 (再来一发)
			`{\|[\s\S]+?\|}|` + // 去除所有{||}标签 (再来一发)
			`(?:\n|^)\|.+|` + // 去除表格
			`== (?:基本信息|登场作品|人物档案|人物信息|外部链接) ==[\s\S]+?(==|$)|` + // 去除无用区块 $1
			`==+? ?.+? ?==+|` + // 去除标题
			`\* ?|` + // 去除列表
			`\[\[([^|\s]+?)\]\]|` + // 第二次替换[[明文|词条名]]词条链接为明文 即$2
			`<[\w "%%=/\-:]+?>` + // 去除HTML标签如<br /><references />
			``)
	rep_mw_round3 = regexp.MustCompile(
		`}}|` + // 去除剩余}}
			`===+ (?:基本资料) ===+[\s\S]+?(==|$)|` + // 去除无用小区块 $1
			`(\n)\n+|` + // 去除连续两个以上的换行符 $2
			`^\s+|` + // 去除开头空白
			`（待补完）|` + // 去除（待补完）
			`[^\S\n]\s+` + // 去除两个以上的空白符 行文中分词的单个空格不会被替换 中文空格不会被替换
			``)
	flp_mw = regexp.MustCompile(`R18`)

	// {{[\s\S]+?}}|\[\[[\s\S]+?\]\]|'''(.+?)'''|== (.+?) ==|<references />|\s+
)

func (this *FilterMoegirlwiki) FilterContent(src string) (dst string) {
	const LEN_LIMIT = 100
	dst = rep_mw_round3.ReplaceAllString(
		rep_mw_round2.ReplaceAllString(
			util.RemoveBlock(rep_mw_round1.ReplaceAllString(src, "$1$2$3$4"), "{{", "}}"),
			"$1$2"),
		"$1$2")
	ds := []rune(dst)
	if len(ds) > LEN_LIMIT {
		dst = string(ds[0:LEN_LIMIT])
	}
	return dst
}

func (this *FilterMoegirlwiki) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	dst = make([]*source.FeedInfo, 0, 10)
	for _, v := range src {
		if !flp_mw.MatchString(v.Title) && !flp_mw.MatchString(v.Content) {
			nv := *v
			nv.Content = this.FilterContent(nv.Content)
			nv.Link = this.WikiUrl + url.QueryEscape(nv.Title) // TODO: 可能不工作
			dst = append(dst, &nv)
		}
	}
	return
}
