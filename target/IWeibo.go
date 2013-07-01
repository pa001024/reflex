package target

import (
	"errors"
	"io"
	"net/url"
)

// 微博API
type IWeibo interface {
	Authorize() (authurl string)                                      // 获取验证URL
	AccessToken(code string) (token string)                           // 换取AccessToken
	PostStatus(api string, args *url.Values) (rst IStatus, err error) // 微博POST接口
	Update(status string) (rst IStatus, err error)                    // 发表微博
	Repost(status string, oriId int64) (rst IStatus, err error)       // 转发微博
	Destroy(oriId int64) (rst IStatus, err error)                     // 删除微博
	Upload(status string, pic io.Reader) (rst IStatus, err error)     // 发表带图片的微博
	UploadUrl(status string, urlText string) (rst IStatus, err error) // 用图片URL发表带图片的微博
}
type IStatus interface {
	Url() (urlText string) // 获取对应网页URL
}
type IUser interface {
	// Url() (urlText string) // 获取微博主页URL
}

var (
	NoAccessTokenError     = errors.New("Need AccessToken")
	AccessTokenLapsedError = errors.New("AccessToken Lapsed")
	RemoteError            = errors.New("Remote Error")
)
