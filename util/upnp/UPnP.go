package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/pa001024/MoeWorker/util"
)

var (
	search_msg = []byte("M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 3\r\nST: upnp:rootdevice\r\n\r\n")
	mcast_addr = &net.UDPAddr{net.ParseIP("239.255.255.250"), 1900, ""}
)

type UPnPService struct {
	LocalIP    net.IP
	LocalIfi   net.Interface
	ServerHost string
}

// 创建新UPnP控制器
func NewUPnPService() (this *UPnPService) {
	this = &UPnPService{}
	this.GetLocalIP(0)
	return
}

// 获取内网IP
func (this *UPnPService) GetLocalIP(skip int) {
	if is, err := net.Interfaces(); err == nil {
		for _, v := range is {
			if as, err := v.Addrs(); err == nil {
				for _, v2 := range as {
					l := net.ParseIP(v2.String())
					if !l.Equal(net.IPv4zero) && skip <= 0 {
						this.LocalIP = l
						this.LocalIfi = v
						return
					}
					skip--
				}
			}
		}
	}
}

/*
 获取 UPnP 服务器
 1. 对 239.255.255.250:1900 进行 HTTP over UDP 广播
 M-SEARCH * HTTP/1.1
 HOST: 239.255.255.250:1900
 MAN: "ssdp:discover"
 MX: 3
 ST: upnp:rootdevice

 2. 等待MX秒 读取回包 (302 跳转)
 HTTP/1.1 200 OK
 CACHE-CONTROL: max-age=100
 DATE: Fri, 09 Aug 2013 08:40:47 GMT
 EXT:
 LOCATION: http://192.168.1.253:1900/igd.xml
 SERVER: Wireless AP WR700N, UPnP/1.0
 ST: upnp:rootdevice
 USN: uuid:upnp-InternetGatewayDevice-192168125378900001::upnp:rootdevice

 3. 跟踪LOCATION到Host 可直接解析
*/
func (this *UPnPService) GetUPnPServer() (err error) {
	util.DEBUG.Log("listening...")
	conn, err := net.ListenMulticastUDP("udp4", &this.LocalIfi, mcast_addr)
	util.DEBUG.Log("writing search_msg...")
	_, err = conn.WriteTo(search_msg, mcast_addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	util.DEBUG.Log("sleeping...")
	time.Sleep(3 * time.Second) // 按MX设置的3秒响应时间等待服务器响应
	util.DEBUG.Log("reading...")
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	s := string(buf[:n])
	util.DEBUG.Log("\n[response]\n", s)
	ms := regexp.MustCompile(`LOCATION:\s*(\S+)`).FindStringSubmatch(s)
	if ms != nil {
		u, _ := url.Parse(ms[1])
		this.ServerHost = u.Host
		util.DEBUG.Log("set ServerHost to", u.Host)
	}
	return
}

type AddPortMapping struct {
	xml.Name                  `xml:"u:AddPortMapping"`
	Xmlns                     string `xml:"xmlns:u,attr"`
	NewRemoteHost             string // 外网IP 一般为空
	NewExternalPort           int    // 外网端口
	NewProtocol               string // TCP or UDP
	NewInternalPort           int    // 内网端口
	NewInternalClient         string // 内网IP
	NewEnabled                int    // 启用=1 禁用=0
	NewPortMappingDescription string // 描述
	NewLeaseDuration          int    // 结束时间 0 为无限
}

/*
映射端口

POST /ipc HTTP/1.1
HOST: 192.168.1.253:1900
Content-Type: text/xml; charset="utf-8"
Content-Length: 615
SOAPACTION: "urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"

<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<s:Body>
<u:AddPortMapping xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewRemoteHost></NewRemoteHost>
<NewExternalPort>5678</NewExternalPort>
<NewProtocol>TCP</NewProtocol>
<NewInternalPort>3389</NewInternalPort>
<NewInternalClient>192.168.1.100</NewInternalClient>
<NewEnabled>1</NewEnabled>
<NewPortMappingDescription>xxxxxxxx</NewPortMappingDescription>
<NewLeaseDuration>0</NewLeaseDuration>
</u:AddPortMapping>
</s:Body>
</s:Envelope>
*/
func (this *UPnPService) Link(v *AddPortMapping) (err error) {
	server_url := "http://" + this.ServerHost + "/ipc"
	util.DEBUG.Logf("linking %s : %+v", server_url, v)
	buf := bytes.NewBufferString("<?xml version=\"1.0\"?>\r\n")
	buf.WriteString("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\">\r\n")
	buf.WriteString("<s:Body>\r\n")
	buf.WriteString("<u:AddPortMapping xmlns:u=\"urn:schemas-upnp-org:service:WANIPConnection:1\">\r\n")
	buf.WriteString("<NewRemoteHost>")
	xml.Escape(buf, []byte(v.NewRemoteHost))
	buf.WriteString("</NewRemoteHost>\r\n")
	buf.WriteString("<NewExternalPort>" + util.ToString(v.NewExternalPort) + "</NewExternalPort>\r\n")
	buf.WriteString("<NewProtocol>")
	xml.Escape(buf, []byte(v.NewProtocol))
	buf.WriteString("</NewProtocol>\r\n")
	buf.WriteString("<NewInternalPort>" + util.ToString(v.NewInternalPort) + "</NewInternalPort>\r\n")
	buf.WriteString("<NewInternalClient>")
	xml.Escape(buf, []byte(v.NewInternalClient))
	buf.WriteString("</NewInternalClient>\r\n")
	buf.WriteString("<NewEnabled>" + util.ToString(v.NewEnabled) + "</NewEnabled>\r\n")
	buf.WriteString("<NewPortMappingDescription>")
	xml.Escape(buf, []byte(v.NewPortMappingDescription))
	buf.WriteString("</NewPortMappingDescription>\r\n")
	buf.WriteString("<NewLeaseDuration>" + util.ToString(v.NewLeaseDuration) + "</NewLeaseDuration>\r\n")
	buf.WriteString("</u:AddPortMapping>\r\n</s:Body>\r\n</s:Envelope>\r\n\r\n")
	res, err := http.Post(server_url, "text/xml", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	util.DEBUG.Log("\n[response]\n", string(b))
	return
}
func (this *UPnPService) LinkTCP(local_addr, remote_addr string, local_port, remote_port int, desc string) (err error) {
	return this.Link(&AddPortMapping{xml.Name{}, "urn:schemas-upnp-org:service:WANIPConnection:1", remote_addr, remote_port, "TCP", local_port, local_addr, 1, desc, 0})

}
func (this *UPnPService) LinkUDP(local_addr, remote_addr string, local_port, remote_port int, desc string) (err error) {
	return this.Link(&AddPortMapping{xml.Name{}, "urn:schemas-upnp-org:service:WANIPConnection:1", remote_addr, remote_port, "UDP", local_port, local_addr, 1, desc, 0})
}

// 临时的
func main() {
	util.DEBUG.SetEnable(true)
	u := NewUPnPService()
	if err := u.GetUPnPServer(); err == nil {
		u.LinkTCP(u.LocalIP.String(), "", 25565, 25565, "minecraft")
	}
}
