package target

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	tarqq = &QQWeibo{}
)

func init() {
	json.Unmarshal([]byte(`
		{
			"type":"qq",
			"client_id":"236071xxxx",
			"redirect_uri":"http://xxxx.com/api/weibo/",
			"access_token":"2.006dgysDaxSlZCd6a36b15dbxxxxxx",
			"openid":"2.0xxxxxxxxxxxxxxxxxx",
			"method":[
				{
					"action":"update",
					"source":["wiki"],
					"filter":["moewikitext","url","weibobase","zhconv"],
					"limit":1
				}
			]
		}
		`), tarsina)
}

func TestQQWeiboUpdate(t *testing.T) {
	t.Log("[单元测试]target.QQWeibo.Update()")
	fmt.Printf("%#v\n", tarsina)
	t.Log("[]Authorize")
	fmt.Println(tarsina.Authorize())
	tarqq.Update("[单元测试]target.QQWeibo.Update()")
}
