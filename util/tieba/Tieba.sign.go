package tieba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pa001024/reflex/util"
)

type SignResult struct {
	ErrorCode int32  `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Time      int32  `json:"time"`
	Ctime     int32  `json:"ctime"`
	Logid     int32  `json:"logid"`
}

func (this *Tieba) Sign(kw string) (result *Result, err error) {
	this.checkTbs()
	parm := []string{
		this._BDUSS(),
		"kw=" + kw,
		this._TBS(),
	}
	str := fmt.Sprint(strings.Join(parm, "&"), "&sign=", this.getSign(parm))
	req := bytes.NewBufferString(str)
	res, err := this.post(_SIGN_URL, req, "")
	bin := util.MustReadAll(res.Body)
	util.DEBUG.Log("[Sign]", str, "\n========\n", string(bin))
	result = &Result{}
	err = json.Unmarshal(bin, result)
	return
}

func (this *Tieba) SignAll() {
	if this.Liked == nil || len(this.Liked) == 0 {
		this.GetLiked()
	}
	for _, v := range this.Liked {
		this.Sign(v.Name)
		time.Sleep(time.Second)
	}
}
