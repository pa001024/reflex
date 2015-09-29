package util

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	arrA = genArray()
	arrB = genArray()
)

func TestRemoveBlock(t *testing.T) {
	os := `{{
| {{xxx}}
* sss{{xx}}
}}

asdasd

{{
sdasdasd
s}}`
	fmt.Println(RemoveBlock(os, "{{", "}}"))
}

func genArray() (rst []string) {
	g := []rune(`abcdefgijklmnopqrstuvwxyzABCDEFGIJKLMNOPQRSTUVWXYZ1234567890-=,./;'[]<>?:"{}`)
	gl := len(g)
	rst = make([]string, 1000)
	for i, _ := range rst {
		for j := 0; j < 10+(rand.Int()%90); j++ {
			rst[i] += string(g[rand.Int()%gl])
		}
	}
	return
}

// Run benchmark with `go test -bench=Benchmark.+`
func BenchmarkIsArrayDuplicate(b *testing.B) {
	t := 0
	for i := 0; i < b.N; i++ {
		if IsArrayDuplicate(arrA, arrB) {
			t++
		}
	}
	b.Log("Matched:", t, "/", b.N)
}

func BenchmarkIsArrayDuplicateOpt(b *testing.B) {
	t := 0
	for i := 0; i < b.N; i++ {
		if IsArrayDuplicateOpt(arrA, arrB) {
			t++
		}
	}
	b.Log("Matched:", t, "/", b.N)
}

func TestTryCatch(t *testing.T) {
	defer Catch()
	Try(fmt.Errorf("test err %s", "[test]"))
}

func TestTryTimeoutCatch(t *testing.T) {
	defer Catch()
	dosomething := func() <-chan time.Time {
		return time.After(3 * time.Second)
	}
	select {
	case <-time.After(1 * time.Second):
		Try(fmt.Errorf("Exec Timeout %s", "Error"))
	case <-dosomething():
	}
}

func TestEncodeJsUint64BE(t *testing.T) {
	v := uint64(2735284921)
	r := EncodeJsUint64BE(v)
	fmt.Println(r)
	if r != `\x00\x00\x00\x00\xa3\x09\x22\xb9` {
		t.Fail()
	}
}

func TestUint64BEString(t *testing.T) {
	v := uint64(2735284921)
	r := EncodeJsHex([]byte(Uint64BEString(v)))
	fmt.Println(r)
	de := BEStringUint64(string(DecodeJsHex(`\x00\x00\x00\x00\xa3\x09\x22\xb9`)))
	if de != v || r != `\x00\x00\x00\x00\xa3\x09\x22\xb9` {
		t.Fail()
	}
}

func TestGetIPLocal(t *testing.T) {
	fmt.Println(GetIPLocal())
}

func TestAESCBCStringX(t *testing.T) {
	origintext := "123456"
	key := Md5("WTF?MASK")
	opwd := AESCBCStringX(key, []byte(origintext), true)
	if opwd != "3CCBEDF2AB9727A7D3CF80425F1EE936" {
		t.Log(opwd)
		t.Fail()
	}
	bin, _ := hex.DecodeString(opwd)
	rpwd := string(AESCBC(key, bin, false))
	if rpwd != "123456" {
		t.Log(rpwd)
		t.Fail()
	}
}

func TestJsCurrentTime(t *testing.T) {
	fmt.Println(JsCurrentTime())
}
