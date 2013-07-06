package target

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	tarsina = &SinaWeibo{}
)

func init() {
	json.Unmarshal([]byte(`
		{
			"name":"test",
			"type":"sina",
			"client_id":"236071xxxx",
			"redirect_uri":"http://xx.com/api/weibo/",
			"access_token":"2.006dgysDaxSlZCd6a36b15dbExxxxx",
			"method":[
				{
					"action":"update",
					"source":["wiki"],
					"filter":["moewikitext","weibobase"],
					"limit":"1"
				}
			]
		}
		`), tarsina)
}

func TestSinaWeiboUpdate(t *testing.T) {
	t.Log("[单元测试]target.SinaWeibo.Update()")
	fmt.Printf("%#v\n", tarsina)
	t.Log("[]Authorize")
	fmt.Println(tarsina.Authorize())
	tarsina.Update("[单元测试]target.SinaWeibo.Update()")
}
