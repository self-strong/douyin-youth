package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Name        string `gorm:"unique_index:UserName;not null;unique"`
	Password    string `gorm:"type:varchar(64);not null"`
	Followcount int64  `gorm:"default:0"`
	Fancount    int64  `gorm:"default:0"`
}

type Video struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title        string `gorm:"index:VdeioTitle;not null"`
	CreateUid    int64  `gorm:"not null"`
	Timestamp    string `gorm:"not null"`
	PlayUrl      string `gorm:"not null"`
	CoverUrl     string `gorm:"not null"`
	ThumbCount   int64  `gorm:"default:0"`
	CommentCount int64  `gorm:"default:0"`
}

type Thumb struct {
	Uid       int64  `gorm:"not null"`
	Vid       int64  `gorm:"not null"`
	Timestamp string `gorm:"not null"`
}

type Following struct {
	FansId int64 `gorm:"index:FansId;not null"`
	IdolId int64 `gorm:"index:IdolId;not null"`
}

type Comment struct {
	CmId      int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Vid       int64  `gorm:"not null"`
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

	// 迁移 schema
	err = db.AutoMigrate(&User{}, &Video{}, &Thumb{}, &Comment{}, &Following{})

	if err != nil {
		fmt.Println(err)
		return
	}

}
