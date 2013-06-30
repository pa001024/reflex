package target

import (
// "net/http"
// "net/url"
)

const (
	FACEBOOK_OAUTH_VERSION = "2.0"
)

type Facebook struct { // Facebook API 实现接口IWeibo
	// IWeibo 暂未实现
	// ITarget
	Target

	AppKey      string `json:client_id`     // Consumer key
	AppSecret   string `json:client_secret` // Consumer secret
	CallbackUrl string `json:redirect_uri`  // 验证URL
}
