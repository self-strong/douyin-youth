package main

import (
	"github.com/gin-gonic/gin"
	"github.com/self-strong/douyin-youth/controller"
)

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run()
}

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)                //获取视频流
	apiRouter.GET("/user/", controller.UserInfo)            //用户信息			// 计算用户的粉丝、关注、是否关注
	apiRouter.POST("/user/register/", controller.Register)  //用户注册接口		基本功能
	apiRouter.POST("/user/login/", controller.Login)        //用户登陆接口		基本功能？是否要初始化userlogininfo那个map
	apiRouter.POST("/publish/action/", controller.Publish)  //视频投稿
	apiRouter.GET("/publish/list/", controller.PublishList) //发布列表

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
