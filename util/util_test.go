package util

import (
	"fmt"
	"math/rand"
	"testing"
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
