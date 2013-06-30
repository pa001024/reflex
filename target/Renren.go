package target

import ()

const (
	RENREN_OAUTH_VERSION = "2.0"
)

type Renren struct { // 人人网 API 实现接口IWeibo
	// IWeibo 暂未实现
	// ITarget
	Target

	AppKey      string `json:client_id`     // Consumer key
	AppSecret   string `json:client_secret` // Consumer secret
	CallbackUrl string `json:redirect_uri`  // 验证URL
}
