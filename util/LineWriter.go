package util

/*
import (
	"io"
	"os"
	"sync"
	"time"
)

var (
	poollock = &sync.Mutex{}
	pool     = make(map[LineWriter]*os.File)
	poollife = make(map[LineWriter]time.Time)
)

type LineWriter string

func (this LineWriter) WriteLine(str string) (err error) {
	// TODO: replace with channal
	poollock.Lock()
	defer poollock.Unlock()

	w := pool[this]
	if w == nil {
		w, err = os.OpenFile(string(this), os.O_APPEND, os.ModeAppend)
		pool[this] = w
		poollife[this] = time.Now()
	}
	if err != nil {
		return
	}
	_, err = io.WriteString(w, str+"\n")
	poollife[this] = time.Now()
	return
}

func init() {
	// file handle GC 定期释放文件句柄
	go func() {
		for {
			poollock.Lock()
			for i, v := range poollife {
				// 写入时间大于5分钟
				if v.Before(time.Now().Add(-5 * time.Minute)) {
					pool[i].Close()
					delete(pool, i)
					delete(poollife, i)
				}
			}
			poollock.Unlock()
			time.Sleep(time.Minute)
		}
	}()
}
*/
