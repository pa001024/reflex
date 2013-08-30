package webqq

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/pa001024/MoeWorker/util"
)

// 正常结构
type ContentModel struct {
	Font     Font
	Message  string
	HasImage bool // 是否有图片
}

// 变态结构 [["font",{...}],["offpic",{...},],["cface",{...}],["face",100],"msg",["cface_idx",0]]
type Content []json.RawMessage

var (
	contentPatten = regexp.MustCompile(`[(offpic|cface|face)=(.+?)]`)
	DefaultFont   = Font{"微软雅黑", 10, []int{0, 0, 0}, "000000"} // 默认字体
	DefaultServer = "112.90.78.75:0"
)

// 字体 {"name":"宋体","size":"10","style":[0,0,0],"color":"000000"}
type Font struct {
	Name  string `json:"name"`
	Size  uint32 `json:"size"`
	Style []int  `json:"style"`
	Color string `json:"color"`
}

// 离线图片 success:1 file_path:/32d369de-e2a1-49b7-8b27-249d57bbd5d0
type Offpic struct {
	FilePath string `json:"file_path"`
	Success  int32  `json:"success"`
}

// 自定义表情 {"file_id":4086436741,"key":"j9v3yUNifMrByvWC","name":"{DC854AD1-A1DB-0855-E267-2CC14D14345A}.gif","server":"112.90.78.75:0"}
type Cface struct {
	Name   string `json:"name"`
	FileId string `json:"file_id"`
	Key    string `json:"key"`
	Server string `json:"server"`
}

// 转化到string
func (this Content) EncodeString() string {
	b, _ := json.Marshal(this)
	return string(b)
}

// 解码
func (this Content) Decode() (v ContentModel) {
	var index []string
	if len(this) > 3 { // msg pic picindex font
		index = make([]string, 0, len(this)-3)
	}
	for _, data := range this {
		var block interface{}
		json.Unmarshal(data, &block)
		switch v1 := block.(type) {
		case string:
			if v1 != "" {
				v.Message += v1
			}
		case []interface{}:
			var block []json.RawMessage
			json.Unmarshal(data, &block)
			var v2 string
			json.Unmarshal(block[0], &v2)
			switch v2 {
			case "font":
				json.Unmarshal(block[1], &v.Font)
			case "cface":
				v.HasImage = true
				cface := Cface{}
				json.Unmarshal(block[1], &cface)
				if index != nil {
					index = append(index, cface.Name)
				}
				v.Message += fmt.Sprintf("[cface=%s]", cface.Name)
			case "offpic":
				v.HasImage = true
				offpic := Offpic{}
				json.Unmarshal(block[1], &offpic)
				if index != nil {
					index = append(index, offpic.FilePath)
				}
				if offpic.Success != 1 { // 忽略裂图
					v.Message += fmt.Sprintf("[offpic=%s]", offpic.FilePath)
				}
			case "face":
				var face int32
				json.Unmarshal(block[1], &face)
				v.Message += fmt.Sprintf("[face=%s]", face)
			case "cface_idx":
				var cface_idx int32
				json.Unmarshal(block[1], &cface_idx)
				v.Message += fmt.Sprintf("[cface=%s]", index[cface_idx])
			case "offpic_idx":
				var offpic_idx int32
				json.Unmarshal(block[1], &offpic_idx)
				v.Message += fmt.Sprintf("[offpic=%s]", index[offpic_idx])
			}
		}
	}
	return
}

// 转换到下层对象
func (this ContentModel) Encode(webqq *WebQQ) (v Content) {
	matches := contentPatten.FindAllStringSubmatch(this.Message, -1)
	if matches == nil { // 没有图片
		v = make([]json.RawMessage, 1, 2)
		b, err := json.Marshal(this.Message)
		util.Try(err)
		v[0] = json.RawMessage(b)
	} else {
		v = make([]json.RawMessage, 0, len(matches)*2+2) // 123[]1[]2[]3 = 3*2+1 = 7
		for _, v1 := range matches {
			switch v1[1] {
			case "offpic":
				b, err := json.Marshal([]interface{}{v1[1], Offpic{Success: 1, FilePath: v1[2]}})
				if err == nil {
					v = append(v, b)
				}
			case "cface":
				b, err := json.Marshal([]interface{}{v1[1], Cface{Key: webqq.cface_key, Name: v1[2], Server: DefaultServer}})
				if err == nil {
					v = append(v, json.RawMessage(b))
				}
			case "face":
				b, err := json.Marshal([]interface{}{v1[1], util.ToInt(v1[2])})
				if err == nil {
					v = append(v, json.RawMessage(b))
				}
			}
		}
	}

	b, err := json.Marshal([]interface{}{"font", this.Font})
	util.Try(err)
	v = append(v, json.RawMessage(b))
	return v
}

// func (this ContentModel) String() string {
// 	return this.Encode().EncodeString()
// }
