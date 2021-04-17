package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"hackathon/config"
	"hackathon/controller"
	_ "hackathon/docs"
	"hackathon/middleware"
	"io"
	"os"
)

func Init() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	r.Use(cors.Default())
	r = CollectRoute(r)
	panic(r.Run(fmt.Sprintf(":%d", config.Conf.Port)))
}

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//用户
	r.POST("api/register", controller.Register)
	r.POST("api/login", controller.Login)
	r.GET("api/info/:id", controller.Info)
	r.PUT("api/auth/me", middleware.AuthMiddleware(), controller.EditInfo)
	r.PUT("api/auth/pwd", middleware.AuthMiddleware(), controller.ChangePwd)

	//版块
	r.GET("api/tags", controller.RetrieveTags)
	r.POST("api/auth/tag", middleware.AuthMiddleware(), controller.CreateTag)
	r.DELETE("api/auth/tag/:id", middleware.AuthMiddleware(), controller.DeleteTag)

	//树洞
	r.GET("api/post/:id", controller.RetrievePost)
	r.POST("api/auth/post", middleware.AuthMiddleware(), controller.CreatePost)
	r.DELETE("api/auth/post/:id", middleware.AuthMiddleware(), controller.DeletePost)

	//管理员
	r.GET("api/admin/posts", middleware.AuthMiddleware(), controller.ShowAllPosts)
	r.PUT("api/admin/post/:id/:status", middleware.AuthMiddleware(), controller.UpdatePostStatus)

	return r
}
