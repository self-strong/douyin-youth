package core

import "C"
import (
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// HttpContentType 文件类型
var HttpContentType = map[string]string{
	".avi":  "video/avi",
	".mp3":  "   audio/mp3",
	".mp4":  "video/mp4",
	".wmv":  "   video/x-ms-wmv",
	".asf":  "video/x-ms-asf",
	".rm":   "application/vnd.rn-realmedia",
	".rmvb": "application/vnd.rn-realmedia-vbr",
	".mov":  "video/quicktime",
	".m4v":  "video/mp4",
	".flv":  "video/x-flv",
	".jpg":  "image/jpeg",
	".png":  "image/png",
}

// Feed same demo video list for every request
// 不限制登陆状态，返回按投稿时间倒序的视频列表，由服务端控制；最多30个
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	token := c.Query("token") // 获得token
	latestTime, _ := strconv.ParseInt(latestTimeStr, 10, 64)
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	//userLoginInfo := DbFindUserInfoByToken(token) // 根据token获取用户信息

	videoList, nextTime := DbFeed(latestTime, token)
	if videoList == nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1,
				StatusMsg: "this is a bug",
			},
		})
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0, StatusMsg: "Successfully!"},
			VideoList: videoList,
			NextTime:  nextTime,
		})
	}

	//if userLoginInfo == nil {
	//
	//} else {
	//	videoList := DbFeedWithLogin(userLoginInfo.Id)
	//	c.JSON(http.StatusOK, FeedResponse{
	//		Response:  Response{StatusCode: 0, StatusMsg: "Successful!"},
	//		VideoList: videoList,
	//		NextTime:  time.Now().Unix(),
	//	})
	//}
}

// GetVideo 获取视频
func GetVideo(c *gin.Context) {
	videoName := c.Query("videoName")
	videoPath := "./public/video/" + videoName
	//获取文件的后缀
	fileType := path.Ext(videoName)
	//获取文件类型对应的http ContentType 类型
	fileContentType := HttpContentType[fileType]

	c.Header("Content-Type", fileContentType)
	c.File(videoPath)
}

// GetCover 获取封面
func GetCover(c *gin.Context) {
	coverName := c.Query("coverName")
	coverPath := "./public/cover/" + coverName
	//获取文件的后缀
	fileType := path.Ext(coverName)
	//获取文件类型对应的http ContentType 类型
	fileContentType := HttpContentType[fileType]
	//fmt.Println(coverPath, coverName, fileContentType)

	c.Header("Content-Type", fileContentType)
	c.File(coverPath)

}
