package util

import (
	"encoding/json"
	"io"
)

type Unique interface { // 禁止重复
	Put(id string) bool    // 检测唯一并设置
	CanPut(id string) bool // 检测唯一
}

func NewUnrepeatable(typ string, args map[string]string) (this Unique) { // 工厂函数
	switch typ {
	default:
	case "file_mem":
		data, _ := json.Marshal(args)
		this = &UnrepeatableFileMem{}
		json.Unmarshal(data, this)
		break
	}
	return
}

type UnrepeatableFileMem struct { // 使用内存HashSet来保证唯一性(非进程安全)
	Unique

	CacheContainer       map[string]bool
	PersistenceContainer string `json:"stonefile"`
}

func (this *UnrepeatableFileMem) Put(id string) bool {
	if this.CanPut(id) {
		_, err := io.WriteString(this.PersistenceContainer, id+"\n")
		if err != nil {
			Log("Fail to Write File:", err)
			return false
		}
		this.CacheContainer[id] = true
	}
	return false
}
func (this *UnrepeatableFileMem) CanPut(id string) bool {
	if this.CacheContainer[id] {
		return false
	}
	return true
}
