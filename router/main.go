package router

import (
	"clipboard/controllers"
	"clipboard/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	auth := e.Group("/auth")
	{
		auth.POST("/login", controllers.LoginRoute)
		auth.POST("/register", controllers.RegisterRoute)
	}
	clip := e.Group("/clip")
	clip.Use(middlewares.PermissionMiddleware())
	{
		clip.GET("/", controllers.GetClipRoute)
		clip.POST("/", controllers.CreateClipRoute)
	}
}
