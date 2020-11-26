package router

import (
	"clipboard/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	auth := e.Group("/auth")
	{
		auth.GET("/login")
	}
	clip := e.Group("/clip")
	{
		clip.GET("/", controllers.GetClipRoute)
		clip.POST("/", controllers.CreateClipRoute)
	}
}
