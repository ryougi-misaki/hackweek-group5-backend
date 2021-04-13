package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hackathon/config"
	"hackathon/controller"
	"hackathon/middleware"
)

func Init() {
	r := gin.Default()
	r.Use(cors.Default())
	r = CollectRoute(r)
	panic(r.Run(fmt.Sprintf(":%d", config.Conf.Port)))
}

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("api/register", controller.Register)
	r.POST("api/login", controller.Login)
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.GET("api/auth/sendcode", controller.SendCode)
	r.PUT("api/forget", controller.ChangePwd)
	return r
}
