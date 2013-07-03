package target

import (
	"bytes"
	"encoding/json"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	TSINA_OAUTH_VERSION = "2.0"
)

func (this *SinaWeibo) Send(src *source.FeedInfo) (rid string, e error) {
	if util.DEBUG {
		c := []rune(src.Content)
		if len(c) > 100 {
			c = c[0:100]
		}
		util.Log(src.SourceId, ":", src.Id, string(c))
	}
	if src.RepostId != "" {
		r, err := this.Repost(src.Content, src.RepostId)
		if err != nil {
			e = err
			return
		}
		util.Log("Repost sent:", r.Url())
		return util.ToString(r.Id), nil
	} else if src.PicUrl != nil && len(src.PicUrl) > 0 {
		if this.EnableUploadUrl {
			r, err := this.UploadUrl(src.Content, src.PicUrl[0])
			if err != nil {
				e = err
				return
			}
			util.Log("UploadUrl sent:", r.Url())
			return util.ToString(r.Id), nil
		} else {
			pic := util.FetchImageAsStream(src.PicUrl[0])
			r, err := this.Upload(src.Content, pic)
			if err != nil {
				e = err
				return
			}
			util.Log("Upload sent:", r.Url())
			return util.ToString(r.Id), nil
		}
	} else {
		r, err := this.Update(src.Content)
		if err != nil {
			e = err
			return
		}
		util.Log("Update sent:", r.Url())
		return util.ToString(r.Id), nil
	}
	return
}

func (this *SinaWeibo) GetMethod() []*TargetMethod { return this.Method }
func (this *SinaWeibo) GetId() string              { return this.Name }

type SinaWeibo struct { // 新浪微博API 实现接口IWeibo
	IWeibo
	ITarget
	Target

	AppKey      string    `json:"client_id"`     // AppKey
	AppSecret   string    `json:"client_secret"` // AppSecret
	CallbackUrl string    `json:"redirect_uri"`  // 验证URL
	Token       string    `json:"access_token"`  // OAuth2.0 验证码
	ExpiresIn   time.Time `json:"expires_in"`    // 失效时间

	EnableUploadUrl bool `json:"enable_upload_url"` // 启用高级接口 直接使用URL上传
}
type SinaWeiboError struct { // 错误
	Request   string `json:"request"`    // 请求
	ErrorCode string `json:"error_code"` // 错误代码
	Error     string `json:"error"`      // 错误信息
}
type SinaWeiboStatus struct { // 微博 实现接口IStatus
	IStatus
	SinaWeiboError

	CreatedAt       string           `json:"created_at"`                 // 创建时间
	Id              int64            `json:"id"`                         // 微博ID
	Mid             string           `json:"mid"`                        // 微博MID
	Idstr           string           `json:"idstr"`                      // 字符串型微博ID
	Text            string           `json:"text"`                       // 微博内容
	Source          string           `json:"source"`                     // 微博来源
	Favorited       bool             `json:"favorited"`                  // 是否已收藏
	Truncated       bool             `json:"truncated"`                  // 是否被截断
	ThumbnailPic    string           `json:"thumbnail_pic"`              // 图片
	BmiddlePic      string           `json:"bmiddle_pic"`                // 图片
	OriginalPic     string           `json:"original_pic"`               // 图片
	User            *SinaWeiboUser   `json:"user,omitempty"`             // PO主
	RetweetedStatus *SinaWeiboStatus `json:"retweeted_status,omitempty"` // 被转发微博信息
	RepostsCount    int              `json:"reposts_count"`              // 转发数
	CommentsCount   int              `json:"comments_count"`             // 评论数
	AttitudesCount  int              `json:"attitudes_count"`            // 表态数
	// Mlevel              int64         `json:"mlevel"`                     // 官方未支持
	// InReplyToStatusId   string        `json:"in_reply_to_status_id"`      // 官方未支持
	// InReplyToUserId     string        `json:"in_reply_to_user_id"`        // 官方未支持
	// InReplyToScreenName string        `json:"in_reply_to_screen_name"`    // 官方未支持
	// Geo                 *WeiboGeo     `json:"geo"`                        // 地理位置信息(无需使用)
}
type SinaWeiboUser struct { // 用户
	IUser
	// SinaWeiboError

	Id               int64            `json:"id"`                 // UID
	Idstr            string           `json:"idstr"`              // 字符串型UID
	Screen_name      string           `json:"screen_name"`        // 昵称
	Name             string           `json:"name"`               // 友好显示名称
	Province         string           `json:"province"`           // 所在省[代码]
	City             string           `json:"city"`               // 所在城市[代码]
	Location         string           `json:"location"`           // 所在地
	Description      string           `json:"description"`        // 描述
	Url              string           `json:"url"`                // 主页地址
	ProfileImageUrl  string           `json:"profile_image_url"`  // 头像
	ProfileUrl       string           `json:"profile_url"`        // 统一URL地址
	Domain           string           `json:"domain"`             // 个性域名
	Weihao           string           `json:"weihao"`             // 微号
	Gender           string           `json:"gender"`             // 性别
	FollowersCount   int              `json:"followers_count"`    // 粉丝数量
	FriendsCount     int              `json:"friends_count"`      // 好友数量
	StatusesCount    int              `json:"statuses_count"`     // 微博数量
	FavouritesCount  int              `json:"favourites_count"`   // 关注数量
	CreatedAt        string           `json:"created_at"`         // 注册时间
	Following        bool             `json:"following"`          // 是否关注
	AllowAllActMsg   bool             `json:"allow_all_act_msg"`  // 允许所有动态
	AllowAllComment  bool             `json:"allow_all_comment"`  // 允许所有评论
	Geo_enabled      bool             `json:"geo_enabled"`        // 区域
	Verified         bool             `json:"verified"`           // 验证
	VerifiedType     int              `json:"verified_type"`      // 验证类型
	VerifiedReason   string           `json:"verified_reason"`    // 验证信息
	Remark           string           `json:"remark"`             // 备注
	AvatarLarge      string           `json:"avatar_large"`       // 头像
	FollowMe         bool             `json:"follow_me"`          // 关注我
	OnlineStatus     int              `json:"online_status"`      // 在线状态
	BiFollowersCount int              `json:"bi_followers_count"` // 互粉数
	Lang             string           `json:"lang"`               // 语言
	Status           *SinaWeiboStatus `json:"status,omitempty"`   // 微博
}
type SinaWeiboComment struct { // 评论
	CreatedAt    string            `json:"created_at"`
	Id           int64             `json:"id"`
	Text         string            `json:"text"`
	Source       string            `json:"source"`
	User         *SinaWeiboUser    `json:"user,omitempty"`
	Mid          string            `json:"mid"`
	Idstr        string            `json:"idstr"`
	Status       *SinaWeiboStatus  `json:"status,omitempty"`
	ReplyComment *SinaWeiboComment `json:"reply_comment,omitempty"`
}
type SinaWeiboReposts struct { // 转发
	Reposts        []*SinaWeiboStatus `json:"reposts"`
	Hasvisible     bool               `json:"hasvisible"`
	PreviousCursor int                `json:"previous_cursor"`
	NextCursor     int                `json:"next_cursor"`
	TotalNumber    int                `json:"total_number"`
}
type SinaWeiboPrivacy struct { // 隐私设置
	Comment  int `json:"comment"`  // 是否可以评论我的微博，0：所有人、1：关注的人、2：可信用户
	Geo      int `json:"geo"`      // 是否开启地理信息，0：不开启、1：开启
	Message  int `json:"message"`  // 是否可以给我发私信，0：所有人、1：我关注的人、2：可信用户
	Realname int `json:"realname"` // 是否可以通过真名搜索到我，0：不可以、1：可以
	Badge    int `json:"badge"`    // 勋章是否可见，0：不可见、1：可见
	Mobile   int `json:"mobile"`   // 是否可以通过手机号码搜索到我，0：不可以、1：可以
	Webim    int `json:"webim"`    // 是否开启webim， 0：不开启、1：开启
}
type SinaWeiboRemind struct { // 消息未读数
	Status        int `json:"status"`         // 新微博未读数
	Follower      int `json:"follower"`       // 新粉丝数
	Cmt           int `json:"cmt"`            // 新评论数
	Dm            int `json:"dm"`             // 新私信数
	MentionStatus int `json:"mention_status"` // 新提及我的微博数
	MentionCmt    int `json:"mention_cmt"`    // 新提及我的评论数
	Group         int `json:"group"`          // 微群消息未读数
	PrivateGroup  int `json:"private_group"`  // 私有微群消息未读数
	Notice        int `json:"notice"`         // 新通知未读数
	Invite        int `json:"invite"`         // 新邀请未读数
	Badge         int `json:"badge"`          // 新勋章数
	Photo         int `json:"photo"`          // 相册消息未读数
}
type SinaWeiboUrlShort struct { // 短链
	UrlShort string `json:"url_short"` // 短链接
	UrlLong  string `json:"url_long"`  // 原始长链接
	Type_    int    `json:"type"`      // 链接的类型，0：普通网页、1：视频、2：音乐、3：活动、5、投票
	Result   bool   `json:"result"`    // 短链的可用状态，true：可用、false：不可用。
}
type SinaWeiboGeo struct { // 地理信息
	Longitude    string `json:"longitude"`     // 经度坐标
	Latitude     string `json:"latitude"`      // 维度坐标
	City         string `json:"city"`          // 所在城市的城市代码
	Province     string `json:"province"`      // 所在省份的省份代码
	CityName     string `json:"city_name"`     // 所在城市的城市名称
	ProvinceName string `json:"province_name"` // 所在省份的省份名称
	Address      string `json:"address"`       // 所在的实际地址，可以为空
	Pinyin       string `json:"pinyin"`        // 地址的汉语拼音，不是所有情况都会返回该字段
	More         string `json:"more"`          // 更多信息，不是所有情况都会返回该字段
}

func NewSinaWeibo(client_id, client_secret, access_token, redirect_uri string) (this *SinaWeibo) {
	this = &SinaWeibo{
		AppKey:      client_id,
		AppSecret:   client_secret,
		Token:       access_token,
		CallbackUrl: redirect_uri,
	}
	return
}
func (this *SinaWeibo) Authorize() (authurl string) {
	return "https://api.weibo.com/oauth2/authorize?" + (url.Values{
		"client_id":     {this.AppKey},
		"redirect_uri":  {this.CallbackUrl},
		"response_type": {"code"},
		"display":       {"client"},
	}).Encode()
}
func (this *SinaWeibo) AccessToken(code string) (token string) {
	res, err := http.PostForm("https://api.weibo.com/oauth2/access_token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {this.AppKey},      // yourappkey
			"client_secret": {this.AppSecret},   // yourpppsecret
			"code":          {code},             // xxxxxxxxxxxxxx
			"redirect_uri":  {this.CallbackUrl}, // http://some/weibocb.php
		})
	if err != nil {
		util.Log("Fail to AccessToken:", err)
		return
	}

	defer res.Body.Close()
	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)
	if body["error"] != nil || body["access_token"] == nil {
		util.Log("Fail to AccessToken(Remote):", body["error"])
		return
	}
	this.Token = body["access_token"].(string)
	i, _ := strconv.Atoi(body["expires_in"].(string))
	ex := time.Now().Add(time.Duration(i) * time.Second)
	this.ExpiresIn = ex
	return this.Token
}
func (this *SinaWeibo) PostStatus(api string, args *url.Values) (rst *SinaWeiboStatus, err error) {
	args.Set("access_token", this.Token)
	res, err := http.PostForm("https://api.weibo.com/2/statuses/"+api+".json", *args)
	if err != nil {
		util.Log("Error on call", api+":", err)
		return
	}
	defer res.Body.Close()
	rst = &SinaWeiboStatus{}
	json.NewDecoder(res.Body).Decode(rst)
	if rst.Error != "" {
		util.Log("Error on call", api+"(Remote):", args.Encode(), ":", rst.Error, "\nOn:", rst.Request)
		return nil, RemoteError(rst.Error)
	}
	return
}
func (this *SinaWeibo) Update(status string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.PostStatus("update", &url.Values{
		"status": {status},
	})
	return
}
func (this *SinaWeibo) Repost(status string, oriId string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.PostStatus("repost", &url.Values{
		"status": {status},
		"id":     {oriId},
	})
	return
}
func (this *SinaWeibo) Destroy(oriId string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.PostStatus("destroy", &url.Values{
		"id": {oriId},
	})
	return
}
func (this *SinaWeibo) Upload(status string, pic io.Reader) (rst *SinaWeiboStatus, err error) {
	// multipart/form-data
	var bpic bytes.Buffer
	formdata := multipart.NewWriter(&bpic)
	formdata.WriteField("access_token", this.Token)
	formdata.WriteField("status", status)
	picdata, _ := formdata.CreateFormFile("pic", "image.png")
	io.Copy(picdata, pic)
	formdata.Close()

	res, err := http.Post("https://api.weibo.com/2/statuses/upload.json", formdata.FormDataContentType(), &bpic)
	if err != nil {
		util.Log("Error on call upload :", err)
		return
	}
	defer res.Body.Close()
	rst = &SinaWeiboStatus{}
	json.NewDecoder(res.Body).Decode(rst)
	if rst.Error != "" {
		util.Log("Error on call upload (Remote):", rst.ErrorCode, ":", rst.Error, "\nOn:", rst.Request)
		return nil, RemoteError(rst.Error)
	}
	return
}
func (this *SinaWeibo) UploadUrl(status string, urlText string) (rst *SinaWeiboStatus, err error) {
	rst, err = this.PostStatus("upload_url_text", &url.Values{
		"status": {status},
		"url":    {urlText},
	})
	return
}
func (this *SinaWeiboStatus) Url() (urlText string) {
	if this == nil || this.User == nil {
		b, _ := json.Marshal(this)
		util.Log("Bad response!", string(b))
	}
	urlText = "http://weibo.com/" + this.User.Idstr + "/" + util.Base62Url(this.Mid)
	return
}
