package filter

import (
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
	"net/url"
	"regexp"
	"strings"
)

type FilterMoegirlwiki struct { // 萌娘百科 wikitext过滤器
	IFilter
	Filter

	WikiUrl string `json:"wiki_url"`
}

var (
	rep_mw_round1 = regexp.MustCompile(
		`^[\s\S]+?\s==? ?(?:基本介绍|简介|簡介|.{2,4}设定|.{2,4}介绍) ?==\n|` + // 去除第一个区块前的所有内容 替换为空 如果没有任何区块则不会替换
			`{{(?:[Bb][Dd]|[Ll]ang[-\|].+?)\|(.+?)}}|` + // 替换生日/lang为明文 $1
			// `{{.+}}|{{[^{]+?}}|` + // 去除所有{{}}标签 - 因为Go的正则不支持平衡组(也就.NET支持了)已经换成util.RemoveBlock()
			`{\|[\s\S]+?\|}|` + // 去除所有{||}标签
			`\[\[(\s*[^|\]]+?)\|.+?\]\]|` + // 替换[[明文|词条名]]词条链接为明文 即$2
			`\[\[([^[|]+?)\]\]|` + // 替换[[明文]]词条链接为明文 $3
			`<s>.+?</s>|` + // 去除所有<s></s>标签
			`<ref.*?>.+?</ref>|` + // 去除所有<ref></ref>标签
			`<\w+>(.+?)</\w+>` + // 替换所有HTML标签为明文 即$4
			``)
	rep_mw_round2 = regexp.MustCompile(
		`(?:\n|^).{1,12}[：:].+|` + // 去除说明体文字 如 姓名: xxx  `\n(?:本名|姓名|身高|机体|声优|种别|种类|萌点|年龄)[：:].+` +
			`''+(.+?)''+|` + // 替换粗体为明文 即$1
			`(?:\n|^)\|.+|` + // 去除表格
			`== (?:基本信息|登场作品|人物档案|人物信息|外部链接) ==[\s\S]+?(==|$)|` + // 去除无用区块 $2
			`==+? ?.+? ?==+|` + // 去除标题
			`[\*;] ?|` + // 去除列表
			`\[\[(\s*[^\|]+?)[\s\S]+?\]\]|` + // 替换[[明文|词条名]]词条链接为明文 即$3
			`\[\[([^\[]+?)\]\]|` + // 替换[[明文]]词条链接为明文 $4
			`<[A-z/].+?>` + // 去除HTML标签如<br /><references />
			``)
	rep_mw_round3 = regexp.MustCompile(
		`===+ (?:基本资料) ===+[\s\S]+?(==|$)|` + // 去除无用小区块 $1
			`(\n)\n+|` + // 去除连续两个以上的换行符 $2
			`\|\|+|` + // 去除连续两个以上的|
			`^\s+|` + // 去除开头空白
			`<[A-z/].+?>` + // 去除HTML标签如<br /><references /> (再来一发)
			`（待补完）|` + // 去除（待补完）
			`[^\S\n]\s+` + // 去除两个以上的空白符 行文中分词的单个空格不会被替换 中文空格不会被替换
			``)

	// {{[\s\S]+?}}|\[\[[\s\S]+?\]\]|'''(.+?)'''|== (.+?) ==|<references />|\s+
)

func (this *FilterMoegirlwiki) FilterContent(src string) (dst string) {
	dst = rep_mw_round1.ReplaceAllString(src, "$1$2$3$4")
	dst = util.RemoveBlock(dst, "{{", "}}")    // 去除所有{{}}标签
	dst = util.RemoveBlock(dst, "<!--", "-->") // 去除所有HTML注释
	dst = rep_mw_round2.ReplaceAllString(dst, "$1$2$3$4")
	dst = rep_mw_round3.ReplaceAllString(dst, "$1$2")
	return dst
}

func (this *FilterMoegirlwiki) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	dst = make([]*source.FeedInfo, 0, len(src))
	for _, v := range src {
		nv := *v
		nv.Content = this.FilterContent(nv.Content)
		nv.Link = this.WikiUrl + url.QueryEscape(strings.Replace(nv.Title, " ", "_", -1)) // TODO: 可能不工作
		dst = append(dst, &nv)
	}
	return
}
