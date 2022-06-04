package repository

import (
	"time"

	"github.com/self-strong/douyin-youth/create_db"
)

func Publish(paly_url string, title string, cover_url string, uid int64) (create_db.Video, error) {

	if err := connect_db(); err != nil {
		// println(err)
		return create_db.Video{}, err
	}

	q := db.Table("videos")
	// q := query.Use(db).User
	// //插入姓名
	video := create_db.Video{Title: title, CreateUid: uid, PlayUrl: paly_url, CoverUrl: cover_url, Timestamp: time.Now().String()}
	q.Create(&video)
	// err = q.WithContext(context.Background()).Create(&user)

	return video, nil
}

// func CheckUsername(name string) (bool, error) {
// 	dsn := "root:hallo2014@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.New(mysql.Config{
// 		DSN: dsn}), &gorm.Config{})

// 	if err != nil {
// 		// println(err)
// 		return false, err
// 	}
// 	tb := db.Table("users")

// 	// user := tb.Where("name = ?", "name") // 查找表中名字为name的用户
// 	var user_ create_db.User
// 	// fmt.Println(err, user_)

// 	// 这个怎么判断呢
// 	err = tb.Where("name = ?", name).Find(&user_).Error

// 	fmt.Println(err, user_)
// 	if user_.Name != name {
// 		return false, nil //不存在用户名
// 	} else {
// 		return true, nil
// 	}

// }

func Feed() ([]create_db.Video, error) {
	if err := connect_db(); err != nil {
		return []create_db.Video{}, err
	}

	// 获取最近30条视频的列表

	tb := db.Table("videos")

	var videos []create_db.Video

	tb.Limit(30).Order("timestamp desc").Find(&videos)

	return videos, nil

}

func SearchVideoById(videoId []int64) ([]create_db.Video, error) {
	// 根据视频id序列查找
	var video_list []create_db.Video
	return video_list, nil
}

// 更新视频的评论数，flag为true表示+1
func updateVideoComment(videoId int64, flag bool) error {
	var video create_db.Video
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("videos")

	// err := tb.Where("id = ?", userId).Find(&user)

	// fanscnt := user.FanCount + 1
	// true表示粉丝加1
	if flag {
		tb.Where("id = ?", videoId).Find(&video)
		res := tb.Where("id = ?", videoId).Update("comment_count", video.CommentCount+1)
		return res.Error
	} else {
		tb.Where("id = ?", videoId).Find(&video)
		res := tb.Where("id = ?", videoId).Update("comment_count", video.CommentCount-1)
		return res.Error
	}
}

// 更新视频的评论数，flag为true表示+1
func updateVideoThumb(videoId int64, flag bool) error {
	var video create_db.Video
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("videos")

	// err := tb.Where("id = ?", userId).Find(&user)

	// fanscnt := user.FanCount + 1
	// true表示粉丝加1
	if flag {
		tb.Where("id = ?", videoId).Find(&video)
		res := tb.Where("id = ?", videoId).Update("comment_count", video.ThumbCount+1)
		return res.Error
	} else {
		tb.Where("id = ?", videoId).Find(&video)
		res := tb.Where("id = ?", videoId).Update("comment_count", video.ThumbCount-1)
		return res.Error
	}
}
