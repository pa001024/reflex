package webqq

import (
	"encoding/json"
)

// 正常结构
type ContentM struct {
	Msg  []string
	Font Font
}

// 变态结构 ["msg",["font",{...}]]
type Content []interface{}

// ["阿斯达\n阿斯达\n阿斯达",["font",{"name":"宋体","size":"10","style":[0,0,0],"color":"000000"}]]
type Font struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Style []int  `json:"style"`
	Color string `json:"color"`
}

// 提取字体
func (this Content) Font() (v Font) {
	for _, v1 := range this {
		if v2, ok := v1.([]interface{}); ok {
			for _, v3 := range v2 {
				if v4, ok := v3.(map[string]interface{}); ok {
					v.Name, _ = v4["name"].(string)
					v.Size, _ = v4["size"].(string)
					v.Style, _ = v4["style"].([]int)
					v.Color, _ = v4["color"].(string)
					return
				}
			}
		}
	}
	return
}

// 提取消息
func (this Content) Msg() (v []string) {
	v = make([]string, 0, len(this)-1)
	for _, v1 := range this {
		if v2, ok := v1.(string); ok {
			v = append(v, v2)
		}
	}
	return
}

// 转化到string
func (this Content) EncodeString() string {
	b, _ := json.Marshal(this)
	return string(b)
}

// 解码
func (this Content) Decode() (v ContentM) {
	v.Msg = make([]string, 0, len(this)-1)
	for _, v1 := range this {
		switch v2 := v1.(type) {
		case string:
			v.Msg = append(v.Msg, v2)
		case []interface{}:
			for _, v3 := range v2 {
				if v4, ok := v3.(map[string]interface{}); ok {
					v.Font.Name, _ = v4["name"].(string)
					v.Font.Size, _ = v4["size"].(string)
					v.Font.Style, _ = v4["style"].([]int)
					v.Font.Color, _ = v4["color"].(string)
					break
				}
			}
		}
	}
	return
}

// 转换到下层对象
func (this ContentM) Encode() (v Content) {
	v = make([]interface{}, 0, len(this.Msg)+1)
	v = append(v, []interface{}{"font", this.Font})
	for _, v1 := range this.Msg {
		v = append(v, v1)
	}
	return v
}

func (this ContentM) String() string {
	return this.Encode().EncodeString()
}
