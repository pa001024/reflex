package target

import (
	"github.com/pa001024/reflex/source"
	"github.com/pa001024/reflex/util"
	"github.com/pa001024/reflex/util/weibo"
)

func (this *TargetSinaWeibo) Send(src *source.FeedInfo) (rid string, e error) {
	util.DEBUG.Logf("SinaWeibo.Send(%v:%v,repost_id:%v,title:%v,content:%v,picurl:%v)", src.SourceId, src.Id, src.RepostId, src.Title, src.Content, src.PicUrl)
	if src.RepostId != "" {
		r, err := this.Repost(src.Content, src.RepostId)
		if err != nil {
			e = err
			return
		}
		util.INFO.Log("[sina."+this.Name+"] Repost sent:", r.Url())
		return util.ToString(r.Id), nil
	} else if src.PicUrl != nil && len(src.PicUrl) > 0 {
		if this.EnableUploadUrl {
			r, err := this.UploadUrl(src.Content, src.PicUrl[0])
			if err != nil {
				e = err
				return
			}
			util.INFO.Log("[sina.%s] UploadUrl sent: %v", this.Name, r.Url())
			return util.ToString(r.Id), nil
		} else {
			pic := util.FetchImageAsStream(src.PicUrl[0])
			r, err := this.Upload(src.Content, pic)
			if err != nil {
				e = err
				return
			}
			util.INFO.Log("[sina.%s] Upload sent: %v", this.Name, r.Url())
			return util.ToString(r.Id), nil
		}
	} else {
		r, err := this.Update(src.Content)
		if err != nil {
			e = err
			return
		}
		util.INFO.Log("Update sent:", r.Url())
		return util.ToString(r.Id), nil
	}
	return
}

func (this *TargetSinaWeibo) GetMethod() []*TargetMethod { return this.Method }
func (this *TargetSinaWeibo) GetId() string              { return this.Name }

// 新浪微博API
type TargetSinaWeibo struct {
	ITarget
	Target
	weibo.SinaWeibo
}
