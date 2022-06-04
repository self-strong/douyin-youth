package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/self-strong/douyin-youth/repository"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
// 插入数据，需要更新users数据
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {

		to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64) // 被关注人的id
		action_type := c.Query("action_type")

		if action_type == "1" {
			// 关注
			if err := repository.CreateRelation(user.Id, to_user_id); err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Create failed!"})
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Create Successful!"})
		} else {
			// 取消关注
			if err := repository.DeleteRelation(user.Id, to_user_id); err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Delete failed!"})
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Delete Successful!"})

		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	// var user_list []User

	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		fansid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64) // 此id即粉丝id，返回偶像列表

		idol_list, err := repository.SearchIdols(fansid)
		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				UserList: nil,
			})
		}

		// 根据关注列表的id查询已关注的用户
		var userIdList []int64
		for i := range idol_list {
			userIdList = append(userIdList, idol_list[i].IdolId)
		}
		fmt.Println(userIdList)
		users, _ := repository.SearchUserById(userIdList)
		fmt.Println(users)
		user_list := make([]User, len(users))

		for i := 0; i < len(users); i++ {
			user := users[i]
			fmt.Println(user)
			user_list[i] = User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FanCount,
				IsFollow:      true, // 肯定都是关注的
			}
		}
		// for i := range users {
		// 	user := users[i]

		// 	user_ := User{
		// 		Id:            user.Id,
		// 		Name:          user.Name,
		// 		FollowCount:   user.FollowCount,
		// 		FollowerCount: user.FanCount,
		// 		IsFollow:      true, // 查表
		// 	}
		// 	user_list = append(user_list, user_)
		// }

		// 需要转换为array返回
		// l := len(user_list)

		// var return_list [length]User
		// return_list := [...]User{&user_list[0:]}
		// var return_list = (*[10]User(user_list))
		// var p = (*[3]int)(b) // ok，*p = [11, 12, 13]

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "OK",
			},
			UserList: user_list,
		})

	} else {
		// 鉴权失败
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
			UserList: nil,
		})
	}

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	var user_list []User

	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		idolid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64) // 此id即偶像id， 通过idol查看谁关注了我

		fans_list, err := repository.SearchFans(idolid) // 获取关注列表
		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				UserList: user_list,
			})
		}

		// 根据关注列表的id查询已关注的用户
		var userIdList []int64
		for i := range fans_list {
			userIdList = append(userIdList, fans_list[i].FansId)
		}
		// 根据用户id查询用户
		users, _ := repository.SearchUserById(userIdList)

		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				UserList: nil,
			})
		}

		user_list := make([]User, len(users))

		for i := 0; i < len(users); i++ {
			user := users[i]
			user_list[i] = User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FanCount,
				IsFollow:      true, // 查表
			}
		}
		// for i := range users {
		// 	user := users[i]

		// 	user_ := User{
		// 		Id:            user.Id,
		// 		Name:          user.Name,
		// 		FollowCount:   user.FollowCount,
		// 		FollowerCount: user.FanCount,
		// 		IsFollow:      true, // 查表
		// 	}
		// 	user_list = append(user_list, user_)
		// }

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "OK",
			},
			UserList: user_list,
		})

	} else {
		// 鉴权失败
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
			UserList: user_list,
		})
	}

}
