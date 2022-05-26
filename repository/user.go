package repository

import (
	"fmt"

	"github.com/self-strong/douyin-youth/create_db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Register(name string, pwd string) (create_db.User, error) {

	dsn := "root:hallo2014@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		// println(err)
		return create_db.User{}, err
	}

	q := db.Table("users")
	// q := query.Use(db).User
	// //插入姓名
	user := create_db.User{Name: name, Password: pwd, FollowCount: 0, FanCount: 0}
	q.Create(&user)
	// err = q.WithContext(context.Background()).Create(&user)

	if err != nil {
		// println(err)
		return user, err
	}

	return user, nil
}

func CheckUsername(name string) (bool, error) {
	dsn := "root:hallo2014@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn}), &gorm.Config{})

	if err != nil {
		// println(err)
		return false, err
	}
	tb := db.Table("users")

	// user := tb.Where("name = ?", "name") // 查找表中名字为name的用户
	var user_ create_db.User
	// fmt.Println(err, user_)

	// 这个怎么判断呢
	err = tb.Where("name = ?", name).Find(&user_).Error

	fmt.Println(err, user_)
	if user_.Name != name {
		return false, nil //不存在用户名
	} else {
		return true, nil
	}

}
