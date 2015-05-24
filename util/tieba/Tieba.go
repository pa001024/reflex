package tieba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/pa001024/reflex/util"
)

// 常量
const (
	_CLIENT_TYPE    = "_client_type=2"              // 客户端类型(安卓)
	_CLIENT_VERSION = "_client_version=5.1.3"       // 版本
	_PHONE_IMEI     = "_phone_imei=321694623354592" // 手机识别号
	_CUID           = "cuid=99A2533A10E616A5629FC5D341B58736|896958336496733"
	_ANONYMOUS      = "anonymous=0"
	_FROM           = "from=1317a"
	_NET_TYPE       = "net_type=3" // 网络类型
	_FID_URL        = "http://m.tieba.com/f?kw="
	_TBS_URL        = "http://tieba.baidu.com/dc/common/tbs"
	_USERINFO_URL   = "http://tieba.baidu.com/f/user/json_userinfo"
	_LOGIN_URL      = "http://c.tieba.baidu.com/c/s/login"          // 登录
	_THREAD_URL     = "http://c.tieba.baidu.com/c/c/thread/add"     // 发帖
	_POST_URL       = "http://c.tieba.baidu.com/c/c/post/add"       // 回帖
	_MYLIKE_URL     = "http://c.tieba.baidu.com/c/f/forum/favolike" // 获取喜欢
	_SIGN_URL       = "http://c.tieba.baidu.com/c/c/forum/sign"     // 签到
)

type Tieba struct {
	Uid      uint32
	Username string
	Password string
	BDUSS    string
	// 瞬态
	Liked        []Forum
	tbs          string
	client_id    string
	client       *http.Client
	VcodeHandler func(img image.Image) (vcode string)
}

func NewTieba(username, password string) (this *Tieba) {
	jar, _ := cookiejar.New(nil)
	this = &Tieba{
		Username:  username,
		Password:  password,
		client:    &http.Client{nil, nil, jar},
		client_id: "wappc_138" + fmt.Sprint(util.Random(10), "_", util.Random(3)),
	}
	return
}

func (this *Tieba) getSign(parm []string) (sign string) {
	sign = util.Md5String(fmt.Sprint(strings.Join(parm, ""), "tiebaclient!!!"))
	return
}

func (this *Tieba) checkTbs() {
	if this.tbs == "" {
		this.tbs = this.GetTbs()
	}
}

// 检查登录是否有效
func (this *Tieba) GetTbs() (tbs string) {
	res, err := this.post(_TBS_URL, nil, "BDUSS="+this.BDUSS+";")
	if err != nil {
		util.ERROR.Log(err)
		return
	}
	defer res.Body.Close()
	v := &TbsResult{}
	json.NewDecoder(res.Body).Decode(v)
	return v.Tbs
}

func (this *Tieba) _CLIENT_ID() string { return "_client_id=" + this.client_id }
func (this *Tieba) _BDUSS() string     { return "BDUSS=" + this.BDUSS }
func (this *Tieba) _TBS() string       { return "tbs=" + this.tbs }
func (this *Tieba) _COMMON(a []string) []string {
	return append(a, this._CLIENT_ID(), _CLIENT_TYPE, _CLIENT_VERSION, _PHONE_IMEI, _NET_TYPE)
}

func (this *Tieba) getLikedOnePage(pn int) (result *LikedResult, err error) {
	parm := []string{
		this._BDUSS(),
		this._CLIENT_ID(),
		_CLIENT_TYPE,
		_CLIENT_VERSION,
		_PHONE_IMEI,
		_NET_TYPE,
		fmt.Sprint("pn=", pn),
	}
	str := fmt.Sprint(strings.Join(parm, "&"), "&sign=", this.getSign(parm))
	req := bytes.NewBufferString(str)
	res, err := this.post(_MYLIKE_URL, req, "")
	bin := util.MustReadAll(res.Body)
	util.DEBUG.Log("[getLikedOnePage]", str, "\n========\n", string(bin))
	result = &LikedResult{}
	err = json.Unmarshal(bin, result)
	return
}

// 获取喜欢的贴吧
func (this *Tieba) GetLiked() (list []Forum) {
	result, _ := this.getLikedOnePage(1)
	list = result.ForumList
	for i := 2; result.Page.HasMore != 0; i++ {
		result, _ = this.getLikedOnePage(i)
		list = append(list, result.ForumList...)
	}
	this.Liked = list
	return
}
