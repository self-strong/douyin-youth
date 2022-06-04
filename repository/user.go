package repository

import (
	"fmt"

	"github.com/self-strong/douyin-youth/create_db"
)

func Register(name string, pwd string) (create_db.User, error) {

	// dsn := "root:hallo2014@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN:                       dsn,
	// 	DefaultStringSize:         256,   // string 类型字段的默认长度
	// 	DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	// }), &gorm.Config{})

	if err := connect_db(); err != nil {
		// println(err)
		return create_db.User{}, err
	}

	q := db.Table("users")
	// q := query.Use(db).User
	// //插入姓名
	user := create_db.User{Name: name, Password: pwd, FollowCount: 0, FanCount: 0}
	q.Create(&user)
	// err = q.WithContext(context.Background()).Create(&user)

	// if err != nil {
	// 	// println(err)
	// 	return user, err
	// }

	return user, nil
}

func SearchUsername(name string) (create_db.User, error) {
	// dsn := "root:hallo2014@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN: dsn}), &gorm.Config{})

	if err := connect_db(); err != nil {
		// println(err)
		return create_db.User{}, err
	}
	// if err != nil {
	// 	// println(err)
	// 	return create_db.User{}, err
	// }

	tb := db.Table("users")

	// user := tb.Where("name = ?", "name") // 查找表中名字为name的用户
	var user_ create_db.User
	// fmt.Println(err, user_)

	// 这个怎么判断呢
	// err = tb.Where("name = ?", name).Find(&user_).Error
	res := tb.Where("name = ?", name).Find(&user_)
	// fmt.Println(err, user_)
	// if user_.Name != name {
	// 	return false, nil //不存在用户名
	// } else {
	// 	return true, nil
	// }
	return user_, res.Error

}

//查询单个用户
func SearchOneUserById(userId int64) (create_db.User, error) {
	if err := connect_db(); err != nil {
		return create_db.User{}, err
	}

	tb := db.Table("users")
	var user_ create_db.User
	err := tb.Where("id = ?", userId).Find(&user_)

	return user_, err.Error
}

//查询多个用户
func SearchUserById(userId []int64) ([]create_db.User, error) {
	// var user []create_db.User
	if err := connect_db(); err != nil {
		return nil, err
	}

	tb := db.Table("users")

	user := make([]create_db.User, len(userId))

	for i := 0; i < len(userId); i++ {

		res := tb.Where("id = ?", userId[i]).Find(&(user[i]))
		fmt.Println(userId[i], user[i], res.Error, "数据库")
		if res.Error != nil {
			return nil, res.Error
		}
	}

	// for i := range userId {

	// 	var user_ create_db.User
	// 	res := tb.Where("id = ?", userId[i]).Find(&user_)

	// 	if res.Error != nil {
	// 		return user, res.Error
	// 	}
	// 	user = append(user, user_)
	// }

	return user, nil
}

// 更新粉丝数，flag为true表示增加，
func UpdateFans(userId int64, flag bool) error {

	var user create_db.User
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("users")

	// err := tb.Where("id = ?", userId).Find(&user)

	// fanscnt := user.FanCount + 1
	// true表示粉丝加1
	if flag {
		tb.Where("id = ?", userId).Find(&user)
		cnt := user.FanCount + 1
		res := tb.Where("id = ?", userId).Update("fan_count", cnt)
		return res.Error
	} else {
		tb.Where("id = ?", userId).Find(&user)
		cnt := user.FanCount - 1
		res := tb.Where("id = ?", userId).Update("fan_count", cnt)
		return res.Error
	}

	// db.Model(&user).Where("active = ?", true).Update("name", "hello")

}

// 更新关注数，flag为true表示增加;A关注B，A的偶像数+1
func UpdateIdols(userId int64, flag bool) error {

	var user create_db.User
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("users")

	if flag {

		tb.Where("id = ?", userId).Find(&user)
		cnt := user.FollowCount + 1
		res := tb.Where("id = ?", userId).Update("follow_count", cnt)
		return res.Error
	} else {
		tb.Where("id = ?", userId).Find(&user)
		cnt := user.FollowCount - 1
		res := tb.Where("id = ?", userId).Update("follow_count", cnt)
		return res.Error
	}
	// err := tb.Where("id = ?", userId).Find(&user)

	// fanscnt := user.FanCount + 1

	// tb.Where()

	// db.Model(&user).Where("active = ?", true).Update("name", "hello")
	// return res.Error
}
