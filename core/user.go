package core

import (
	"fmt"
	"net/http"
	"strconv"

	// "github.com/self-strong/douyin-youth/repository"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register 用户注册提供 用户名、密码，昵称即可，用户名保证唯一。创建后返回id和权限token
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断token是否存在，存在说明账户存在了？应该是username
	user := DbFindUserInfoByName(username)
	if user != nil { // 如果用户已经存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			Token:    "",
			UserId:   0,
		})
	} else {

		uId, err := DbRegister(username, password)

		if err != nil {

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
				UserId:   0,
			})
			return
		}
		// newUser := User{
		// 	Id:            user.Id,
		// 	Name:          user.Name,
		// 	FollowCount:   user.FollowCount,
		// 	FollowerCount: user.FanCount,
		// 	IsFollow:      false,
		// }

		// userLoginInfo更新
		DbInsertUserLoginInfo(uId, username, username+password)
		fmt.Println(LoginInfo)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   uId,
			Token:    username + password,
		})
	}
}

// Login 使用用户名和密码登陆，返回用户id和token，进行页面信息显示
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 通过token进行登陆 ?
	userLoginInfo := DbFindUserInfoByToken(token) // 根据token获取用户信息

	if userLoginInfo != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   userLoginInfo.Id,
			Token:    token,
		})
	} else {
		//查表，检查用户是否存在
		fmt.Println(username, password, "++++++++++++++")
		var ret = DbCheckUser(username, password)
		if ret == -1 { // 用户不存在
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
				UserId:   0,
				Token:    "",
			})
		} else if ret == 0 { // 密码不正确
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 2, StatusMsg: "Password Error"},
				UserId:   0,
				Token:    "",
			})
		} else {
			DbInsertUserLoginInfo(ret, username, token)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Login Successful!"},
				UserId:   ret,
				Token:    token,
			})
		}

		//user := DbFindUserInfoByName(username)
		//
		//if user.Username == username {
		//
		//	if !DbCheckPwd(username, password) {
		//
		//		return
		//	}
		//	// newUser := User{
		//	// 	Id:            user.Id,
		//	// 	Name:          username,
		//	// 	FollowCount:   user.FollowCount,
		//	// 	FollowerCount: user.FanCount,
		//	// 	IsFollow:      false, // 需要查表
		//	// }
		//	// usersLoginInfo[token] = newUser
		//	LoginInfo[token] = UserLoginInfo{
		//		Id:       user.Uid,
		//		username: user.Username,
		//	}

		//
		//
		//} else {
		//	//查找不到该用户
		//
		//}
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	uIdStr := c.Query("user_id")

	uId, _ := strconv.ParseInt(uIdStr, 10, 64)

	userLoginInfo := DbFindUserInfoByToken(token) // 根据token获取用户信息

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			User:     User{},
		})
	} else {

		user := DbFindUserInfoById(uId)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful"},
			User:     *user,
		})
	}
}
