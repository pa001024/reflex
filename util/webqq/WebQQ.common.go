package webqq

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pa001024/MoeWorker/util"
)

// 生成随机数
func rand_r() string {
	return fmt.Sprint(rand.ExpFloat64())
}

// 进行三次加盐(uin的16b LE hex值)MD5之类...什么的算法 [2013.8.27]
func (this *WebQQ) genPwd(code string) string {
	salt := util.EncodeJsUint64LE(uint64(this.Uin))
	vSaltedPwd := util.Md5StringX(this.PasswdMd5 + salt)
	return util.Md5StringX(vSaltedPwd + strings.ToUpper(code))
}

/*
 [算法] 获取好友列表的hash算法
 -----------------------------

 "h":"hello"
 1. 这是一个32位分组密码
 2. 取十进制的uin逐位与
 2. ptwebqq中的每4个char值分组 OR 入c , 溢出则从 0 开始计算
 3. 再 XOR 入d
 4. 得出最后32位hash d 以hex形式返回

 function(b, i) {
     for (var a = [], s = 0; s < b.length; s++)
         a[s] = b.charAt(s) - 0;
     for (var j = 0, d = -1, s = 0; s < a.length; s++) {
         j += a[s];
         j %= i.length;
         var c = 0;
         if (j + 4 > i.length)
             for (var l = 4 + j - i.length, x = 0; x < 4; x++)
                 c |= x < l ? (i.charCodeAt(j + x) & 255) << (3 - x) * 8 : (i.charCodeAt(x - l) & 255) << (3 - x) * 8;
         else
             for (x = 0; x < 4; x++)
                 c |= (i.charCodeAt(j + x) & 255) << (3 - x) * 8;
         d ^= c
     }
     a = [];
     a[0] = d >> 24 & 255;
     a[1] = d >> 16 & 255;
     a[2] = d >> 8 & 255;
     a[3] = d & 255;
     d = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"];
     s = "";
     for (j = 0; j < a.length; j++)
         s += d[a[j] >> 4 & 15], s += d[a[j] & 15];
     return s
 }
*/
func (this *WebQQ) genGetUserFriendsHash() string {
	uin := []byte(this.Id.String())
	for i, v := range uin {
		uin[i] = v - 48
	}
	j := uint32(0)
	pt := []byte(this.PtWebQQ)
	ptlen := uint32(len(pt))
	d := uint32(0xffffffff)
	for i, _ := range uin {
		j = (j + uint32(uin[i])) % ptlen
		c := uint32(0)
		if j+4 > ptlen {
			l := 4 + j - ptlen
			for x := uint32(0); x < 4; x++ {
				if x < l {
					c |= uint32(pt[j+x]) << ((3 - x) * 8)
				} else {
					c |= uint32(pt[x-l]) << ((3 - x) * 8)
				}
			}
		} else {
			for x := uint32(0); x < 4; x++ {
				c |= uint32(pt[j+x]) << ((3 - x) * 8)
			}
		}
		d ^= c
	}
	return fmt.Sprintf("%X", d)
}
