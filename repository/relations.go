package repository

import "github.com/self-strong/douyin-youth/create_db"

// 粉丝关注偶像
func CreateRelation(fansId, idolId int64) error {
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("followings")

	relation := create_db.Following{
		FansId: fansId,
		IdolId: idolId,
	}
	// 是否要加锁
	res := tb.Create(&relation)
	// 更新user表

	UpdateFans(idolId, true)
	UpdateIdols(fansId, true)

	return res.Error

}

func DeleteRelation(fansId, idolId int64) error {
	if err := connect_db(); err != nil {
		return err
	}

	tb := db.Table("followings")
	var relation create_db.Following
	res := tb.Where("fans_id = ? AND idol_id = ?", fansId, idolId).Delete(&relation)

	// 用户表关注数减1
	UpdateFans(idolId, false)
	UpdateIdols(fansId, false)

	return res.Error
}

// 根据偶像id查询粉丝列表
func SearchFans(idolId int64) ([]create_db.Following, error) {
	var fans_list []create_db.Following

	if err := connect_db(); err != nil {
		return fans_list, err
	}

	tb := db.Table("followings")

	res := tb.Where("idol_id = ?", idolId).Find(&fans_list)

	return fans_list, res.Error
}

// 根据用户id查询关注列表
func SearchIdols(fansId int64) ([]create_db.Following, error) {
	var idol_list []create_db.Following

	if err := connect_db(); err != nil {
		return idol_list, err
	}

	tb := db.Table("followings")

	res := tb.Where("fans_id = ?", fansId).Find(&idol_list)

	return idol_list, res.Error

}
