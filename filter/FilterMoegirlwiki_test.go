package filter

import (
	"fmt"
	"github.com/pa001024/MoeCron/source"
	"testing"
)

const (
	ex_f2 = `
		{
			"type":"moegirlwiki",
			"wiki_url":"http://zh.wikipedia.org/wiki/"
		}
		`
)

func TestMoegirl(t *testing.T) {
	f := New("name", []byte(ex_f2))
	src := []*source.FeedInfo{
		&source.FeedInfo{
			Link: "http://wwww",
			Content: `{{Infobox football biography
|name=高本詞史
|image=
|fullname=高本詞史
|birth_date={{Birth date and age|1967|12|31}}
|cityofbirth=[[宮崎縣]]
|countryofbirth=[[日本]]
|position=[[後衛 (足球)|後衛]]
|currentclub=
|years=?<br/>1993<br/>1994-1996
|clubs=[[札幌岡薩多|東芝]]<br/>[[名古屋鯨魚]]<br/>[[京都不死鳥]]
|caps(goals)=
|nationalyears=
|nationalteam=
|nationalcaps(goals)=
}}
'''高本詞史'''（{{Bd|1967年|12月31日||}}），前[[日本]][[足球]][[運動員]]。

==成績==
{{Football player club statistics 1|YYN}}
{{Football player club statistics 2|JPN|YYN}}
|-
|1993||[[名古屋鯨魚]]||[[日本職業足球甲級聯賽|甲級]]||9||0||0||0||1||0||10||0
|-
|1994||rowspan="3"|[[京都不死鳥]]||rowspan="2"|[[Japan Football League|Football League]]||8||1||0||0||colspan="2"|-||8||1
|-
|1995||10||2||1||0||colspan="2"|-||11||2
|-
|1996||[[日本職業足球甲級聯賽|甲級]]||3||0||0||0||0||0||3||0
{{Football player club statistics 3|1|JPN}}30||3||1||0||1||0||32||3
{{Football player club statistics 5}}30||3||1||0||1||0||32||3
|}

{{Japan-sport-bio-stub}}
{{Soccer-player-stub}}
[[Category:日本足球運動員]]`,
		},
	}
	ns := f.Process(src)
	fmt.Printf("%#v\n", src[0])
	fmt.Printf("%#v\n", ns[0])
}
