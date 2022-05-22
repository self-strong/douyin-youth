package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户
type User struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT"`            //用户id，设置为primary_key主键，AUTO_INCREMENT自增
	Name        string `gorm:"unique_index:UserName;not null;unique"` //设置为唯一索引
	Password    string `gorm:"type:varchar(64);not null"`
	FollowCount int64  `gorm:"default:0"` // 默认为0
	FanCount    int64  `gorm:"default:0"`
}

// 视频
type Video struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title        string `gorm:"index:VdeioTitle;not null"` //设置为普通索引
	CreateUid    int64  `gorm:"not null"`
	Timestamp    string `gorm:"not null"`
	PlayUrl      string `gorm:"not null"`
	CoverUrl     string `gorm:"not null"`
	ThumbCount   int64  `gorm:"default:0"`
	CommentCount int64  `gorm:"default:0"`
}

// 点赞
type Thumb struct {
	Uid       int64  `gorm:"index:Uid;not null"`
	Vid       int64  `gorm:"index:Vid;not null"`
	Timestamp string `gorm:"not null"`
}

// 关注
type Following struct {
	FansId int64 `gorm:"index:FansId;not null"` //粉丝数，设置为普通索引
	IdolId int64 `gorm:"index:IdolId;not null"` //关注数，设置为普通索引
}

// 评论
type Comment struct {
	CmId      int64  `gorm:"primary_key;AUTO_INCREMENT"` //评论id，设置为primary_key主键，AUTO_INCREMENT自增
	Vid       int64  `gorm:"index:Vid;not null"`         //粉丝数，设置为普通索引
	Uid       int64  `gorm:"not null"`
	Content   string `gorm:"not null"`
	Timestamp string `gorm:"not null"`
}

func main() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:hallo2014@tcp(127.0.0.1:3306)/douyin",
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		println(err)
		return
	}

	// 自动迁移创建表格
	err = db.AutoMigrate(&User{}, &Video{}, &Thumb{}, &Comment{}, &Following{})

	if err != nil {
		fmt.Println(err)
		return
	}

}
