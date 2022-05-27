package main

import (
	"douyin/core"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	// 建立数据库连接
	err := core.DbConnect()

	if err != nil {
		fmt.Println("Connect database error!")
	}

	// 
	r := gin.Default()

	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	apiRouter.POST("/publish/action/", core.PublishAction)
	apiRouter.GET("/publish/list/", core.PublishList)

	r.Run()
}
