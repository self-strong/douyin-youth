package controller

import (
	"fmt"
	"net/http"

	"github.com/self-strong/douyin-youth/repository"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

// 需要初始化这样一个变量，再有用户添加时可以插入，
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// 用户id的序列号，用mysql主键自增获取
var userIdSequence = int64(1)

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

	// 判断token是否存在，存在说明账户存在了？应该是username   存在check为true
	check, err := repository.CheckUsername(username)

	fmt.Println(check)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}

	if check {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {

		// if _, exist := usersLoginInfo[token]; exist {
		// 	c.JSON(http.StatusOK, UserLoginResponse{
		// 		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		// 	})
		// } else {

		//  插入id 姓名，关注列表，点赞数
		user, err := repository.Register(username, password)

		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}

		newUser := User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FanCount,
			IsFollow:      false,
		}

		// userlogininfo更新
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
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
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

// 获取用户的id,昵称，返回关注数和粉丝数
func UserInfo(c *gin.Context) {
	// 访问的用户，自身的呢？
	token := c.Query("token")
	// user_id := c.Query("id")

	if user, exist := usersLoginInfo[token]; exist {

		// 查看是否关注
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

// func Register_db(name string) (int64, error) {
// 	db, err := gorm.Open(mysql.New(mysql.Config{
// 		DSN:                       "root:hallo2014@tcp(127.0.0.1:3306)/douyin",
// 		DefaultStringSize:         256,   // string 类型字段的默认长度
// 		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
// 		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
// 		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
// 		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
// 	}), &gorm.Config{})
// 	if err != nil {
// 		// println(err)
// 		return 0, err
// 	}
// 	// q := query.Use(db).User
// 	// //插入姓名
// 	// user := user_model.User{Name: name, FollowCount: 0, FollowerCount: 0, IsFollow: false}
// 	// err = q.WithContext(context.Background()).Create(&user)

// 	if err != nil {
// 		// println(err)
// 		return 0, err
// 	}

// 	// return user.ID, nil
// }
