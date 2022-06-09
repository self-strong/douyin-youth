package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	Response
	Comment Comment `json:"comment"`
}
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
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
	videoId := c.Query("video_id")
	if actionType == "" || videoId == "" {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
			Comment:  Comment{},
		})
		return
	}

	var result error
	vID, _ := strconv.ParseInt(videoId, 10, 64)
	if actionType == "1" {
		content := c.Query("comment_text")
		if content == "" {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
				Comment:  Comment{},
			})
			return
		}
		var comment Comment
		result, comment = DbPostComment(userLoginInfo.Id, vID, content)
		if result == nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Comment successfully"},
				Comment:  comment,
			})
			return
		}
	} else if actionType == "2" {
		cmId := c.Query("comment_id")
		if cmId == "" {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
				Comment:  Comment{},
			})
			return
		}
		cmIdInt, _ := strconv.ParseInt(cmId, 10, 64)
		result = DbDeleteComment(cmIdInt, vID)
		if result == nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Remove comment successfully"},
				Comment:  Comment{},
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Illegal Parameter"},
			Comment:  Comment{},
		})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: result.Error()})
	return
}

func CommentList(c *gin.Context) {
	token := c.Query("token")

	userLoginInfo := DbFindUserInfoByToken(token)

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "User not Logged in or Not Exist"},
			CommentList: nil,
		})
		return
	}

	videoId := c.Query("video_id")
	if videoId == "" {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "Missing Parameter"},
			CommentList: nil,
		})
		return
	}

	vID, _ := strconv.ParseInt(videoId, 10, 64)
	comments := DbCommentList(userLoginInfo.Id, vID)

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "Query finished"},
		CommentList: comments,
	})
}
