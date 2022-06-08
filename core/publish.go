package core

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// PublishAction Publish check token then save upload file to public directory
func PublishAction(c *gin.Context) {
	token := c.PostForm("token") // 获取token
	title := c.PostForm("title") // 获取title

	// 根据token获取用户信息
	userLoginInfo := DbFindUserInfoByToken(token)
	// 返回用户不存在
	if userLoginInfo == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't log in"})
		return
	}

	// 判断用户上传的文件是否是视频文件
	contentType := c.GetHeader("Content-Type")
	fmt.Println("Content-Type :", contentType)

	// 获取用户上传的数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 获取用户上传的视频文件名，并设置存放的路径
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(data.Filename))
	videoPath := filepath.Join("./public/video/", fileName)

	// 保存视频文件
	if err := c.SaveUploadedFile(data, videoPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 获取缩略图
	coverName, err := GetVideoCover(videoPath)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		// 上传文件失败后，删去已经保存的文件

		if err := os.Remove(videoPath); err != nil {
			fmt.Println("删除文件", videoPath, "失败")
		}
		return
	}

	// 存入数据库
	if err := DbInsertVideoInfo(userLoginInfo.Id, title, fileName, coverName); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		if err := os.Remove(videoPath); err != nil {
			fmt.Println("删除文件", videoPath, "失败")
		}
		return
	}

	// 返回上传成功消息
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  filepath.Base(data.Filename) + " uploaded successfully",
	})
}

// PublishList 获取发布视频列表
func PublishList(c *gin.Context) {
	uIdStr := c.Query("user_id") // 获取用户ID
	token := c.Query("token")    // 用户登录token

	// 检查token
	if DbFindUserInfoByToken(token) == nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't log in",
			},
			VideoList: nil,
		})
		return
	}

	uId, _ := strconv.ParseInt(uIdStr, 10, 64)
	// 根据用户ID获取用户信息
	user := DbFindUserInfoById(uId)
	// 根据用户ID获取投稿视频
	videoList := DbFindVideoList(user)

	if videoList == nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "no videos!",
			},
			VideoList: nil,
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
	return
}

//http://localhost:8080/douyin/publish/video/?videoName=1_bear.mp4
//http://localhost:8080/douyin/publish/cover/?coverName=1_bear.jpeg
