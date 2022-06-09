package core

// core 中共用的类型部分

// Response response
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type User struct {
	Uid       int64  `json:"id,omitempty"`
	Username  string `json:"name,omitempty"`
	Follow    int64  `json:"follow_count,omitempty"`
	Following int64  `json:"follower_count,omitempty"`
	IsFollow  bool   `json:"is_follow,omitempty"`
}

type Video struct {
	Id           int64  `json:"id,omitempty"`
	Author       User   `json:"author"`
	PlayUrl      string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl     string `json:"cover_url,omitempty"`
	ThumbCount   int64  `json:"favorite_count,omitempty"`
	CommentCount int64  `json:"comment_count,omitempty"`
	IsFavorite   bool   `json:"is_favorite,omitempty"`
	Title        string `json:"title,omitempty"`
}

type Comment struct {
	CmId       int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
