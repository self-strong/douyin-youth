package core

// core 中共用的类型部分


// response
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// User
type User struct {
	Uid       int64  `json:"id"`
	Username  string `json:"name"`
	Follow    int64  `json:"follow_count"`
	Following int64  `json:"follower_count"`
	Is_follow bool   `json:"is_follow"`
}

// video
type Video struct {
	Id           int64  `json:"id"`
	User         User   `json:"author"`
	PlayUrl      string `json:"play_url"`
	CoverUrl     string `json:"cover_url"`
	ThumbCount   int64  `json:"favorite_count"`
	CommentCount int64  `json:"comment_count"`
	Is_favorite  bool   `json:"is_favorite"`
	Title        string `json:"title"`
}