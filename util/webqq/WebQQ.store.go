package webqq

import (
	// "encoding/json"

	"code.google.com/p/leveldb-go/leveldb"
	"github.com/pa001024/reflex/util"
)

// 储存好友信息的容器
type FriendJar interface {
	GetByUin(uin Uin) *FriendInfo
	GetByAccount(acc Account) *FriendInfo
	All() []*FriendInfo
	Put(model *FriendInfo)
	Detele(uin Uin)
}

// 储存群信息的容器
type GroupJar interface {
	GetByGin(gin GroupId) *GroupInfo
	GetByGroupId(gid GroupId) *GroupInfo
	All() []*GroupInfo
	put(model *GroupInfo)
	Detele(gin GCode)
}

// 储存自定义表情信息的容器
type CustomFaceJar interface {
	GetByGin()
}

// 储存离线图片信息的容器
type OfflinePicJar interface {
	GetByGin()
}

// the implementation of FriendJar
type LevelDbFriendJar struct {
	FriendJar
	db *leveldb.DB
}

func NewLevelDbFriendJar(filename string) (this *LevelDbFriendJar) {
	this = &LevelDbFriendJar{}
	db, err := leveldb.Open(filename, nil)
	util.Try(err)
	this.db = db
	return this
}
func (this *LevelDbFriendJar) Get() {
}

func (this *LevelDbFriendJar) Put() {
}
