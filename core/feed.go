package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  string  `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
// 不限制登陆状态，返回按投稿时间倒序的视频列表，由服务端控制；最多30个
func Feed(c *gin.Context) {

	token := c.Query("token")
	// video_list := []Video{}

	userLoginInfo := DbFindUserLoginInfo(token) // 根据token获取用户信息

	if userLoginInfo == nil {
		video_list := DbFeed()
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0, StatusMsg: "Successful!"},
			VideoList: video_list,
			NextTime:  time.Now().String(),
		})
	} else {
		video_list := DbFeedWithLogin(userLoginInfo.Id)
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0, StatusMsg: "Successful!"},
			VideoList: video_list,
			NextTime:  time.Now().String(),
		})
	}

}
