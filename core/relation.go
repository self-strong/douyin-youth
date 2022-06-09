package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User not Logged in or Not Exist"},
			Comment:  Comment{},
		})
		return
	}

	actionType := c.Query("action_type")
	toId := c.Query("to_user_id")
	if actionType == "" || toId == "" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Missing Parameter"})
		return
	}

	toIdInt, _ := strconv.ParseInt(toId, 10, 64)

	if userLoginInfo.Id == toIdInt {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "can not follow yourself"})
		return
	}
	var result error
	if actionType == "1" {
		result = DbFollowAction(userLoginInfo.Id, toIdInt)
	} else if actionType == "2" {
		result = DbUnFollowAction(userLoginInfo.Id, toIdInt)
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Illegal Parameter"})
		return
	}

	if result != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: result.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "Action succeeded",
	})
}

func FollowList(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User not Logged in or Not Exist"},
			UserList: nil,
		})
		return
	}

	userIdStr := c.Query("user_id")
	if userIdStr == "" {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
			UserList: nil,
		})
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	followList := DbFollowList(userId, userLoginInfo.Id)
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Query succeeded"},
		UserList: followList,
	})
}

func FollowerList(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User not Logged in or Not Exist"},
			UserList: nil,
		})
		return
	}

	userIdStr := c.Query("user_id")
	if userIdStr == "" {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
			UserList: nil,
		})
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	followerList := DbFollowerList(userId, userLoginInfo.Id)
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Query succeeded"},
		UserList: followerList,
	})
}
