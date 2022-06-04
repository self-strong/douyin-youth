package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/self-strong/douyin-youth/repository"
)

type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	// 数据库的插入和删除
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {

		userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		videoid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
		action_type := c.Query("action_type")
		// 评论
		if action_type == "1" {
			content := c.Query("commnet_text")
			comment, err := repository.CreateComment(content, userid, videoid) // 从数据库返回的comment

			if err != nil {
				c.JSON(http.StatusOK, CommentResponse{
					Response: Response{StatusCode: 1},
					Comment:  Comment{},
				})
			} else {
				user, err := repository.SearchOneUserById(comment.Uid) // 查找用户信息
				if err != nil {
					c.JSON(http.StatusOK, CommentResponse{
						Response: Response{StatusCode: 1, StatusMsg: "Failed!"},
						Comment:  Comment{},
					})
				}
				user_ := User{
					Id:            user.Id,
					Name:          user.Name,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FanCount,
					IsFollow:      false, // 需要查表
				}

				comment_ := Comment{
					Id:         comment.CmId,
					User:       user_,
					Content:    comment.Content,
					CreateDate: comment.Timestamp,
				}

				c.JSON(http.StatusOK, CommentResponse{
					Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
					Comment:  comment_,
				})
			}
		} else {
			// 删除评论
			commentid, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
			if err := repository.DeleteComment(commentid); err != nil {

				c.JSON(http.StatusOK, CommentResponse{
					Response: Response{StatusCode: 1, StatusMsg: err.Error()},
					Comment:  Comment{},
				})

			} else {
				c.JSON(http.StatusOK, CommentResponse{
					Response: Response{StatusCode: 0, StatusMsg: "Successful Delete!"},
					Comment:  Comment{},
				})
			}
		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// token := c.Query("token")
	videoId := c.Query("video_id")

	var comment_list []Comment

	// 根据videoid查询得到评论结构体，根据评论获取用户id

	videoid, _ := strconv.ParseInt(videoId, 10, 64) // string转换为int64
	comments, err := repository.SearchComment(videoid)

	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1},
			CommentList: comment_list,
		})
	}

	// 根据评论的uid来获取用户
	var userIdlist []int64
	for i := range comments {
		tempid := comments[i].Uid
		userIdlist = append(userIdlist, tempid)
	}

	// 根据用户id查询
	users, err := repository.SearchUserById(userIdlist) // 通过序列查询
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1},
			CommentList: comment_list,
		})
	}
	// 根据comment和user返回数据
	for i := range comments {

		user := users[i]
		user_ := User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FanCount,
			IsFollow:      true, // 查表看是否关注!!!!
		}

		comment := comments[i]

		comment_ := Comment{
			Id:         comment.CmId,
			User:       user_,
			Content:    comment.Content,
			CreateDate: comment.Timestamp,
		}

		comment_list = append(comment_list, comment_)
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "Successful!"},
		CommentList: comment_list,
	})
}
