package webqq

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"

	asc "github.com/pa001024/MoeWorker/util/ascgen"
)

const CAPTCHA_URL = "http://captcha.qq.com/"

// 获取验证码
func (this *WebQQ) ptlogin_getimage(vc_type string) (code string, err error) {
	res, err := this.client.Get(CAPTCHA_URL + "getimage?" + (url.Values{"uin": {this.Id}, "aid": {WEBQQ_APPID}, "r": {fmt.Sprint(rand.ExpFloat64())}, "vc_type": {vc_type}}).Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()
	rw, err := os.Create("vc.jpg")
	if err != nil {
		return
	}
	io.Copy(rw, res.Body)
	rw.Seek(0, 0)
	asc.ShowFile(os.Stdout, rw, asc.Console{6, 14, 120}, true, false) // TODO: 缺乏配置
	rw.Close()
	fmt.Print("Enter Verify Code: ")
	fmt.Scanf("%s", &code)
	return
}