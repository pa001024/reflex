package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

// 获取当前内网IP
func GetIPLocal() []string {
	res := make([]string, 0, 4)
	ips, _ := net.InterfaceAddrs()
	for _, ip := range ips {
		if ipnet, ok := ip.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.String()[:7] != "169.254" {
				res = append(res, ipnet.IP.String())
			}
		}
	}
	return res
}

// 当前IP的缓存
var IP string

// 获取当前IP, 有缓存
func GetIP() string {
	if IP != "" {
		return IP
	}
	IP, _ = CheckIP()
	return IP
}

// IP信息
type IPInfo struct {
	IP          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
}

// 获取当前IP信息
func GetIPInfo() *IPInfo {
	res, err := http.Get("http://whois.pconline.com.cn/ipJson.jsp")
	if err != nil {
		return nil
	}
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	str := string(bin)
	v := &IPInfo{}
	ex := regexp.MustCompile(`\{"ip".+?\}`)
	bin = []byte(ex.FindString(str))
	json.Unmarshal(bin, v)
	return v
}

// 获取当前IP, 无缓存
func FlushIP() string {
	IP, _ = CheckIP()
	return IP
}

// 获取当前IP
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

// 下载图片(带缓存)
func FetchImageAsStream(url string) (r *bytes.Buffer) {
	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)
	return buf
}
