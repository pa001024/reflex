package smtp

import (
	"crypto/tls"
	"net/smtp"
	"reflect"
	"strings"
	"unsafe"

	"github.com/pa001024/reflex/util"
)

// 邮件发送
type SMTPClient struct {
	User      string // user@exmaple.com
	Server    string // smtp.exmaple.com:25 or ssl smtp.exmaple.com:465
	EnableSSL bool   // 465 port ssl conn
	auth      smtp.Auth
}

// 初始化
func NewSMTPClient(user, pass string, ssl bool, server string) (this *SMTPClient) {
	this = &SMTPClient{
		user, server, ssl, nil,
	}
	if pass != "" {
		this.auth = smtp.PlainAuth("", user, pass, server[:strings.Index(server, ":")])
	}
	return
}

// SSL方式连接服务器 (与STARTTLS方式不同 适用于465端口)
func (this *SMTPClient) sendSSL(to []string, title, context string) (err error) {
	defer util.Catch(&err)
	host := this.Server[:strings.Index(this.Server, ":")]
	conn, err := tls.Dial("tcp", this.Server, nil)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	*(*bool)(unsafe.Pointer(reflect.ValueOf(c).Elem().FieldByName("tls").UnsafeAddr())) = true // set tls to true for AUTH

	if this.auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(this.auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(this.User); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	body := []byte(strings.Join([]string{
		"From: " + this.User,
		"To: " + strings.Join(to, ","),
		"Subject: " + title,
		"Content-Type: text/html; charset=utf-8;",
		"",
		context,
	}, "\r\n"))
	_, err = w.Write(body)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

// 同步发送
func (this *SMTPClient) Send(to []string, title, context string) (err error) {
	if this.EnableSSL {
		return this.sendSSL(to, title, context)
	}
	body := []byte(strings.Join([]string{
		"From: " + this.User,
		"To: " + strings.Join(to, ","),
		"Subject: " + title,
		"Content-Type: text/html; charset=utf-8;",
		"",
		context,
	}, "\r\n"))
	err = smtp.SendMail(
		this.Server,
		this.auth,
		this.User,
		to,
		body)
	return
}
