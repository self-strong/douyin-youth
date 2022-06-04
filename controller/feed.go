package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/self-strong/douyin-youth/repository"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  string  `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
// 不限制登陆状态，返回按投稿时间倒序的视频列表，由服务端控制；最多30个
func Feed(c *gin.Context) {
	video_list := []Video{}
	var userIdlist []int64
	videos, err := repository.Feed() // 需要创建一个videolist
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1},
			VideoList: nil,
			NextTime:  time.Now().String(),
		})
	}

	// 根据video创建list
	for i := range videos {
		tempid := videos[i].CreateUid
		userIdlist = append(userIdlist, tempid)
	}
	// repository.SearchUserById(userIdlist)
	user_, err := repository.SearchUserById(userIdlist) // 通过序列查询

	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1},
			VideoList: nil,
		})
	}
	for i := range videos {
		user := User{
			Id:            user_[i].Id,
			Name:          user_[i].Name,
			FollowCount:   user_[i].FollowCount,
			FollowerCount: user_[i].FanCount,
			IsFollow:      false, // 查表看是否关注!!!!
		}

		video := videos[i]
		video_ := Video{
			Id:            video.Id,
			Author:        user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.ThumbCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true, // 查表确认是都点赞!!!!
			Title:         video.Title,
		}

		video_list = append(video_list, video_)
	}
	// for i := range videos {
	// 	video := videos[i]
	// 	userId := video.CreateUid

	// 	user_, err := repository.SearchUserById(userId) // 应该直接使用序列查询
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, FeedResponse{
	// 			Response:  Response{StatusCode: 1},
	// 			VideoList: video_list,
	// 			NextTime:  time.Now().Unix(),
	// 		})
	// 	}
	// 	user := User{
	// 		Id:            user_.Id,
	// 		Name:          user_.Name,
	// 		FollowCount:   user_.FollowCount,
	// 		FollowerCount: user_.FanCount,
	// 		IsFollow:      true, // 查表看是否关注!!!!

	// 	}

	// 	video_ := Video{
	// 		Id:            video.Id,
	// 		Author:        user,
	// 		PlayUrl:       video.PlayUrl,
	// 		CoverUrl:      video.CoverUrl,
	// 		FavoriteCount: video.ThumbCount,
	// 		CommentCount:  video.CommentCount,
	// 		IsFavorite:    true, // 查表确认是都点赞!!!!
	// 		Title:         video.Title,
	// 	}

	// 	video_list = append(video_list, video_)
	// }

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "Successful!"},
		VideoList: video_list,
		NextTime:  time.Now().String(),
	})
}
