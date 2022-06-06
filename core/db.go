package core

// 对数据库的操作

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库连接
var db *gorm.DB

// 数据库中单独使用的表

// DbUser 用户表
type DbUser struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT"`            //用户id，设置为primary_key主键，AUTO_INCREMENT自增
	Name        string `gorm:"unique_index:UserName;not null;unique"` //设置为唯一索引
	Password    string `gorm:"type:varchar(64);not null"`
	FollowCount int64  `gorm:"default:0"` // 默认为0
	FanCount    int64  `gorm:"default:0"`
}

// TableName 设置数据库表名
func (*DbUser) TableName() string {
	return "users"
}

// DbVideo 视频表
type DbVideo struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title        string `gorm:"index:VideoTitle;not null"` //设置为普通索引
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

// DbThumb 点赞表
type DbThumb struct {
	Uid       int64  `gorm:"index:Uid;not null"`
	Vid       int64  `gorm:"index:Vid;not null"`
	Timestamp string `gorm:"not null"`
}

func (*DbThumb) TableName() string {
	return "thumbs"
}

// DbFollowing 关注表
type DbFollowing struct {
	FansId int64 `gorm:"index:FansId;not null"` //粉丝数，设置为普通索引
	IdolId int64 `gorm:"index:IdolId;not null"` //关注数，设置为普通索引
}

func (*DbFollowing) TableName() string {
	return "followings"
}

// DbComment 评论表
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

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	Id       int64
	UserName string
}

// DbInsertVideoInfo 向数据库中插入用户发布的video记录
func DbInsertVideoInfo(uId int64, fileName, coverName string) *gorm.DB {
	playUrl := "http://192.168.1.4:8080/douyin/publish/video/?videoName=" + fileName
	coverUrl := "http://192.168.1.4:8080/douyin/publish/cover/?coverName=" + coverName
	// 可不可以将数据库插入请求存下来，待缓存区满之后再批量插入
	videoInfo := DbVideo{Title: fileName, CreateUid: uId, Timestamp: time.Now().String(), PlayUrl: playUrl, CoverUrl: coverUrl}
	result := db.Create(&videoInfo)
	return result
}

// LoginInfo 模拟用户登录时的Token和用户信息
var LoginInfo map[string]UserLoginInfo

// DbFindUserInfoByToken 根据Token获取登录用户的信息
func DbFindUserInfoByToken(token string) *UserLoginInfo {
	// 在数据库表中查询登录用户的TOKEN
	userLoginInfo, ok := LoginInfo[token]
	if ok { // 存在
		return &userLoginInfo
	}
	return nil
}

// DbInsertUserLoginInfo 插入用户登录的信息
func DbInsertUserLoginInfo(id int64, userName, token string) {
	LoginInfo[token] = UserLoginInfo{Id: id, UserName: userName}
}

// DbFindVideoList 获取发布视频列表 是不是可以利用缓存的思想来优化以下
func DbFindVideoList(user *User) []Video {
	var dbVideos []DbVideo
	db.Where("create_uid = ?", user.Uid).Find(&dbVideos)
	if dbVideos == nil {
		return nil
	}

	videos := make([]Video, len(dbVideos))
	for i := 0; i < len(dbVideos); i++ {
		videos[i].User = *user
		// 恢复Title的原名
		videos[i].Title = dbVideos[i].Title
		videos[i].PlayUrl = dbVideos[i].PlayUrl
		videos[i].CoverUrl = dbVideos[i].CoverUrl
		videos[i].Id = dbVideos[i].Id
		videos[i].CommentCount = dbVideos[i].CommentCount
		videos[i].ThumbCount = dbVideos[i].ThumbCount
		videos[i].Is_favorite = true
	}
	return videos
}

// DbFindUserInfoById 根据uId查找用户信息
func DbFindUserInfoById(uId int64) *User {
	var dbUser DbUser
	db.First(&dbUser, uId)
	if dbUser.Id == 0 {
		return nil
	}
	var user User
	user.Uid = dbUser.Id
	user.Username = dbUser.Name
	user.Follow = dbUser.FollowCount
	user.Following = dbUser.FanCount
	user.Is_follow = false // 这需要查表z
	return &user
	// 判断用户是否存在
}

// DbFindUserInfoByName 根据username查找用户信息
func DbFindUserInfoByName(username string) *User {
	var dbUser DbUser
	db.Table("users").Where("name = ?", username).Find(&dbUser)
	if dbUser.Id == 0 {
		return nil
	}
	var user User
	user.Uid = dbUser.Id
	user.Username = dbUser.Name
	user.Follow = dbUser.FollowCount
	user.Following = dbUser.FanCount
	user.Is_follow = false // 这需要查表
	return &user
}

// DbCheckUser 检查用户是否存在
// 返回值-1代表用户不存在，返回值0代表用户密码错误，其他则返回用户的ID
func DbCheckUser(username, password string) int64 {
	var dbUser DbUser
	db.Table("users").Where("name = ?", username).First(&dbUser)
	fmt.Println(dbUser)
	var ret int64 = -1
	if dbUser.Name == username {
		if dbUser.Password == password {
			ret = dbUser.Id
		} else {
			ret = 0
		}
	}
	return ret
}

// DbConnect 连接数据库
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
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	LoginInfo = map[string]UserLoginInfo{} // 初始化登录信息表

	return nil
	// 自动迁移创建表格
	// err = db.AutoMigrate(&User{}, &Video{}, &Thumb{}, &Comment{}, &Following{})
}

// DbFavoriteAction Thumb Up
func DbFavoriteAction(uId int64, vId int64) error {
	tx := db.Begin()
	thumbInfo := DbThumb{Uid: uId, Vid: vId, Timestamp: time.Now().String()}
	result := tx.Create(&thumbInfo)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	dbVideo := DbVideo{Id: vId}
	result = tx.Model(&dbVideo).Update("ThumbCount", gorm.Expr("ThumbCount + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// DbUnFavoriteAction Cancel Thumb Up
func DbUnFavoriteAction(uId int64, vId int64) error {
	tx := db.Begin()
	thumbInfo := DbThumb{Uid: uId, Vid: vId}
	result := tx.Delete(&thumbInfo)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	dbVideo := DbVideo{Id: vId}
	result = tx.Model(&dbVideo).Update("ThumbCount", gorm.Expr("ThumbCount - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// DbFavoriteList Fetch List
func DbFavoriteList(uId int64) []Video {
	var dbThumbs []DbThumb
	db.Where("Uid = ?", uId).Find(&dbThumbs)
	if dbThumbs == nil {
		return nil
	}

	favoriteVideos := make([]Video, len(dbThumbs))
	for i := 0; i < len(dbThumbs); i++ {
		var dbVideo DbVideo
		db.First(&dbVideo, dbThumbs[i].Vid)

		favoriteVideos[i].Id = dbVideo.Id
		favoriteVideos[i].Title = dbVideo.Title
		favoriteVideos[i].PlayUrl = dbVideo.PlayUrl
		favoriteVideos[i].CoverUrl = dbVideo.CoverUrl
		favoriteVideos[i].CommentCount = dbVideo.CommentCount
		favoriteVideos[i].ThumbCount = dbVideo.ThumbCount

		var author DbUser
		db.First(&author, dbVideo.CreateUid)

		var relation DbFollowing
		following := db.Where("FansId = ? AND IdolId = ?", uId, author.Id).First(&relation)

		favoriteVideos[i].User = User{
			Uid:       author.Id,
			Username:  author.Name,
			Follow:    author.FanCount,
			Following: author.FollowCount,
			Is_follow: following.RowsAffected > 0,
		}
	}

	return favoriteVideos
}

// DbPostComment by vID and content
func DbPostComment(uId int64, vId int64, text string) (error, Comment) {
	tx := db.Begin()
	comment := DbComment{Uid: uId, Vid: vId, Content: text, Timestamp: time.Now().String()}
	result := tx.Create(&comment)
	if result.Error != nil {
		tx.Rollback()
		return result.Error, Comment{}
	}

	dbVideo := DbVideo{Id: vId}
	result = tx.Model(&dbVideo).Update("CommentCount", gorm.Expr("CommentCount + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error, Comment{}
	}

	var user DbUser
	result = tx.First(&user, uId)
	if result.Error != nil {
		tx.Rollback()
		return result.Error, Comment{}
	}

	returnComment := Comment{
		CmId: comment.CmId,
		User: User{
			Uid:       user.Id,
			Username:  user.Name,
			Follow:    user.FanCount,
			Following: user.FollowCount,
			Is_follow: true,
		},
		Content:    text,
		CreateDate: comment.Timestamp,
	}
	return tx.Commit().Error, returnComment
}

// DbDeleteComment by cmId
func DbDeleteComment(cmId int64, vId int64) error {
	tx := db.Begin()
	comment := DbComment{CmId: cmId}
	if result := tx.Delete(&comment).Error; result != nil {
		tx.Rollback()
		return result
	}

	dbVideo := DbVideo{Id: vId}
	result := tx.Model(&dbVideo).Update("CommentCount", gorm.Expr("CommentCount - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// DbCommentList by vId
func DbCommentList(uId int64, vId int64) []Comment {
	var dbComments []DbComment
	db.Where("Vid = ?", vId).Find(&dbComments)
	if dbComments == nil {
		return nil
	}

	comments := make([]Comment, len(dbComments))
	for i := 0; i < len(dbComments); i++ {
		var author DbUser
		db.First(&author, dbComments[i].Uid)

		var relation DbFollowing
		following := db.Where("FansId = ? AND IdolId = ?", uId, author.Id).First(&relation)

		user := User{
			Uid:       author.Id,
			Username:  author.Name,
			Follow:    author.FanCount,
			Following: author.FollowCount,
			Is_follow: following.RowsAffected > 0,
		}

		comments[i] = Comment{
			CmId:       dbComments[i].CmId,
			User:       user,
			Content:    dbComments[i].Content,
			CreateDate: dbComments[i].Timestamp,
		}
	}

	return comments
}

// DbFollowAction uId -> toID
func DbFollowAction(uId int64, toId int64) error {
	tx := db.Begin()
	relation := DbFollowing{
		FansId: uId,
		IdolId: toId,
	}
	if err := tx.Create(&relation).Error; err != nil {
		tx.Rollback()
		return err
	}

	target := DbUser{Id: toId}
	if err := tx.Model(&target).Update("FanCount", gorm.Expr("FanCount + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	self := DbUser{Id: uId}
	if err := tx.Model(&self).Update("FollowCount", gorm.Expr("FollowCount + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DbUnFollowAction uID -> toId
func DbUnFollowAction(uId int64, toId int64) error {
	tx := db.Begin()
	relation := DbFollowing{
		FansId: uId,
		IdolId: toId,
	}
	if err := tx.Delete(&relation).Error; err != nil {
		tx.Rollback()
		return err
	}

	target := DbUser{Id: toId}
	if err := tx.Model(&target).Update("FanCount", gorm.Expr("FanCount - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	self := DbUser{Id: uId}
	if err := tx.Model(&self).Update("FollowCount", gorm.Expr("FollowCount - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DbFollowList get uId follows whom
func DbFollowList(uId int64, opId int64) []User {
	var dbUsers []DbFollowing
	result := db.Where("FansId = ?", uId).Find(&dbUsers)

	if result.RowsAffected == 0 {
		return nil
	}

	followList := make([]User, len(dbUsers))
	for i := 0; i < len(dbUsers); i++ {
		var dbUser DbUser
		db.First(&dbUser, dbUsers[i].IdolId)

		var followTestRelation DbFollowing
		followTest := db.Model(DbFollowing{FansId: opId, IdolId: dbUsers[i].IdolId}).First(&followTestRelation)
		followList[i] = User{
			Uid:       dbUser.Id,
			Username:  dbUser.Name,
			Follow:    dbUser.FanCount,
			Following: dbUser.FollowCount,
			Is_follow: followTest.RowsAffected > 0,
		}
	}

	return followList
}

func DbFollowerList(uId int64, opId int64) []User {
	var dbUsers []DbFollowing
	result := db.Where("IdolId = ?", uId).Find(&dbUsers)

	if result.RowsAffected == 0 {
		return nil
	}

	followerList := make([]User, len(dbUsers))
	for i := 0; i < len(dbUsers); i++ {
		var dbUser DbUser
		db.First(&dbUser, dbUsers[i].FansId)

		var followTestRelation DbFollowing
		followTest := db.Model(DbFollowing{FansId: opId, IdolId: dbUsers[i].FansId}).First(&followTestRelation)
		followerList[i] = User{
			Uid:       dbUser.Id,
			Username:  dbUser.Name,
			Follow:    dbUser.FanCount,
			Following: dbUser.FollowCount,
			Is_follow: followTest.RowsAffected > 0,
		}
	}

	return followerList
}

func DbRegister(username, password string) (DbUser, error) {

	tb := db.Table("users")
	// q := query.Use(db).User
	// //插入姓名
	user := DbUser{Name: username, Password: password, FollowCount: 0, FanCount: 0}
	res := tb.Create(&user)
	// err = q.WithContext(context.Background()).Create(&user)

	// if err != nil {
	// 	// println(err)
	// 	return user, err
	// }

	return user, res.Error
}

// DbFeed 未登陆时刷视频
func DbFeed() []Video {
	// var video_list []Video

	tb := db.Table("videos")

	var videos []DbVideo

	tb.Limit(30).Order("timestamp desc").Find(&videos) // 查找video信息

	videoList := make([]Video, len(videos))

	for i := 0; i < len(videos); i++ {
		// var dbVideo DbVideo
		dbVideo := videos[i]

		videoList[i].Id = dbVideo.Id
		videoList[i].Title = dbVideo.Title
		videoList[i].PlayUrl = dbVideo.PlayUrl
		videoList[i].CoverUrl = dbVideo.CoverUrl
		videoList[i].CommentCount = dbVideo.CommentCount
		videoList[i].ThumbCount = dbVideo.ThumbCount

		var author DbUser
		db.First(&author, dbVideo.CreateUid) // 视频发布的id

		// var relation DbFollowing
		// following := db.Where("FansId = ? AND IdolId = ?", uId, author.Id).First(&relation)

		videoList[i].User = User{
			Uid:       author.Id,
			Username:  author.Name,
			Follow:    author.FanCount,
			Following: author.FollowCount,
			Is_follow: false,
		}
	}

	return videoList
}

// DbFeedWithLogin 未登陆时发布视频
func DbFeedWithLogin(uId int64) []Video {
	// var video_list []Video

	tb := db.Table("videos")

	var videos []DbVideo

	tb.Limit(30).Order("timestamp desc").Find(&videos) // 查找video信息

	video_list := make([]Video, len(videos))

	for i := 0; i < len(videos); i++ {
		// var dbVideo DbVideo
		dbVideo := videos[i]

		video_list[i].Id = dbVideo.Id
		video_list[i].Title = dbVideo.Title
		video_list[i].PlayUrl = dbVideo.PlayUrl
		video_list[i].CoverUrl = dbVideo.CoverUrl
		video_list[i].CommentCount = dbVideo.CommentCount
		video_list[i].ThumbCount = dbVideo.ThumbCount

		var author DbUser
		db.First(&author, dbVideo.CreateUid) // 视频发布的id

		var relation DbFollowing
		following := db.Where("FansId = ? AND IdolId = ?", uId, author.Id).First(&relation)

		video_list[i].User = User{
			Uid:       author.Id,
			Username:  author.Name,
			Follow:    author.FanCount,
			Following: author.FollowCount,
			Is_follow: following.RowsAffected > 0,
		}
	}

	return video_list
}
