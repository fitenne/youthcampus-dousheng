package main

import (
	"github.com/fitenne/youthcampus-dousheng/internal/common/mid"
	"github.com/fitenne/youthcampus-dousheng/internal/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)

	//在需要鉴权的接口类似的使用token鉴权
	apiRouter.GET("/user/", mid.JWTAuthMiddleware(), controller.UserInfo)

	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", mid.JWTAuthMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", mid.JWTAuthMiddleware(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", mid.JWTAuthMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", mid.JWTAuthMiddleware(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", mid.JWTAuthMiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", mid.JWTAuthMiddleware(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", mid.JWTAuthMiddleware(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", mid.JWTAuthMiddleware(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", mid.JWTAuthMiddleware(), controller.FollowerList)
}
