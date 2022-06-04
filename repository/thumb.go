package repository

import (
	"time"

	"github.com/self-strong/douyin-youth/create_db"
	"gorm.io/gorm"
)

func CreateThumb(userId, videoId int64) (*gorm.DB, error) {
	if err := connect_db(); err != nil {
		return db, err
	}

	tb := db.Table("thumbs")
	// q := query.Use(db).User
	// 视频
	thumb := create_db.Thumb{Uid: userId, Vid: videoId, Timestamp: time.Now().String()}

	res := tb.Create(&thumb)

	if err := updateVideoThumb(videoId, true); err != nil {
		return res, err
	}

	return res, res.Error
}

func CancelThumb(userId, videoId int64) (*gorm.DB, error) {
	if err := connect_db(); err != nil {
		return db, err
	}

	tb := db.Table("thumbs")
	// q := query.Use(db).User
	// 视频

	var thumb create_db.Thumb
	res := tb.Where("uid = ? AND vid = ?", userId, videoId).Delete(&thumb)
	if err := updateVideoThumb(videoId, true); err != nil {
		return res, err
	}

	return res, res.Error
}

func SearchThumbVideo(userId int64) ([]int64, error) {
	var video_id []int64
	var thumb_list []create_db.Thumb

	if err := connect_db(); err != nil {
		return video_id, err
	}

	tb := db.Table("thumbs")

	res := tb.Where("uid = ?", userId).Find(&thumb_list) // 根据用户的id查询点赞的视频

	if res.Error != nil {
		return video_id, res.Error
	}
	for i := range thumb_list {
		video_id = append(video_id, thumb_list[i].Vid)
	}

	return video_id, nil
}
