package repository

import (
	"time"

	"github.com/self-strong/douyin-youth/create_db"
)

func SearchComment(videoId int64) ([]create_db.Comment, error) {

	var comments []create_db.Comment

	return comments, nil
}

func CreateComment(content string, userId, videoId int64) (create_db.Comment, error) {
	var comment create_db.Comment
	if err := connect_db(); err != nil {
		return comment, err
	}

	tb := db.Table("comments")
	// q := query.Use(db).User
	// 视频

	comment = create_db.Comment{
		Vid:       videoId,
		Uid:       userId,
		Content:   content,
		Timestamp: time.Now().String()}

	if res := tb.Create(&comment); res.Error != nil {
		return create_db.Comment{}, res.Error
	}

	// 更新评论数
	if err := updateVideoComment(videoId, true); err != nil {
		return create_db.Comment{}, err
	}

	return comment, nil
	// return create_db.Comment{}, nil
}

func DeleteComment(commentid int64) error {

	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("comments")
	// q := query.Use(db).User

	var comment create_db.Comment
	if res := tb.Where("cm_id = ? ", commentid).Delete(&comment); res.Error != nil {
		return res.Error
	}
	if err := updateVideoComment(comment.Vid, false); err != nil {
		return err
	}

	return nil
}
