package target

import (
	"net/url"
)

// 微博API {{
type IWeibo interface {
	// TODO: 添加error
	Authorize() (authurl string)                           // 获取验证URL
	AccessToken(code string) (token string)                // 换取AccessToken
	PostStatus(api string, args *url.Values) (rst IStatus) // 微博POST接口
	Update(status string) (rst IStatus)                    // 发表微博
	Repost(status string, oriId int64) (rst IStatus)       // 转发微博
	Destroy(oriId int64) (rst IStatus)                     // 删除微博
	Upload(status string, pic io.Reader) (rst IStatus)     // 发表带图片的微博
	UploadUrl(status string, urlText string) (rst IStatus) // 用图片URL发表带图片的微博
}
type IStatus interface {
	Url() (urlText string) // 获取对应网页URL
}
type IUser interface{}
