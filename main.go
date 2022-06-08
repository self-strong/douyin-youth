package main

import (
	"douyin/core"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	// 建立数据库连接
	err := core.DbConnect()

	// 数据库连接出错
	if err != nil {
		fmt.Println("Connect database error!")
		return
	}
	r := gin.Default()

	//public directory is used to serve static resources
	//r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//基础接口
	apiRouter.GET("/feed/", core.Feed)               // 视频流接口
	apiRouter.GET("/user/", core.UserInfo)           // 用户信息
	apiRouter.POST("/user/register/", core.Register) // 用户注册
	apiRouter.POST("/user/login/", core.Login)       // 用户登录

	//上传视频以及获取发布列表
	apiRouter.POST("/publish/action/", core.PublishAction) // 发布视频
	apiRouter.GET("/publish/list/", core.PublishList)      // 获取视频列表

	//扩展接口一
	apiRouter.POST("/favorite/action/", core.FavoriteAction)
	apiRouter.GET("/favorite/list/", core.FavoriteList)
	apiRouter.POST("/comment/action/", core.CommentAction)
	apiRouter.GET("/comment/list/", core.CommentList)

	// 扩展接口二
	apiRouter.POST("/relation/action/", core.RelationAction)
	apiRouter.GET("/relation/follow/list/", core.FollowList)
	apiRouter.GET("/relation/follower/list/", core.FollowerList)

	apiRouter.GET("/publish/video/", core.GetVideo)
	apiRouter.GET("/publish/cover/", core.GetCover)

	err = r.Run()
	if err != nil {
		return
	}
}
