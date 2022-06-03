package core

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserLoginInfo(token)

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
	var result *gorm.DB
	if actionType == "1" {
		result = DbFollowAction(userLoginInfo.Id, toIdInt)
	} else if actionType == "2" {
		result = DbUnFollowAction(userLoginInfo.Id, toIdInt)
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Illegal Parameter"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "Action succeeded",
	})
}

func FollowList(c *gin.Context) {

}

func FollowerList(c *gin.Context) {

}
