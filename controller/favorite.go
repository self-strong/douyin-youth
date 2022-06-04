package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/self-strong/douyin-youth/repository"
)

// FavoriteAction no practical effect, just check if token is valid
// 用户id, token, video_id, action_type
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	videoId := c.Query("video_id")
	action := c.Query("action_type")
	userid, _ := strconv.ParseInt(userId, 10, 64)
	videoid, _ := strconv.ParseInt(videoId, 10, 64)

	if _, exist := usersLoginInfo[token]; exist {
		if action == "1" {
			//点赞
			if _, err := repository.CreateThumb(userid, videoid); err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "thumb failed!"})
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Cancel Successful"})

		} else {
			// 取消点赞， 删除
			if _, err := repository.CancelThumb(userid, videoid); err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "cancel failed"})
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Cancel Successful"})
		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	userid, _ := strconv.ParseInt(userId, 10, 64)
	var video_list []Video
	var userIdlist []int64
	// var user []User

	if _, exist := usersLoginInfo[token]; exist {

		videoId, err := repository.SearchThumbVideo(userid) // 根据用户id查询已点赞的视频ID
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response:  Response{StatusCode: 1},
				VideoList: nil,
			})
		}
		// 根据视频id返回视频,
		videos, err := repository.SearchVideoById(videoId)

		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response:  Response{StatusCode: 0, StatusMsg: "Successful!"},
				VideoList: video_list,
			})
		}
		// 根据用户id查询用户信息
		for i := range videos {
			tempid := videos[i].CreateUid
			userIdlist = append(userIdlist, tempid) // 发布视频的id列表
		}

		user_, err := repository.SearchUserById(userIdlist) // 通过序列查询

		// 通过用户id和发布视频的id， 查找following表判断是否关注
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response:  Response{StatusCode: 1},
				VideoList: video_list,
			})
		}
		for i := range videos {
			user := User{
				Id:            user_[i].Id,
				Name:          user_[i].Name,
				FollowCount:   user_[i].FollowCount,
				FollowerCount: user_[i].FanCount,
				IsFollow:      true, // 查表看是否关注!!!!
			}

			video := videos[i]
			video_ := Video{
				Id:            video.Id,
				Author:        user,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: video.ThumbCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    true, // 点赞列表，必然所有都点赞!!!!
				Title:         video.Title,
			}

			video_list = append(video_list, video_)
		}

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: video_list,
		})

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}
