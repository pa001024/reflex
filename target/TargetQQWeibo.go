package target

import (
	"github.com/pa001024/reflex/source"
	"github.com/pa001024/reflex/util"
	"github.com/pa001024/reflex/util/tqq"
)

func (this *TargetQQWeibo) Send(src *source.FeedInfo) (rid string, e error) {
	util.DEBUG.Logf("QQWeibo.Send(%v:%v,repost_id:%v,title:%v,content:%v,picurl:%v)", src.SourceId, src.Id, src.RepostId, src.Title, src.Content, src.PicUrl)
	if src.RepostId != "" {
		r, err := this.Repost(src.Content, src.RepostId)
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] Repost sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	} else if src.PicUrl != nil && len(src.PicUrl) > 0 {
		r, err := this.UploadUrl(src.Content, src.PicUrl[0])
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] UploadUrl sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	} else {
		r, err := this.Update(src.Content)
		if err != nil {
			e = err
			return
		}
		util.INFO.Logf("[qq.%v] Update sent: %v", this.Name, r.Url())
		return util.ToString(r.Id), nil
	}
	return
}

func (this *TargetQQWeibo) GetMethod() []*TargetMethod { return this.Method }
func (this *TargetQQWeibo) GetId() string              { return this.Name }

const (
	TQQ_OAUTH_VERSION = "2.a"
)

// 腾讯微博API
type TargetQQWeibo struct {
	ITarget
	Target
	tqq.QQWeibo
}
