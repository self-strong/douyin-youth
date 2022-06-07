package core

import (
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

// 用户注册提供 用户名、密码，昵称即可，用户名保证唯一。创建后返回id和权限token
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 判断token是否存在，存在说明账户存在了？应该是username
	user := DbFindUserName(username)
	if user.Username == username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			UserId:   0,
		})
	} else {

		user, err := DbRegister(username, password)

		if err != nil {

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
				UserId:   0,
			})
		}
		// newUser := User{
		// 	Id:            user.Id,
		// 	Name:          user.Name,
		// 	FollowCount:   user.FollowCount,
		// 	FollowerCount: user.FanCount,
		// 	IsFollow:      false,
		// }

		// userlogininfo更新
		LoginInfo[token] = UserLoginInfo{
			Id:       user.Id,
			username: user.Name,
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   user.Id,
			Token:    username + password,
		})
	}
}

// 使用用户名和密码登陆，返回用户id和token，进行页面信息显示
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// 通过token进行登陆
	userLoginInfo := DbFindUserLoginInfo(token) // 根据token获取用户信息

	if userLoginInfo != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   userLoginInfo.Id,
			Token:    token,
		})
	} else {
		//查表，先看是否存在用户名，再看密码是否正确

		user := DbFindUserName(username)

		if user.Username == username {

			if !DbCheckPwd(username, password) {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 2, StatusMsg: "Password Error"},
					UserId:   0,
				})
				return
			}
			// newUser := User{
			// 	Id:            user.Id,
			// 	Name:          username,
			// 	FollowCount:   user.FollowCount,
			// 	FollowerCount: user.FanCount,
			// 	IsFollow:      false, // 需要查表
			// }
			// usersLoginInfo[token] = newUser
			LoginInfo[token] = UserLoginInfo{
				Id:       user.Uid,
				username: user.Username,
			}

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Login Successful!"},
				UserId:   user.Uid,
				Token:    token,
			})

		} else {
			//查找不到该用户
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
				UserId:   0,
			})
		}
	}
}

func UserInfo(c *gin.Context) {
	// 访问的用户，自身的呢？
	token := c.Query("token")
	uIdStr := c.Query("id")

	uId, _ := strconv.ParseInt(uIdStr, 10, 64)

	userLoginInfo := DbFindUserLoginInfo(token) // 根据token获取用户信息

	if userLoginInfo == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			User:     User{},
		})
	} else {

		user := DbFindUserInfo(uId)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful"},
			User:     *user,
		})
	}

}
