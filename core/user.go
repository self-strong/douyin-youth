package core

import (
	"fmt"
	"net/http"
	"strconv"

	//"douyin/dto"
	"douyin/pkg/jwt"
	// "github.com/self-strong/douyin-youth/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	//密码加密后放入数据库
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hashedPassword)
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

		// userLoginInfo更新
		token, _ := jwt.GenToken(username)
		DbInsertUserLoginInfo(uId, username, token)
		fmt.Println(LoginInfo)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   uId,
			Token:    token,
		})
	}
}

// 登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := DbFindUserInfoByName(username)
	token, _ := jwt.GenToken(username)
	if user != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful!"},
			UserId:   user.Uid,
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

	}
}

func UserInfo(c *gin.Context) {
	//token := c.Query("token")
	uIdStr := c.Query("user_id")

	uId, _ := strconv.ParseInt(uIdStr, 10, 64)

	user := DbFindUserInfoById(uId) // 根据token获取用户信息

	if user == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			User:     User{},
		})
	} else {

		user.IsFollow = DbCheckIsFollow(user.Uid, uId)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful"},
			User:     *user,
		})
	}
}
