package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	IP string
)

func GetIP() string {
	if IP != "" {
		return IP
	}
	IP, _ = CheckIP()
	return IP
}

func CheckIP() (ip string, err error) {
	res, err := http.Get("http://checkip.dyndns.com/")
	if err != nil {
		panic("CheckIP: App is Offine! Dead.")
	}
	bin, err := ioutil.ReadAll(res.Body)
	str := string(bin)
	if len(str) < 92 {
		panic("CheckIP: Bad Response!")
	}
	ip = str[76 : len(str)-14]
	return
}

func FetchImageAsStream(url string) (r io.Reader) {
	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)
	return buf
}
