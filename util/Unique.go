package util

import (
	"encoding/json"
)

type Unique interface { // 禁止重复更新
	Put(id string) (err error) // 检测唯一并设置
	CanPut(id string) bool     // 检测唯一
}

func NewUnique(typ string, args map[string]string) (this Unique) { // 工厂函数
	switch typ {
	default:
	// TODO: case "sql_mem":
	case "file_mem":
		data, _ := json.Marshal(args)
		this = &UniqueFileMem{}
		json.Unmarshal(data, this)
		break
	}
	return
}

type UniqueFileMem struct { // 使用内存HashSet来保证唯一性(非进程安全)
	Unique

	containerCache       map[string]bool // 缓存容器
	containerPersistence LineWriter      `json:"stonefile"` // 持久容器
}

func (this *UniqueFileMem) Put(id string) (err error) {
	if this.CanPut(id) {
		err = this.containerPersistence.WriteLine(id)
		if err != nil {
			Log("Fail to Write File:", err)
			return
		}
		this.containerCache[id] = true
	}
	return
}
func (this *UniqueFileMem) CanPut(id string) bool {
	if this.containerCache[id] {
		return false
	}
	return true
}
