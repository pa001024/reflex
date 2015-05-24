package tieba

type Result struct {
	ErrorCode int32  `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Time      int32  `json:"time"`
	Ctime     int32  `json:"ctime"`
	Logid     int32  `json:"logid"`
}
type TbsResult struct {
	Tbs     string `json:"tbs"`
	IsLogin int32  `json:"is_login"`
}

type Forum struct {
	ForumId  int32  `json:"forum_id"`
	IsLike   int32  `json:"is_like"`
	FavoType int32  `json:"favo_type"`
	Name     string `json:"name"`
	LevelId  string `json:"level_id"`
}
type LikedResult struct {
	Result
	Anti struct {
		Tbs string `json:"tbs"`
	} `json:"anti"`
	ForumList []Forum `json:"forum_list"`
	Page      struct {
		PageSize    string `json:"page_size"`
		Offset      int32  `json:"offset"`
		CurrentPage int32  `json:"current_page"`
		TotalCount  int32  `json:"total_count"`
		TotalPage   int32  `json:"total_page"`
		HasMore     int32  `json:"has_more"`
		HasPrev     int32  `json:"has_prev"`
	} `json:"page"`
}

type LoginResult struct {
	Result
	Anti struct {
		NeedVcode   int32  `json:"need_vcode"`
		VcodeMd5    string `json:"vcode_md5"`
		VcodePicUrl string `json:"vcode_pic_url"`
		Tbs         string `json:"tbs"`
	} `json:"anti"`
	User struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		BDUSS    string `json:"BDUSS"`
		Passwd   string `json:"passwd"`
		Portrait string `json:"portrait"`
	} `json:"user"`
}
