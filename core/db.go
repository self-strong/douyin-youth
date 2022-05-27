package core

// 对数据库的操作

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库连接
var db *gorm.DB

// 数据库中单独使用的表

// 用户表
type DbUser struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT"`            //用户id，设置为primary_key主键，AUTO_INCREMENT自增
	Name        string `gorm:"unique_index:UserName;not null;unique"` //设置为唯一索引
	Password    string `gorm:"type:varchar(64);not null"`
	FollowCount int64  `gorm:"default:0"` // 默认为0
	FanCount    int64  `gorm:"default:0"`
}

// 设置数据库表名
func (*DbUser) TableName() string {
	return "users"
}

// 视频表
type DbVideo struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title        string `gorm:"index:VdeioTitle;not null"` //设置为普通索引
	CreateUid    int64  `json:"create_uid" gorm:"not null"`
	Timestamp    string `json:"timestamp" gorm:"not null"`
	PlayUrl      string `json:"play_url" gorm:"not null"`
	CoverUrl     string `json:"cover_url" gorm:"not null"`
	ThumbCount   int64  `json:"thumb_count" gorm:"default:0"`
	CommentCount int64  `json:"comment_count" gorm:"default:0"`
}
func (*DbVideo) TableName() string {
	return "videos"
}

// 点赞表
type DbThumb struct {
	Uid       int64  `gorm:"index:Uid;not null"`
	Vid       int64  `gorm:"index:Vid;not null"`
	Timestamp string `gorm:"not null"`
}

func (*DbThumb) TableName() string {
	return "thumbs"
}

// 关注表
type DbFollowing struct {
	FansId int64 `gorm:"index:FansId;not null"` //粉丝数，设置为普通索引
	IdolId int64 `gorm:"index:IdolId;not null"` //关注数，设置为普通索引
}

func (*DbFollowing) TableName() string {
	return "followings"
}

// 评论表
type DbComment struct {
	CmId      int64  `gorm:"primary_key;AUTO_INCREMENT"` //评论id，设置为primary_key主键，AUTO_INCREMENT自增
	Vid       int64  `gorm:"index:Vid;not null"`         //粉丝数，设置为普通索引
	Uid       int64  `gorm:"not null"`
	Content   string `gorm:"not null"`
	Timestamp string `gorm:"not null"`
}
func (*DbComment) TableName() string {
	return "comments"
}

// 用户登录信息
type UserLoginInfo struct {
	Id       int64
	username string
}

// 向数据库中插入用户发布的video记录
func DbInsertVideoInfo(uId int64, fileName, videoPath, coverPath string) *gorm.DB {
	// 可不可以将数据库插入请求存下来，待缓存区满之后再批量插入
	videoInfo := DbVideo{Title: fileName, CreateUid: uId, Timestamp: time.Now().String(), PlayUrl: videoPath, CoverUrl: coverPath}
	result := db.Create(&videoInfo)
	return result
}

// 模拟用户登录时的Token和用户信息
var LoginInfo = map[string]UserLoginInfo{
	"token": {
		Id:       1,
		username: "king",
	},
}

// 根据Token获取登录用户的信息
func DbFindUserLoginInfo(token string) *UserLoginInfo {
	// 在数据库表中查询登录用户的TOKEN
	userLoginInfo , ok := LoginInfo[token]
	if ok { // 存在
		return &userLoginInfo
	}
	return nil
}

// 获取发布视频列表 是不是可以利用缓存的思想来优化以下
func DbFindVideoList(user *User) []Video {
	var dbVideos []DbVideo
	db.Where("create_uid = ?", user.Uid).Find(&dbVideos)
	if dbVideos == nil {
		return nil
	}

	videos := make([]Video, len(dbVideos))
	for i := 0; i < len(dbVideos); i++ {
		var vdo Video;
		vdo.User = *user
		// 恢复Title的原名
		vdo.Title = dbVideos[i].Title
		vdo.PlayUrl = dbVideos[i].PlayUrl
		vdo.CoverUrl = dbVideos[i].CoverUrl
		vdo.Id = dbVideos[i].Id
		vdo.CommentCount = dbVideos[i].CommentCount
		vdo.ThumbCount = dbVideos[i].ThumbCount
		vdo.Is_favorite = true
	}
	return videos
}

// 根据uId查找用户信息
func DbFindUserInfo(uId int64) *User {
	var dbUser DbUser
	db.First(&dbUser, uId)
	var user User
	user.Uid = dbUser.Id
	user.Username = dbUser.Name
	user.Follow = 0
	user.Following = 0
	user.Is_follow = false
	return &user
	// 判断用户是否存在
}

// 连接数据库
func DbConnect() error {
	// 是否已有数据库连接
	if db != nil {
		return nil
	}
	// 配置mysql,用户名、密码；
	dsn := "root:wb20010115@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN:                       dsn,
	// 	DefaultStringSize:         256,   // string 类型字段的默认长度
	// 	DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	// }), &gorm.Config{})
	var err error
	db, err = gorm.Open((mysql.Open(dsn)), &gorm.Config{})

	if err != nil {
		return err
	}

	return nil
	// 自动迁移创建表格
	// err = db.AutoMigrate(&User{}, &Video{}, &Thumb{}, &Comment{}, &Following{})
}
