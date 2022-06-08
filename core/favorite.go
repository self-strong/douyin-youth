package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Not Logged in or Not Exist"})
		return
	}

	actionType := c.Query("action_type")
	videoId := c.Query("video_id")
	if actionType == "" || videoId == "" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Missing parameter!"})
		return
	}

	vId, _ := strconv.ParseInt(videoId, 10, 64)
	var result error
	if actionType == "1" {
		result = DbFavoriteAction(userLoginInfo.Id, vId)
	} else if actionType == "2" {
		result = DbUnFavoriteAction(userLoginInfo.Id, vId)
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Illegal parameter!"})
		return
	}

	if result != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: result.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Action " + actionType + " on " + videoId + " succeeded"})
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User Not Logged in or Not Exist"})
		return
	}

	uIdStr := c.Query("user_id")
	uId, _ := strconv.ParseInt(uIdStr, 10, 64)

	result := DbFavoriteList(uId)
	if result == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Favorite list is empty"})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Succeed",
		},
		VideoList: result,
	})
	// return
}
